package manager

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type ManagerConf struct {
	debug            bool
	errorLogging     bool
	errorLoggingPath string
	hashKey          string
	blockKey         string
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

func (m *ManagerConf) Generate32BytesKeys() {
	m.hashKey = generateRandomBytesString(32)
	m.blockKey = generateRandomBytesString(32)
}

func (m *ManagerConf) Get32BytesKeys() map[string]string {
	mp := make(map[string]string, 2)
	mp["HashKey"] = m.hashKey
	mp["BlockKey"] = m.blockKey
	return mp
}

func generateRandomBytesString(length int) string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
