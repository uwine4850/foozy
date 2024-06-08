package debug

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/uwine4850/foozy/pkg/interfaces"
)

func LogError(message string, manager interfaces.IManagerConfig) {
	if manager.GetErrorLoggingFile() == "" {
		panic("Unable to create log file. File path not set")
	}
	f, err := os.OpenFile(manager.GetErrorLoggingFile(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	ilog := log.New(f, "", log.LstdFlags)
	ilog.SetFlags(log.LstdFlags)

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		ilog.Println("Could not retrieve caller information")
		return
	}
	ilog.Printf("%s:%d %s\n", file, line, message)
}

func ErrorLoggingIfEnableAndWrite(w http.ResponseWriter, text []byte, manager interfaces.IManagerConfig) {
	_, err := w.Write(text)
	if err != nil {
		fmt.Println("LoggingIfEnableAndWrite: ", err.Error())
		if manager.IsErrorLogging() {
			LogError(err.Error(), manager)
		}
	}
	if manager.IsErrorLogging() {
		LogError(string(text), manager)
	}
}

func ErrorLogginIfEnable(message string, manager interfaces.IManagerConfig) {
	if manager.IsErrorLogging() {
		LogError(message, manager)
	}
}
