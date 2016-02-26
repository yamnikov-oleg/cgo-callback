package callback

// #include "entry.h"
import "C"
import (
	"fmt"
	"reflect"
)

type kind int

const (
	invalid  kind = iota
	signed        // Signed integers
	unsigned      // Unsigned integers and pointers
	singlePrec    // Single-precision float
	doublePrec    // Double-precision float
)

func kindFromReflect(k reflect.Kind) kind {
	switch k {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return signed

	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Ptr,
		reflect.UnsafePointer:
		return unsigned

	case reflect.Float32:
		return singlePrec

	case reflect.Float64:
		return doublePrec

	default:
		return invalid
	}
}

type value struct {
	kind kind
	size uintptr
	reft reflect.Type
}

func valueFromReflect(t reflect.Type) (v value) {
	v.kind = kindFromReflect(t.Kind())
	v.size = t.Size()
	v.reft = t
	return
}

type context struct {
	port uint
	fn   reflect.Value
	ins  []value
	ret  *value
}

var ctxs = map[uint]*context{}

func findNextPort() uint {
	for i := uint(0); i < maxPort; i++ {
		if _, ok := ctxs[i]; !ok {
			return i
		}
	}
	panic("cgo-callback: ran out of ports")
}

func New(f interface{}) uintptr {
	ctx := context{}

	ctx.fn = reflect.ValueOf(f)
	ft := ctx.fn.Type()
	if ft.Kind() != reflect.Func {
		panic("cgo-callback: New() only accepts functions")
	}

	for i := 0; i < ft.NumIn(); i++ {
		val := valueFromReflect(ft.In(i))
		if val.kind == invalid {
			panic(fmt.Sprintf("cgo-callback: callbacks doesn't support arguments of kind %v", ft.In(i).Kind()))
		}
		ctx.ins = append(ctx.ins, val)
	}

	if ft.NumOut() > 0 {
		val := valueFromReflect(ft.Out(0))
		if val.kind == invalid {
			panic(fmt.Sprintf("cgo-callback: callbacks doesn't support return values of kind %v", ft.Out(0).Kind()))
		}
		ctx.ret = &val
	}

	ctx.port = findNextPort()
	ctxs[ctx.port] = &ctx
	return uintptr(C.cgo_callback_get_port_addr(C.uint(ctx.port)))
}
