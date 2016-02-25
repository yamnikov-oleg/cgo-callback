//+build i386 amd64

.text
.global cgo_callback_ports
.type cgo_callback_ports, @function

// Ports list, list of first entry points for incoming callbacks.
// Each "port" (CALL instruction) simply forwards the call to next entry point.
// Effectively it means, that it pushes next port's address onto the stack
// before calling to the real receiver. This address is later popped and used
// to determine number of the port and retrieve associated go function
// by that number. Vice versa, each go function gets associated with
// a free port to become a callback.
//
// I've borrowed this trick from syscall.NewCallback() implementation
// on windows.
cgo_callback_ports:
	call cgo_callback_asm_entry
	call cgo_callback_asm_entry
	call cgo_callback_asm_entry
	call cgo_callback_asm_entry
	call cgo_callback_asm_entry
