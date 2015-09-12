package unsafer_test

import (
	"github.com/bouk/unsafer"
	"testing"
)

func deferMe(called *bool) {
	*called = true
}

func tryMe() (called bool) {
	defer deferMe(&called)

	if g := unsafer.GetG(); g.Defer != nil {
		g.Defer = g.Defer.Link
	}

	return
}

func TestDefer(t *testing.T) {
	called := tryMe()

	if called {
		t.Error("deferred function was called")
	}
}
