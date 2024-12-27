## Logging
__IMPORTANT:__ for logging to work correctly, the configuration must be set up correctly.

This package contains logging tools. The following 
configurations are used for the functionality of this package:
```
DebugRelativeFilepath: true
ErrorLogging: true
ErrorLoggingPath: “error.log”
RequestInfoLog: true
RequestInfoLogPath: “request_info.log”
```
This package also contains a list of prefixes(*P_*) for logging. 
These prefixes can be used in appropriate locations.

__WriteLog__
```
WriteLog(skipLevel int, filePath string, flag int, prefix string, message string)
```
The main function for logging a message. It should be used for any type of logging.

__LogError__
```
LogError(message string)
```
Logs the error message.

__ErrorLoggingIfEnableAndWrite__
```
ErrorLoggingIfEnableAndWrite(w http.ResponseWriter, errorText string, writeText string)
```
Logging errors if enabled in the configuration and writing text to a web page.

__ErrorLogginIfEnable__
```
ErrorLogginIfEnable(message string)
```
Logging errors if enabled in the configuration.

__ClearRequestInfoLogging__
```
ClearRequestInfoLogging() error
```
Clears the file that logs the request. In the standard implementation, it is used before each request.

__LogRequestInfo__
```
LogRequestInfo(prefix string, message string)
```
Logs information about the request.

__RequestLogginIfEnable__
```
RequestLogginIfEnable(prefix string, message string)
```
Logs information about the request if it is enabled in the settings.