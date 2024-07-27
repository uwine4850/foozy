package manager

import (
	"math/rand"
	"time"

	"github.com/uwine4850/foozy/pkg/interfaces"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Key structure that generates and stores three types of keys:
// hashKey is a key that is used for HMAC and can be dynamic.
// blockKey - a key that is used for encoding and can be dynamic.
// staticKey - a key that cannot change.
// The old keys haskKey and blockKey are also stored here.
type Key struct {
	hashKey     string
	oldHashKey  string
	blockKey    string
	oldBlockKey string
	staticKey   string
	date        time.Time
}

func (k *Key) HashKey() string {
	return k.hashKey
}

func (k *Key) OldHashKey() string {
	return k.oldHashKey
}

func (k *Key) BlockKey() string {
	return k.blockKey
}

func (k *Key) OldBlockKey() string {
	return k.oldBlockKey
}

func (k *Key) StaticKey() string {
	return k.staticKey
}

func (k *Key) Date() time.Time {
	return k.date
}

// GenerateBytesKeys generates keys.
// hashKey and blockKey will be updated. staticKey will only be generated once, cannot be regenerated.
func (k *Key) GenerateBytesKeys(length int) {
	k.oldHashKey = k.hashKey
	k.oldBlockKey = k.blockKey
	k.hashKey = string(k.generateKeys(length))
	k.blockKey = string(k.generateKeys(length))
	if k.staticKey == "" {
		k.staticKey = string(k.generateKeys(length))
	}
	k.date = time.Now()
}

func (k *Key) generateKeys(length int) []byte {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return b
}

type ManagerConf struct {
	debug            bool
	printLog         bool
	errorLogging     bool
	errorLoggingPath string
	skipLoggingLevel int
	key              Key
}

func NewManagerConf() *ManagerConf {
	return &ManagerConf{skipLoggingLevel: -1}
}

func (m *ManagerConf) Debug(enable bool) {
	m.debug = enable
}

func (m *ManagerConf) IsDebug() bool {
	return m.debug
}

func (m *ManagerConf) PrintLog(enable bool) {
	m.printLog = enable
}

func (m *ManagerConf) IsPrintLog() bool {
	return m.printLog
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

func (m *ManagerConf) SkipLiggingLevel(skip int) {
	m.skipLoggingLevel = skip
}

func (m *ManagerConf) LoggingLevel() int {
	return m.skipLoggingLevel
}

func (m *ManagerConf) Generate32BytesKeys() {
	m.key.GenerateBytesKeys(32)
}

func (m *ManagerConf) Get32BytesKey() interfaces.IKey {
	return &m.key
}
