package tests

// #include "calls.h"
import "C"
import (
	"unsafe"

	"github.com/yamnikov-oleg/cgo-callback"
)

type cint C.int
type cuint C.uint

func Void_Void(f func()) {
	ptr := callback.New(f)
	C.void_void(unsafe.Pointer(ptr))
}

func Void_Int(f func(cint), arg1 cint) {
	ptr := callback.New(f)
	C.void_int(unsafe.Pointer(ptr), C.int(arg1))
}

func Void_Uint(f func(cuint), arg1 cuint) {
	ptr := callback.New(f)
	C.void_uint(unsafe.Pointer(ptr), C.uint(arg1))
}

func Void_IntInt(f func(cint, cint), arg1 cint, arg2 cint) {
	ptr := callback.New(f)
	C.void_int_int(unsafe.Pointer(ptr), C.int(arg1), C.int(arg2))
}

func Void_Float(f func(float32), arg1 float32) {
	ptr := callback.New(f)
	C.void_float(unsafe.Pointer(ptr), C.float(arg1))
}

func Void_Double(f func(float64), arg1 float64) {
	ptr := callback.New(f)
	C.void_double(unsafe.Pointer(ptr), C.double(arg1))
}
