package loganalytics

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/gliderlabs/logspout/router"
)

const (
	envWorkspaceID       = "LOGANALYTICS_WORKSPACE_ID"
	envWorkspaceSecret   = "LOGANALYTICS_WORKSPACE_SECRET"
	envIoTHubName        = "IOTEDGE_IOTHUBHOSTNAME"
	envIoTHubDeviceID    = "IOTEDGE_DEVICEID"
	envGatewayHostName   = "IOTEDGE_GATEWAYHOSTNAME"
	envEdgeHubConnString = "EdgeHubConnectionString"
)

var (
	locationGMT = time.FixedZone("GMT", 0)
)

// LogClient is the client for log analytics
type LogClient struct {
	workspaceID     string
	workspaceSecret string
	httpClient      *http.Client
	hostname        string
	iotHubName      string
	iotHubDeviceID  string
	signingKey      []byte
	apiLogsURL      string
}

// ModuleMessage defines a log message from an IoT Edge module.
type ModuleMessage struct {
	V              int       `json:"version"`
	Time           time.Time `json:"timeEmitted"`
	Level          string    `json:"level"`
	Data           string    `json:"msg"`
	ModuleName     string    `json:"moduleName"`
	ContainerID    string    `json:"containerID"`
	ContainerImage string    `json:"containerImage"`
	Hostname       string    `json:"hostname"`
	IoTHubName     string    `json:"iothubname"`
	IoTHubDeviceID string    `json:"iothubdeviceid"`
}

// NewLogClient creates a log client
func NewLogClient(workspaceID, workspaceSecret string) LogClient {
	client := LogClient{
		workspaceID:     workspaceID,
		workspaceSecret: workspaceSecret,
	}

	client.httpClient = &http.Client{Timeout: time.Second * 30}
	client.signingKey, _ = base64.StdEncoding.DecodeString(workspaceSecret)
	client.apiLogsURL = fmt.Sprintf("https://%s.ods.opinsights.azure.com/api/logs?api-version=2016-04-01", workspaceID)
	client.iotHubName = os.Getenv(envIoTHubName)
	if len(client.iotHubName) == 0 {
		connStr := os.Getenv(envEdgeHubConnString)
		reg := regexp.MustCompile("HostName=(.*?);GatewayHostName=(.*?);DeviceId=(.*?);(.*?)")
		matches := reg.FindStringSubmatch(connStr)
		if len(matches) == 5 {
			client.iotHubName = matches[1]
			client.hostname = matches[2]
			client.iotHubDeviceID = matches[3]
		}
	}

	if len(client.iotHubDeviceID) == 0 {
		client.iotHubDeviceID = os.Getenv(envIoTHubDeviceID)
		client.hostname = os.Getenv(envGatewayHostName)
	}

	return client
}

// PostMessage logs an array of messages to log analytics service
func (c *LogClient) PostMessage(message *router.Message, timestamp time.Time) error {
	if timestamp.IsZero() {
		timestamp = time.Now().UTC()
	}

	msg := ModuleMessage{
		V:              0,
		Time:           timestamp,
		Level:          message.Source,
		Data:           message.Data,
		ModuleName:     message.Container.Name,
		ContainerID:    message.Container.ID,
		ContainerImage: message.Container.Config.Image,
		Hostname:       c.hostname,
		IoTHubDeviceID: c.iotHubDeviceID,
		IoTHubName:     c.iotHubName,
	}

	body, _ := json.Marshal(msg)
	req, _ := http.NewRequest(http.MethodPost, c.apiLogsURL, bytes.NewReader(body))

	date := time.Now().In(locationGMT).Format(time.RFC1123)
	stringToSign := "POST\n" + strconv.FormatInt(req.ContentLength, 10) + "\napplication/json\n" + "x-ms-date:" + date + "\n/api/logs"

	signature := computeHmac256(stringToSign, c.signingKey)

	req.Header.Set("Authorization", fmt.Sprintf("SharedKey %s:%s", c.workspaceID, signature))
	req.Header.Add("Log-Type", "container_logs")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-ms-date", date)
	req.Header.Set("time-generated-field", "Timestamp")

	response, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		defer response.Body.Close()
		buf, _ := ioutil.ReadAll(response.Body)

		time.AfterFunc(
			time.Second*15,
			func() {
				c.PostMessage(message, timestamp)
			})

		return fmt.Errorf("[loganalytics][%s] Post log request failed with status: %d %s", time.Now().UTC().Format(time.RFC3339), response.StatusCode, string(buf))

	}

	return nil
}

func init() {
	router.AdapterFactories.Register(NewLogAnalyticsAdapter, "loganalytics")
}

// ComputeHmac256 computes HMAC with given secret
func computeHmac256(message string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func NewLogAnalyticsAdapter(route *router.Route) (router.LogAdapter, error) {
	workspaceID, workspaceSecret := os.Getenv(envWorkspaceID), os.Getenv(envWorkspaceSecret)
	if workspaceID == "" || workspaceSecret == "" {
		return nil,
			fmt.Errorf("Workspace Id and secret not defined in environment variable '%s' and '%s'.\n", envWorkspaceID, envWorkspaceSecret)
	}

	client := NewLogClient(workspaceID, workspaceSecret)

	return &Adapter{
		route:  route,
		client: &client,
	}, nil
}

// Adapter defines a logspout adapter for azure log analytics.
type Adapter struct {
	route  *router.Route
	client *LogClient
}

// Stream waits on a logspout message channel. Upon receiving on it POSTs it to
// Log Analytics endpoint.
func (adapter *Adapter) Stream(logstream chan *router.Message) {
	for message := range logstream {
		adapter.client.PostMessage(message, time.Now().UTC())
	}
}
