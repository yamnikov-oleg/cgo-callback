package callback

// #include "entry.h"
import "C"

import (
	"fmt"
	"reflect"
)

//export cgo_callback_go_entry
func cgo_callback_go_entry(call *C.cgo_callback_call_t) {
	ctx, ok := ctxs[uint(call.port)]
	if !ok {
		panic(fmt.Sprintf("cgo-callback: call to unused port %v", call.port))
	}

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
		}
	}
	ctx.fn.Call(args)
}
