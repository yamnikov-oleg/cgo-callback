package callback

// #include "entry.h"
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"
)

//export cgo_callback_go_entry
func cgo_callback_go_entry(call *C.cgo_callback_call_t) {
	// Get context
	// Get go func
	// Read argument types
	//   Map to raw type
	//   Get arg by raw type from C layer
	// Call go func
	// Map return type to raw type
	// Set return value

	// C.cgo_callback_get_arg_uint(call, C.BITS_16)
	// C.cgo_callback_get_arg_float(call, C.BITS_32)
	// C.cgo_callback_return(call, &val, C.TYPE_INT, C.BITS_16)

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
			// case float:
			// 	f := C.cgo_callback_conv_get_arg_float(call, C.int(val.size*8))
			// 	args = append(args, reflect.Value(f))
		}
	}
	ctx.fn.Call(args)
	// fmt.Printf("Call from port %d\n", call.port)
	// fmt.Printf("Arg1: %d\n", C.cgo_callback_conv_get_arg_int(call, C.BITS_32))
	// fmt.Printf("Arg2: %d\n", C.cgo_callback_conv_get_arg_uint(call, C.BITS_32))
	// fmt.Printf("Arg3: %d\n", C.cgo_callback_conv_get_arg_uint(call, C.BITS_32))
	// fmt.Printf("Arg4: %d\n", C.cgo_callback_conv_get_arg_uint(call, C.BITS_32))
	// fmt.Printf("Arg5: %d\n", C.cgo_callback_conv_get_arg_uint(call, C.BITS_32))
	// fmt.Printf("Arg6: %d\n", C.cgo_callback_conv_get_arg_uint(call, C.BITS_32))
	// fmt.Printf("Arg7: %d\n", C.cgo_callback_conv_get_arg_uint(call, C.BITS_32))
}

func GetAddr() unsafe.Pointer {
	return C.cgo_callback_get_port_addr(3)
}
