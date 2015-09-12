package unsafer

import (
	"testing"
	"unsafe"
)

func TestToReflectType(t *testing.T) {
	var min, max unsafe.Pointer

	for m := FirstModuleData; m != nil; m = m.Next {
		for _, typ := range m.TypeLinks {
			if uintptr(min) == 0 || uintptr(unsafe.Pointer(&typ)) < uintptr(min) {
				min = unsafe.Pointer(&typ)
			}
			if uintptr(max) == 0 || uintptr(unsafe.Pointer(&typ)) > uintptr(max) {
				max = unsafe.Pointer(&typ)
			}
		}
	}
}
