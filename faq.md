# Frequently Asked Questions

1. **Can this module collect metrics like cpu, memory usage etc?**

    A: Support is limited to stdout and stderr container logs. 

1. **Can the output be filtered to specific logs (eg. stderr only) or containers?**

    A: Yes, check the [logspout documentation on filtering](https://github.com/gliderlabs/logspout#including-specific-containers). For example, you can use the following ```cmd``` in **Container Create Options** to send only stderr logs:

    ```
    .
    .
    "Cmd": [
            "loganalytics://?filter.sources=stderr"
    ],
    .
    .
    ```

1. **Can you use the device twin properties for this module to enable/disable tracing?**

    A: This is not yet supported. However, the recommended way to achieve a similar effect is:

    * Create two [IoT Edge deployments](https://docs.microsoft.com/en-us/azure/iot-edge/how-to-deploy-monitor). Both are identical except one of them has the logging module defined and targets devices with a condition like ```tags.logging='enabled'```.
    * Next add the logging tag to the twin of devices you want to send logs from.
    * Presuming all other target conditions match, the device tagged for logging will get the logging module pushed to it as part of the deployment.
    * Remove the logging tag for the device to disable logging on it.

1. **Is this the officially supported monitoring/logging/diagnotics for Azure IoT Edge?**

    A: No, this is just a sample usage for Log analytics for cloud logging from edge devices. This project is best efffort, and is not offically supported.

5. **Why does this not work with Windows containers?**

    A: This project is simply an adapter for [logspout](https://github.com/gliderlabs/logspout). And logspout doesn't currently support Windows containers.