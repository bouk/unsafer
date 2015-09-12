package unsafer

import (
	"unsafe"
)

// runtime.moduledata
type ModuleData struct {
	PclnData     []byte
	Ftab         []FuncTab
	Filetab      []uint32
	FindFuncTab  uintptr
	MinPC, MaxPC uintptr

	Text, EText           uintptr
	NoPtrData, ENoPtrData uintptr
	Data, EData           uintptr
	Bss, EBss             uintptr
	NoPtrBss, ENoPtrBss   uintptr
	End, GCData, GCBss    uintptr

	TypeLinks []*Type

	Name   string
	Hashes []ModuleHash

	GcDataMask, GcBssMask BitVector

	Next *ModuleData
}

func alignType(n uintptr) uintptr {
	if n&0x1f == 0 {
		return n
	} else {
		return n&(^uintptr(0x1f)) + (0x20)
	}
}

func nextType(t *Type) *Type {
	p := uintptr(unsafe.Pointer(t)) + t.ToSpecificType().Sizeof()
	if t.X != nil {
		p += unsafe.Sizeof(*t.X) + unsafe.Sizeof(t.X.Mhdr[0])*uintptr(len(t.X.Mhdr))
	}
	return (*Type)(unsafe.Pointer(alignType(p)))
}

func (m *ModuleData) AllTypes(f func(*Type)) {
	var max, min uintptr = 0, ((1 << 64) - 1)
	for _, t := range m.TypeLinks {
		p := uintptr(unsafe.Pointer(t))
		if p < min {
			min = p
		}
		if p > max {
			max = p
		}
	}

	for t := (*Type)(unsafe.Pointer(min)); ; t = nextType(t) {
		f(t)
	}
}

type FuncTab struct {
	Entry      uintptr
	FuncOffset uintptr
}

type ModuleHash struct {
	ModuleName   string
	LinktimeHash string
	RuntimeHash  *string
}

type BitVector struct {
	N        int32 // # of bits
	ByteData *uint8
}

var (
	FirstModuleData, LastModuleData *ModuleData
)

func init() {
	inittable()
	for _, sym := range syms {
		if sym.Name == "runtime.firstmoduledata" {
			FirstModuleData = (*ModuleData)(unsafe.Pointer(uintptr(sym.Addr)))
		} else if sym.Name == "runtime.lastmoduledata" {
			LastModuleData = *(**ModuleData)(unsafe.Pointer(uintptr(sym.Addr)))
		}
	}
}
