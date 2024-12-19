package debug

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/uwine4850/foozy/pkg/config"
	"github.com/uwine4850/foozy/pkg/utils/fstring"
)

const (
	P_ERROR           = "ERROR"
	P_ROUTER          = "ROUTER"
	P_MIDDLEWARE      = "MIDDLEWARE"
	P_OBJECT          = "OBJECT"
	P_REQUEST         = "REQUEST"
	P_TEMPLATE_ENGINE = "TEMPLATE_ENGINE"
	P_DATABASE        = "DATABASE"
)

func WriteLog(skipLevel int, filePath string, flag int, prefix string, message string) {
	f, err := os.OpenFile(filePath, flag, 0644)
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
	ilog := log.New(f, fmt.Sprintf("[%s] ", prefix), log.LstdFlags)
	ilog.SetFlags(log.LstdFlags)

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

func LogError(message string) {
	if config.LoadedConfig().Default.Debug.ErrorLoggingPath == "" {
		panic("unable to create log file. File path not set")
	}
	WriteLog(config.LoadedConfig().Default.Debug.SkipLoggingLevel+1, config.LoadedConfig().Default.Debug.ErrorLoggingPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, "", message)
}

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

func ErrorLogginIfEnable(message string) {
	if config.LoadedConfig().Default.Debug.ErrorLogging {
		LogError(message)
	}
}

func ClearRequestInfoLogging() error {
	if config.LoadedConfig().Default.Debug.RequestInfoLogPath != "" && fstring.PathExist(config.LoadedConfig().Default.Debug.RequestInfoLogPath) {
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

func LogRequestInfo(prefix string, message string) {
	if config.LoadedConfig().Default.Debug.RequestInfoLogPath == "" {
		panic("unable to create request info log file. File path not set")
	}
	WriteLog(config.LoadedConfig().Default.Debug.SkipLoggingLevel, config.LoadedConfig().Default.Debug.RequestInfoLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, prefix, message)
}

func RequestLogginIfEnable(prefix string, message string) {
	if config.LoadedConfig().Default.Debug.RequestInfoLog {
		LogRequestInfo(prefix, message)
	}
}
