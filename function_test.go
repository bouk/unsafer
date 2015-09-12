package unsafer_test

import (
	"github.com/bouk/unsafer"
	"testing"
)

func tryAndFindMe() bool {
	return true
}

func TestInsertFunction(t *testing.T) {
	var rightHere func() bool = tryAndFindMe
	rightHere = nil
	if err := unsafer.InsertFunction("github.com/bouk/unsafer_test.tryAndFindMe", &rightHere); err != nil {
		t.Error(err)
	}
	if !rightHere() {
		t.Error("didn't call tryAndFindMe")
	}
}
