## IManagerConfig
Saves and provides access to some settings.

__Debug__
```
Debug(enable bool)
```
Enables or disables the debug mod. If this option is enabled, the web page will 
display some errors.

__IsDebug__
```
IsDebug() bool
```
Returns Debug.

__PrintLog__
```
PrintLog(enable bool)
```
Enables the console log. That is, some data will be output to the console.

__IsPrintLog__
```
IsPrintLog() bool
```
Returns the value of PrintLog.

__ErrorLogging__
```
ErrorLogging(enable bool)
```
Enables logging of selected errors to the .log file. The path to the file must be set 
in the __ErrorLoggingFile__ method.

__IsErrorLogging__
```
IsErrorLogging() bool
```
Returns the ErrorLogging value.

__ErrorLoggingFile__
```
ErrorLoggingFile(path string)
```
Sets the path to the log file. If the file does not exist, it will be created.

__GetErrorLoggingFile__
```
ErrorLoggingFile(path string)
```
Returns the path to the log file.

__Generate32BytesKeys__
```
Generate32BytesKeys()
```
Generates 32-byte keys.

__Get32BytesKey__
```
Get32BytesKey() IKey
```
Returns an IKey structure responsible for key operations.