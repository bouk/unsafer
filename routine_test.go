package unsafer_test

import (
	"fmt"
	"github.com/bouk/symme"
	"github.com/bouk/unsafer"
	"testing"
)

func TestAllgs(t *testing.T) {
	table, err := symme.Table()
	if err != nil {
		t.Fatal(err)
	}

	unsafer.AllGs(func(g *unsafer.Goroutine) {
		fileGo, lineGo, _ := table.PCToLine(uint64(g.Gopc))
		fileF, lineF, _ := table.PCToLine(uint64(g.Startpc))
		fmt.Printf("%s:%d ->\n\t%s:%d\n", fileGo, lineGo, fileF, lineF)
	})
}
