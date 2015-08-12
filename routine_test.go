package unsafer

import (
	"fmt"
	"testing"
	"time"
)

func TestGetG(t *testing.T) {
	for i := 0; i < 1; i++ {
		go func() {
			fmt.Printf("%+v\n", GetG().M)
		}()
	}
	time.Sleep(time.Second)
}

func TestTracebackOthers(t *testing.T) {
	TracebackOthers()
}
