package unsafer

import (
	"unsafe"
)

type Defer struct {
	Size    int32
	Started bool
	SP      uintptr // sp at time of defer
	PC      uintptr
	Fn      *FuncVal
	Panic   *Panic // panic that is running defer
	Link    *Defer
}

/*
 * panics
 */
type Panic struct {
	Argp      unsafe.Pointer // pointer to arguments of deferred call run during panic; cannot move - known to liblink
	Arg       interface{}    // argument to panic
	Link      *Panic         // link to earlier panic
	Recovered bool           // whether this panic is over
	Aborted   bool           // the panic was aborted
}

type FuncVal struct {
	PC uintptr
}
