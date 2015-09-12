package unsafer

import (
	"testing"
)

func TestStop(t *testing.T) {
	success := true

	StopTheWorld("testing")

	go func() {
		for i := 0; i < 100000000; i++ {
			success = false
		}
	}()

	for j := 0; j < 100000000; j++ {
		if !success {
			t.Error("Other go routine woke up")
			return
		}
	}

	StartTheWorld()
}
