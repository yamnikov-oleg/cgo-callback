#include "entry.h"
#include "regs_386.h"

#include <stdlib.h>
#include <string.h>

typedef struct {
  // Number of integer (signed, unsigned and pointer) arguments used.
  int int_args;
  // Number of floating-point arguments used.
  int float_args;
} cgo_callback_sysv64_conv_t;

void cgo_callback_conv_init(cgo_callback_call_t *call) {
  cgo_callback_sysv64_conv_t *conv = malloc(sizeof(cgo_callback_sysv64_conv_t));
  conv->int_args = 0;
  conv->float_args = 0;
  call->conv = conv;
}

void cgo_callback_conv_destroy(cgo_callback_call_t *call) {
  free(call->conv);
}

int64_t cgo_callback_conv_get_int(void *addr, int bits) {
  signed char bytes[8] = {0};
  int bcount = bits/8;
  int i;
  for (i = 0; i < bcount; i++) {
    bytes[i] = *((char *)addr + i);
  }
  // For negative values must set leading ones.
  if (bytes[bcount-1] < 0) {
    for (i = bcount; i < 8; i++) {
      bytes[i] = -1;
    }
  }
  return *(int64_t *)bytes;
}


int64_t cgo_callback_conv_pop_int(cgo_callback_call_t* call, int bits) {
  int64_t val = cgo_callback_conv_get_int(call->sp, bits);
  char *csp = (char *)call->sp;
  csp += 4;
  call->sp = (void *)csp;
  return val;
}

uint64_t cgo_callback_conv_get_uint(void *addr, int bits) {
  char bytes[8] = {0};
  int bcount = bits/8;
  int i;
  for (i = 0; i < bcount; i++) {
    bytes[i] = *((char *)addr + i);
  }
  return *(uint64_t *)bytes;
}

uint64_t cgo_callback_conv_pop_uint(cgo_callback_call_t* call, int bits) {
  uint64_t val = cgo_callback_conv_get_uint(call->sp, bits);
  char *csp = (char *)call->sp;
  csp += 4;
  call->sp = (void *)csp;
  return val;
}

float cgo_callback_conv_get_single(void* addr) {
  return *(float *)((char *)addr);
}

float cgo_callback_conv_pop_single(cgo_callback_call_t* call) {
  float val = cgo_callback_conv_get_single((char *)call->sp);
  char *csp = (char *)call->sp;
  csp += 4;
  call->sp = (void *)csp;
  return val;
}

double cgo_callback_conv_get_double(void *addr) {
  return *(double *)((char *)addr);
}

double cgo_callback_conv_pop_double(cgo_callback_call_t* call) {
  double val = cgo_callback_conv_get_double((char *)call->sp);
  char *csp = (char *)call->sp;
  csp += 8;
  call->sp = (void *)csp;
  return val;
}

int64_t cgo_callback_conv_get_arg_int(cgo_callback_call_t *call, int bits) {
  return cgo_callback_conv_pop_int(call, bits);
}

uint64_t cgo_callback_conv_get_arg_uint(cgo_callback_call_t *call, int bits) {
  return cgo_callback_conv_pop_uint(call, bits);
}

float cgo_callback_conv_get_arg_single(cgo_callback_call_t *call) {
  return cgo_callback_conv_pop_single(call);
}

double cgo_callback_conv_get_arg_double(cgo_callback_call_t *call) {
  return cgo_callback_conv_pop_double(call);
}

void cgo_callback_conv_return(cgo_callback_call_t *call, void *val, int type, int bits) {
  memcpy((char *)call->reg + EAX, val, bits/8);
}
