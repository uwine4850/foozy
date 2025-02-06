package utilstest

import (
	"reflect"
	"testing"

	"github.com/uwine4850/foozy/pkg/utils/fstring"
)

func TestSplitUrl(t *testing.T) {
	url := "/test/tee1/tee2"
	splitUrl := fstring.SplitUrl(url)
	if !reflect.DeepEqual(splitUrl, []string{"test", "tee1", "tee2"}) {
		t.Error("The URL is not separated correctly.")
	}
}
