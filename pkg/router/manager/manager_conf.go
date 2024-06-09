package manager

type ManagerConf struct {
	debug            bool
	errorLogging     bool
	errorLoggingPath string
}

func NewManagerConf() *ManagerConf {
	return &ManagerConf{}
}

func (m *ManagerConf) Debug(enable bool) {
	m.debug = enable
}

func (m *ManagerConf) IsDebug() bool {
	return m.debug
}

func (m *ManagerConf) ErrorLogging(enable bool) {
	m.errorLogging = enable
}

func (m *ManagerConf) IsErrorLogging() bool {
	return m.errorLogging
}

func (m *ManagerConf) ErrorLoggingFile(path string) {
	m.errorLoggingPath = path
}

func (m *ManagerConf) GetErrorLoggingFile() string {
	return m.errorLoggingPath
}
