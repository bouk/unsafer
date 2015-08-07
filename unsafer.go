package unsafer

import (
	"reflect"
	"unsafe"
)

var DMA []byte

func init() {
	DMA = *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: 0,
		Len:  (1 << 63) - 1,
		Cap:  (1 << 63) - 1,
	}))
}
