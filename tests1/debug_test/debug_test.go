package debug_test

import (
	"os"
	"testing"
	"time"

	"github.com/uwine4850/foozy/pkg/debug"
	initcnf_t "github.com/uwine4850/foozy/tests1/common/init_cnf"
	"github.com/uwine4850/foozy/tests1/common/tutils"
)

func TestWriteLog(t *testing.T) {
	initcnf_t.InitCnf()
	debug.WriteLog(1, "write_log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, "TEST", "TEE", 0)
	isOk, err := tutils.FilesAreEqual("write_log.log", "exp_write_log.log")
	if err != nil {
		t.Error(err)
	}
	time.Sleep(10 * time.Millisecond)
	if !isOk {
		t.Error("passed and expected logs do not match")
	}
	os.Remove("write_log.log")
}
