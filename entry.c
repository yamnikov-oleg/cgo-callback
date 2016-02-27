#include "entry.h"

// Linked to the first top address from the ports.s
extern char cgo_callback_ports;

// Functions implemented by the accepted calling convention.
void cgo_callback_conv_init(cgo_callback_call_t *call);
void cgo_callback_conv_destroy(cgo_callback_call_t *call);

void cgo_callback_assert_ptr(void *ptr){}

void *cgo_callback_get_port_addr(unsigned id) {
  return &cgo_callback_ports + id*5;
}

unsigned cgo_callback_get_port_id(void *addr) {
  return ((char *)addr - &cgo_callback_ports)/5;
}

// Second, "C" entry point.
// Its responsibility is to provide Go with an interface of reading arguments
// and returning values according the accepted calling convention.
//
// It's forced to use System V ABI to be cross platform.
__attribute__((sysv_abi))
void cgo_callback_c_entry(void *stack, void *reg) {
  cgo_callback_call_t call;

  // Pop the port address of the stack and discard return address.
  void **ptr_stack = (void **)stack;
  void *port_addr = *ptr_stack;
  ptr_stack += 2;
  stack = (void *)ptr_stack;

  call.port = cgo_callback_get_port_id(port_addr) - 1;
  call.sp = stack;
  call.clean_stack = false;
  call.popped_stack = 0;
  call.reg = reg;

  cgo_callback_conv_init(&call);
  cgo_callback_go_entry(&call);
  cgo_callback_conv_destroy(&call);
}
