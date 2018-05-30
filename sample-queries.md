# Sample queries
Sample log analytics queries for data logged by [logspout-loganalytics IoT edge module](https://github.com/veyalla/logspout-loganalytics). Detailed reference for the Log Analytics query language [is also available](https://docs.loganalytics.io/docs/Learn/Getting-Started/Getting-started-with-the-Analytics-portal).

* **Sample 1**: *Show me all logs in the last 36 hours for device with IoT hub device id logspoutTest, sorted in descending order of TimeGenerated:*
    ```
    search *
    | project TimeGenerated, Level, msg_s, moduleName_s, iothubdeviceid_s, hostname_s, iothubname_s
    | sort by TimeGenerated desc nulls last
    | where iothubdeviceid_s == "logspoutTest"
    | where TimeGenerated > ago(36h)
    ```
* **Sample 2**: *Show me all stderr logs from edgeHub module in the last 12 hours for device with IoT hub device id logspoutTest, sorted in descending order of TimeGenerated:*
    ```
    search *
    | sort by TimeGenerated desc nulls last 
    | where Level == "stderr" and moduleName_s == "/edgeHub" and iothubdeviceid_s == "logspoutTest"
    | project TimeGenerated, Level, msg_s, moduleName_s, iothubdeviceid_s, hostname_s, iothubname_s
    | where TimeGenerated > ago(12h)
    ```

* **Sample 3**: *Show me all error logs from tempSensor module in the last 12 days for all devices with IoT hubs named like test-iot-hub, sorted in descending order of TimeGenerated:*
    ```
    search *
    | sort by TimeGenerated desc nulls last 
    | where Level == "stderr" and moduleName_s == "/tempSensor" and iothubname_s contains "test-iot-hub"
    | project TimeGenerated, Level, msg_s, moduleName_s, iothubdeviceid_s, hostname_s, iothubname_s
    | where TimeGenerated > ago(12d)
    ```

* **Sample 4**: *At the time error occured on device with ID logspoutTest, show me logs from all modules on that device 5 mins before and 30 secs after:*
    ```
    let startDatetime = todatetime('2018-05-27T08:51:09.976');
    let incidentDuration = totimespan(300s);

    search *
    | sort by TimeGenerated desc nulls last
    | where iothubdeviceid_s == "logspoutTest"
    | where TimeGenerated  between((startDatetime-incidentDuration) .. (startDatetime+totimespan(30s)) )
    | extend timeFromStart = TimeGenerated  - startDatetime
    | project TimeGenerated, Level, moduleName_s, msg_s
    ```
* **Sample 5**: *Show me the total number of logs received over the past 3 days:*
    ```
    search *  
    | where TimeGenerated > ago(3d)
    | count
    ```
* **Sample 6**: *Show me in MB, the size of all logs over the past 30 days:*
    ```
    Usage
    | where TimeGenerated > ago(30d)
    | summarize sum(Quantity)
    ```




