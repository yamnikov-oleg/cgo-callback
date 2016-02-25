// Interface of the C entry point provided for the Go entry point.
// It declares functions to be used to retrieve arguments, set return values,
// and to translate port numbers and their addresses.

#pragma once

#include <stdint.h>
#include <stdbool.h>

#define BITS_8 8
#define BITS_16 16
#define BITS_32 32
#define BITS_64 64

// #define TYPE_INT 0
// #define TYPE_FLOAT 1

// Calculate address of the port by id and vice versa.
void *cgo_callback_get_port_addr(unsigned id);
unsigned cgo_callback_get_port_id(void *addr);

// Structure containing every information needed by C code to work with a call.
typedef struct {
  // Number of port.
  unsigned port;
  // Pointer to the top of arguments stack. Popping arguments updates this value.
  void *sp;
  // Should the callee clean up the stack? False by default.
  // Ignored on conventions, which don't have callee-cleanup version.
  bool clean_stack;
  // Number of bytes popped of the stack. Used by "callee-cleanup" conventions.
  int popped_stack;
  // Pointer to the register map.
  void *reg;
  // Pointer to the convention-specific data.
  // It's used convention-specific functions.
  void *conv;
} cgo_callback_call_t;

// Get next argument of the specified type and size.
// Each platform must have its implementation of these functions.
int64_t cgo_callback_conv_get_arg_int(cgo_callback_call_t *call, int bits);
uint64_t cgo_callback_conv_get_arg_uint(cgo_callback_call_t *call, int bits);
// double cgo_callback_conv_get_arg_float(cgo_callback_call_t *call, int bits);

// Set the return value of the specified type and size.
// Each platform must have its implementation of this function.
// void cgo_callback_conv_return(cgo_callback_call_t *call, void *val, int type, int bits);
