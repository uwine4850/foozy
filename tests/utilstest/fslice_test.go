package utilstest

import (
	"testing"

	"github.com/uwine4850/foozy/pkg/utils/fslice"
)

func TestContains(t *testing.T) {
	slice := []string{"1", "2", "3"}
	if !fslice.SliceContains(slice, "2") {
		t.Error("The check for the presence of the element was not successful.")
	}
	if fslice.SliceContains(slice, "44") {
		t.Error("Checking for the presence of an element should return false.")
	}
}

func TestAllStringItemsEmpty(t *testing.T) {
	slice := []string{"", "", ""}
	if !fslice.AllStringItemsEmpty(slice) {
		t.Error("The check to see if all slice string elements are empty failed.")
	}
	slice1 := []string{"", "2", ""}
	if fslice.AllStringItemsEmpty(slice1) {
		t.Error("Checking for all slice string elements to be empty should be false.")
	}
}
