package callback

// #include "entry.h"
import "C"

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"unsafe"
)

//export cgo_callback_go_entry
func cgo_callback_go_entry(call *C.cgo_callback_call_t) {
	ctx, ok := ctxs[uint(call.port)]
	if !ok {
		panic(fmt.Sprintf("cgo-callback: call to unused port %v", call.port))
	}

	call.cleanstack = C._Bool(ctx.cleanstack)

	var args []reflect.Value
	for _, val := range ctx.ins {
		switch val.kind {
		case signed:
			s := C.cgo_callback_conv_get_arg_int(call, C.int(val.size*8))
			args = append(args, reflect.ValueOf(s).Convert(val.reft))
		case unsigned:
			u := C.cgo_callback_conv_get_arg_uint(call, C.int(val.size*8))
			args = append(args, reflect.ValueOf(u).Convert(val.reft))
		case singlePrec:
			f := C.cgo_callback_conv_get_arg_single(call)
			args = append(args, reflect.ValueOf(f).Convert(val.reft))
		case doublePrec:
			f := C.cgo_callback_conv_get_arg_double(call)
			args = append(args, reflect.ValueOf(f).Convert(val.reft))
		case pointer:
			u := C.cgo_callback_conv_get_arg_uint(call, C.int(val.size*8))
			ptr := unsafe.Pointer(uintptr(u))
			args = append(args, reflect.ValueOf(ptr).Convert(val.reft))
		}
	}

	rets := ctx.fn.Call(args)
	if ctx.ret == nil {
		return
	}

	var reti interface{}
	if ctx.ret.kind == pointer {
		// Pass the pointer to the dummy C function to check if it suits Cgo
		// pointer passing rules on go 1.6+.
		// Cgo will panic if there's something wrong.
		ptr := rets[0].Convert(reflect.TypeOf(unsafe.Pointer(nil))).Interface().(unsafe.Pointer)
		C.cgo_callback_assert_ptr(ptr)
		reti = uint64(uintptr(ptr))
	} else {
		reti = rets[0].Convert(ctx.ret.reft).Interface()
	}
	var arr [16]byte
	buf := bytes.NewBuffer(arr[0:0])
	if err := binary.Write(buf, binary.LittleEndian, reti); err != nil {
		panic("cgo-callback: " + err.Error())
	}
	C.cgo_callback_conv_return(call, unsafe.Pointer(&arr), ctx.ret.kind.toCType(), C.int(ctx.ret.size*8))
}
