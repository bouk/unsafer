package unsafer

import (
	"reflect"
	"unsafe"
)

const (
	_KindBool = 1 + iota
	_KindInt
	_KindInt8
	_KindInt16
	_KindInt32
	_KindInt64
	_KindUint
	_KindUint8
	_KindUint16
	_KindUint32
	_KindUint64
	_KindUintptr
	_KindFloat32
	_KindFloat64
	_KindComplex64
	_KindComplex128

	_KindArray
	_KindChan
	_KindFunc
	_KindInterface
	_KindMap
	_KindPtr
	_KindSlice
	_KindString
	_KindStruct

	_KindUnsafePointer

	_KindDirectIface = 1 << 5
	_KindGCProg      = 1 << 6 // Type.gc points to GC program
	_KindNoPointers  = 1 << 7
	_KindMask        = (1 << 5) - 1
)

// runtime._type
type Type struct {
	Size       uintptr
	PtrData    uintptr // size of memory prefix holding all pointers
	Hash       uint32
	_unused    uint8
	Align      uint8
	FieldAlign uint8
	Kind       uint8
	Alg        uintptr
	// gcdata stores the GC type data for the garbage collector.
	// If the KindGCProg bit is set in kind, gcdata is a GC program.
	// Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
	GCData *byte
	String *string
	X      *UncommonType
	PtrTo  *Type
	Zero   *byte // ptr to the zero value for this type
}

func (t *Type) Reflect() reflect.Type {
	var i interface{} = t
	*(*uintptr)(unsafe.Pointer(&i)) = uintptr(unsafe.Pointer(t))
	return reflect.TypeOf(i)
}

func (t *Type) ToSpecificType() SpecificType {
	switch t.Kind & _KindMask {
	case _KindInterface:
		return (*InterfaceType)(unsafe.Pointer(t))
	case _KindMap:
		return (*MapType)(unsafe.Pointer(t))
	case _KindChan:
		return (*ChanType)(unsafe.Pointer(t))
	case _KindFunc:
		return (*FuncType)(unsafe.Pointer(t))
	case _KindPtr:
		return (*PtrType)(unsafe.Pointer(t))
	case _KindSlice:
		return (*SliceType)(unsafe.Pointer(t))
	case _KindArray:
		return (*ArrayType)(unsafe.Pointer(t))
	case _KindStruct:
		return (*StructType)(unsafe.Pointer(t))
	default:
		return (*GenericType)(unsafe.Pointer(t))
	}
}

type SpecificType interface {
	ToType() *Type
	Sizeof() uintptr
}

func (t *GenericType) ToType() *Type   { return &t.Type }
func (t *InterfaceType) ToType() *Type { return &t.Type }
func (t *MapType) ToType() *Type       { return &t.Type }
func (t *ChanType) ToType() *Type      { return &t.Type }
func (t *FuncType) ToType() *Type      { return &t.Type }
func (t *PtrType) ToType() *Type       { return &t.Type }
func (t *SliceType) ToType() *Type     { return &t.Type }
func (t *StructType) ToType() *Type    { return &t.Type }
func (t *ArrayType) ToType() *Type     { return &t.Type }

func (t *GenericType) Sizeof() uintptr { return unsafe.Sizeof(*t) }
func (t *InterfaceType) Sizeof() uintptr {
	return unsafe.Sizeof(*t) + unsafe.Sizeof(t.Mhdr[0])*uintptr(len(t.Mhdr))
}
func (t *MapType) Sizeof() uintptr  { return unsafe.Sizeof(*t) }
func (t *ChanType) Sizeof() uintptr { return unsafe.Sizeof(*t) }
func (t *FuncType) Sizeof() uintptr {
	return unsafe.Sizeof(*t) + unsafe.Sizeof(t.In[0])*uintptr(len(t.In)) + unsafe.Sizeof(t.Out[0])*uintptr(len(t.Out))
}
func (t *PtrType) Sizeof() uintptr   { return unsafe.Sizeof(*t) }
func (t *SliceType) Sizeof() uintptr { return unsafe.Sizeof(*t) }
func (t *StructType) Sizeof() uintptr {
	return unsafe.Sizeof(*t) + unsafe.Sizeof(t.Fields[0])*uintptr(len(t.Fields))
}
func (t *ArrayType) Sizeof() uintptr { return unsafe.Sizeof(*t) }

type Method struct {
	Name    *string
	PkgPath *string
	MType   *Type
	Type    *Type
	Ifn     unsafe.Pointer
	Tfn     unsafe.Pointer
}

type UncommonType struct {
	Name    *string
	PkgPath *string
	Mhdr    []Method
}

type IMethod struct {
	Name    *string
	PkgPath *string
	Type    *Type
}

type GenericType struct {
	Type
}

type InterfaceType struct {
	Type
	Mhdr []IMethod
}

type MapType struct {
	Type
	Key           *Type
	Element       *Type
	bucket        *Type  // internal type representing a hash bucket
	hmap          *Type  // internal type representing a hmap
	keysize       uint8  // size of key slot
	indirectkey   bool   // store ptr to key instead of key itself
	valuesize     uint8  // size of value slot
	indirectvalue bool   // store ptr to value instead of value itself
	bucketsize    uint16 // size of bucket
	reflexivekey  bool   // true if k==k for all keys
}

type ArrayType struct {
	Type
	Element *Type
	Slice   *Type
	Length  uintptr
}

type ChanType struct {
	Type
	Element   *Type
	Direction uintptr
}

type SliceType struct {
	Type
	Element *Type
}

type FuncType struct {
	Type
	DotDotDot bool
	In        []*Type
	Out       []*Type
}

type PtrType struct {
	Type
	Element *Type
}

type StructField struct {
	Name    *string
	PkgPath *string
	Type    *Type
	Tag     *string
	Offset  uintptr
}

type StructType struct {
	Type
	Fields []StructField
}

type eface struct {
	typ *Type
	val unsafe.Pointer
}

func itoe(i *interface{}) *eface {
	return (*eface)(unsafe.Pointer(i))
}
