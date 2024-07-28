## IManagerConfig
Saves and provides access to some settings.

__DebugConfig()__
```
DebugConfig() IManagerDebugConfig
```
Provides access to the dabug configuration.

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

__Key__
```
Key() IManagerDebugConfig
```
Provides access to the configuration of the key object.

## IManagerDebugConfig

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

__SkipLoggingLevel__
```
SkipLoggingLevel(skip int)
```
Setting the logging level. Default is 3.
This parameter shows the location of the logging function call in the project. 
For example, the number 0 will show the location of the logging function call in the framework implementation.

__LoggingLevel__
```
LoggingLevel() int
```
Returns the value of the logging level.

## IKey

__HashKey__
```
HashKey() string
```
Generates a hash key.

__OldHashKey__
```
OldHashKey() string
```
Returns the old login key.
This method is used after the hash key is generated again.
It is important to note that this method returns the previous key, starting with the active one.

__BlockKey__
```
BlockKey() string
```
Generates a block key.

__OldBlockKey__
```
OldBlockKey() string
```
Returns the old login key.
This method is used after the block key is generated again.
It is important to note that this method returns the previous key, starting with the active one.

__StaticKey()__
```
StaticKey() string
```
Creates a static key. This key is generated once when the server is started and does not change.

__Date() time.Time__
```
Date() time.Time
```
Returns the time of the last key generation.

__GenerateBytesKeys__
```
GenerateBytesKeys(length int)
```
Generates block key and hash key

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