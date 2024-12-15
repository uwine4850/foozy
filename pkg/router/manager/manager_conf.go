package manager

import (
	"github.com/uwine4850/foozy/pkg/interfaces"
)

type DebugConfig struct {
	debug            bool
	relPath          bool
	errorLogging     bool
	errorLoggingPath string
	requestInfo      bool
	requestInfoPath  string
	skipLoggingLevel int
}

func (d *DebugConfig) Debug(enable bool) {
	d.debug = enable
}

func (d *DebugConfig) IsDebug() bool {
	return d.debug
}

func (d *DebugConfig) RelativeFilePath(enable bool) {
	d.relPath = enable
}

func (d *DebugConfig) IsRelativeFilePath() bool {
	return d.relPath
}

func (d *DebugConfig) ErrorLogging(enable bool) {
	d.errorLogging = enable
}

func (d *DebugConfig) IsErrorLogging() bool {
	return d.errorLogging
}

func (d *DebugConfig) ErrorLoggingFile(path string) {
	d.errorLoggingPath = path
}

func (d *DebugConfig) GetErrorLoggingFile() string {
	return d.errorLoggingPath
}

func (d *DebugConfig) RequestInfo(enable bool) {
	d.requestInfo = enable
}

func (d *DebugConfig) IsRequestInfo() bool {
	return d.requestInfo
}

func (d *DebugConfig) RequestInfoFile(path string) {
	d.requestInfoPath = path
}

func (d *DebugConfig) GetRequestInfoFile() string {
	return d.requestInfoPath
}

func (d *DebugConfig) SkipLoggingLevel(skip int) {
	d.skipLoggingLevel = skip
}

func (d *DebugConfig) LoggingLevel() int {
	return d.skipLoggingLevel
}

type ManagerCnf struct {
	debugConfig *DebugConfig
	printLog    bool
}

func NewManagerCnf() *ManagerCnf {
	return &ManagerCnf{
		debugConfig: &DebugConfig{
			skipLoggingLevel: -1,
		},
	}
}

func (m *ManagerCnf) DebugConfig() interfaces.IManagerDebugConfig {
	return m.debugConfig
}

func (m *ManagerCnf) PrintLog(enable bool) {
	m.printLog = enable
}

func (m *ManagerCnf) IsPrintLog() bool {
	return m.printLog
}
