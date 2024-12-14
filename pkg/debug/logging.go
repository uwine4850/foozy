package debug

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/uwine4850/foozy/pkg/interfaces"
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

func WriteLog(skipLevel int, filePath string, flag int, prefix string, message string, managerConfig interfaces.IManagerConfig) {
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
	if managerConfig.DebugConfig().IsRelativeFilePath() {
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

func LogError(message string, managerConfig interfaces.IManagerConfig) {
	if managerConfig.DebugConfig().GetErrorLoggingFile() == "" {
		panic("unable to create log file. File path not set")
	}
	WriteLog(managerConfig.DebugConfig().LoggingLevel()+1, managerConfig.DebugConfig().GetErrorLoggingFile(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, "", message, managerConfig)
}

func ErrorLoggingIfEnableAndWrite(w http.ResponseWriter, errorText string, writeText string, managerConfig interfaces.IManagerConfig) {
	_, err := w.Write([]byte(writeText))
	if err != nil {
		if managerConfig.DebugConfig().IsErrorLogging() {
			LogError(err.Error(), managerConfig)
		}
	}
	if managerConfig.DebugConfig().IsErrorLogging() {
		LogError(string(errorText), managerConfig)
	}
}

func ErrorLogginIfEnable(message string, managerConfig interfaces.IManagerConfig) {
	if managerConfig.DebugConfig().IsErrorLogging() {
		LogError(message, managerConfig)
	}
}

func ClearRequestInfoLogging(managerConfig interfaces.IManagerConfig) error {
	if managerConfig.DebugConfig().GetRequestInfoFile() != "" && fstring.PathExist(managerConfig.DebugConfig().GetRequestInfoFile()) {
		f, err := os.OpenFile(managerConfig.DebugConfig().GetRequestInfoFile(), os.O_WRONLY, 0644)
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

func LogRequestInfo(prefix string, message string, managerConfig interfaces.IManagerConfig) {
	if managerConfig.DebugConfig().GetRequestInfoFile() == "" {
		panic("unable to create request info log file. File path not set")
	}
	WriteLog(managerConfig.DebugConfig().LoggingLevel(), managerConfig.DebugConfig().GetRequestInfoFile(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, prefix, message, managerConfig)
}

func RequestLogginIfEnable(prefix string, message string, managerConfig interfaces.IManagerConfig) {
	if managerConfig.DebugConfig().IsRequestInfo() {
		LogRequestInfo(prefix, message, managerConfig)
	}
}
