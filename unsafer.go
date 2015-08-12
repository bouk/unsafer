package unsafer

import (
	"debug/gosym"
	"github.com/bouk/symme"
	"sync"
)

var table *gosym.Table
var tableOnce sync.Once

func inittable() {
	tableOnce.Do(func() {
		var err error
		table, err = symme.Table()
		if err != nil {
			panic(err)
		}
	})
}
