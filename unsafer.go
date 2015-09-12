package unsafer

import (
	"debug/gosym"
	"github.com/bouk/symme"
	"github.com/bouk/symme/objfile"
	"sync"
)

var (
	table *gosym.Table
	syms  []objfile.Sym

	tableOnce sync.Once
)

func inittable() {
	tableOnce.Do(func() {
		var err error
		table, err = symme.Table()
		if err != nil {
			panic(err)
		}
		syms, err = symme.Symbols()
		if err != nil {
			panic(err)
		}
	})
}
