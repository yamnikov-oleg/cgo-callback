//+build ignore

package main

import (
	"fmt"
	"os"
)

const maxPorts = 128

func main() {
	portsDotS, err := os.Create("ports_x86.s")
	if err != nil {
		panic(err)
	}
	defer portsDotS.Close()

	portsDotS.Write([]byte(portsPreambule))
	for i := 0; i < maxPorts; i++ {
		portsDotS.Write([]byte("\tcall cgo_callback_asm_entry\n"))
	}

	portsmaxDotGo, err := os.Create("portsmax.go")
	if err != nil {
		panic(err)
	}
	defer portsmaxDotGo.Close()
	portsmaxDotGo.Write([]byte(fmt.Sprintf(portsmaxText, maxPorts)))
}

const portsPreambule = `//+build i386 amd64

// Generated using gen_ports.go.
// Do not edit.

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
`

const portsmaxText = `// Generated using gen_ports.go.
// Do not edit.

package callback

const maxPort = %v
`
