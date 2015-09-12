package unsafer

import (
	"fmt"
	"unsafe"
)

func InsertFunction(name string, receiver interface{}) error {
	inittable()
	f := table.LookupFunc(name)
	if f == nil {
		return fmt.Errorf("couldn't find function '%s'", name)
	}
	location := uintptr(f.Value)
	*(*uintptr)(itoe(&receiver).val) = uintptr(unsafe.Pointer(&location))
	return nil
}
