## Debug
This package contains tools for debugging.

#### WriteLog
WriteLog writes the message to a log file.

`skipLevel` - skips levels of runtime.Caller.<br>
This is used to output the path on which the log is written.

`filePath` - path to the log file.

`flag` - flags for opening a file from the os package.

`prefix` - prefix that will be shown in the log. It is desirable to use
constants “P_...” from this package.

`message` - the message that will be recorded.
```golang
func WriteLog(skipLevel int, filePath string, flag int, prefix string, message string, logFlags int) {
	var _logFlags int
	if logFlags == -1 {
		_logFlags = log.LstdFlags
	}
	f, err := os.OpenFile(filePath, flag, 0666)
	if err != nil {
		fmt.Println("LogError: ", err.Error())
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println("LogError: ", err.Error())
			return
		}
	}(f)
	ilog := log.New(f, fmt.Sprintf("[%s] ", prefix), _logFlags)
	ilog.SetFlags(_logFlags)

	if skipLevel < 0 {
		skipLevel = 3
	}

	_, file, line, ok := runtime.Caller(skipLevel)
	if !ok {
		ilog.Println("Could not retrieve caller information")
		return
	}
	loggingFilePath := file
	if config.LoadedConfig().Default.Debug.DebugRelativeFilepath {
		wd, err := os.Getwd()
		if err != nil {
			ilog.Println("Could not retrieve working directory:", err)
			return
		}
		relPath, err := filepath.Rel(wd, file)
		if err != nil {
			ilog.Println("Could not calculate relative path:", err)
			return
		}
		loggingFilePath = relPath
	}
	ilog.Printf("%s:%d %s\n", loggingFilePath, line, message)
}
```

#### LogError
Logs errors to the error log.
```golang
func LogError(message string) {
	if config.LoadedConfig().Default.Debug.ErrorLoggingPath == "" {
		panic("unable to create log file. File path not set")
	}
	WriteLog(
		config.LoadedConfig().Default.Debug.SkipLoggingLevel+1,
		config.LoadedConfig().Default.Debug.ErrorLoggingPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		"",
		message, -1,
	)
}
```

#### ErrorLoggingIfEnableAndWrite
Writes a message to the log if error logging is enabled.<br>
This function also writes a message to the browser page. It is
convenient for displaying the error on the page.
```golang
func ErrorLoggingIfEnableAndWrite(w http.ResponseWriter, errorText string, writeText string) {
	_, err := w.Write([]byte(writeText))
	if err != nil {
		if config.LoadedConfig().Default.Debug.ErrorLogging {
			LogError(err.Error())
		}
	}
	if config.LoadedConfig().Default.Debug.ErrorLogging {
		LogError(string(errorText))
	}
}
```

#### ErrorLogginIfEnable
Writes a message to the log if error logging is enabled.
```golang
func ErrorLogginIfEnable(message string) {
	if config.LoadedConfig().Default.Debug.ErrorLogging {
		LogError(message)
	}
}
```

#### ClearRequestInfoLogging
Clears the request log.
```golang
func ClearRequestInfoLogging() error {
	if config.LoadedConfig().Default.Debug.RequestInfoLogPath != "" && fpath.PathExist(config.LoadedConfig().Default.Debug.RequestInfoLogPath) {
		f, err := os.OpenFile(config.LoadedConfig().Default.Debug.RequestInfoLogPath, os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		if err := f.Truncate(0); err != nil {
			return err
		}
	}
	return nil
}
```

#### LogRequestInfo
Logs information about the request.
```golang
func LogRequestInfo(prefix string, message string) {
	if config.LoadedConfig().Default.Debug.RequestInfoLogPath == "" {
		panic("unable to create request info log file. File path not set")
	}
	WriteLog(
		config.LoadedConfig().Default.Debug.SkipLoggingLevel,
		config.LoadedConfig().Default.Debug.RequestInfoLogPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		prefix,
		message,
		-1,
	)
}
```

#### RequestLogginIfEnable
Logs request information if request logging is enabled.
```golang
func RequestLogginIfEnable(prefix string, message string) {
	if config.LoadedConfig().Default.Debug.RequestInfoLog {
		LogRequestInfo(prefix, message)
	}
}
```