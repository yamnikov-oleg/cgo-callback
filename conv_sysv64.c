//+build linux,amd64

#include "entry.h"
#include "regs_amd64.h"

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

int64_t cgo_callback_conv_get_reg_int(cgo_callback_call_t* call, int reg, int bits) {
  signed char bytes[8] = {0};
  int bcount = bits/8;
  int i;
  for (i = 0; i < bcount; i++) {
    bytes[i] = *((char *)call->reg + reg + i);
  }
  // For negative values must set leading ones.
  if (bytes[bcount-1] < 0) {
    for (i = bcount; i < 8; i++) {
      bytes[i] = -1;
    }
  }
  return *(int64_t *)bytes;
}

uint64_t cgo_callback_conv_get_reg_uint(cgo_callback_call_t* call, int reg, int bits) {
  char bytes[8] = {0};
  int bcount = bits/8;
  int i;
  for (i = 0; i < bcount; i++) {
    bytes[i] = *((char *)call->reg + reg + i);
  }
  return *(uint64_t *)bytes;
}

float cgo_callback_conv_get_reg_single(cgo_callback_call_t* call, int reg) {
  return *(float *)((char *)call->reg + reg);
}

double cgo_callback_conv_get_reg_double(cgo_callback_call_t* call, int reg) {
  return *(double *)((char *)call->reg + reg);
}

// According to System V ABI, integer args are passed through
// RDI, RSI, RDX, RCX, R8, R9.
static int int_regs[6] = {RDI, RSI, RDX, RCX, R8, R9};

int64_t cgo_callback_conv_get_arg_int(cgo_callback_call_t *call, int bits) {
  cgo_callback_sysv64_conv_t *conv = call->conv;

  // TODO: support more than 6 integer arguments.
  if (conv->int_args >= 6) {
    return 0;
  }
  return cgo_callback_conv_get_reg_int(call, int_regs[conv->int_args++], bits);
}

uint64_t cgo_callback_conv_get_arg_uint(cgo_callback_call_t *call, int bits) {
  cgo_callback_sysv64_conv_t *conv = call->conv;

  // TODO: support more than 6 integer arguments.
  if (conv->int_args >= 6) {
    return 0;
  }
  return cgo_callback_conv_get_reg_uint(call, int_regs[conv->int_args++], bits);
}

// According to System V ABI, float args are passed through
// XMM0-7
static int float_regs[8] = {XMM0, XMM1, XMM2, XMM3, XMM4, XMM5, XMM6, XMM7};

float cgo_callback_conv_get_arg_single(cgo_callback_call_t *call) {
  cgo_callback_sysv64_conv_t *conv = call->conv;

  // TODO: support more than 8 float arguments.
  if (conv->float_args >= 8) {
    return 0;
  }

  return cgo_callback_conv_get_reg_single(call, float_regs[conv->float_args++]);
}

double cgo_callback_conv_get_arg_double(cgo_callback_call_t *call) {
  cgo_callback_sysv64_conv_t *conv = call->conv;
  double ret;

  // TODO: support more than 8 float arguments.
  if (conv->float_args >= 8) {
    return 0;
  }

  return cgo_callback_conv_get_reg_double(call, float_regs[conv->float_args++]);
}

void cgo_callback_conv_return(cgo_callback_call_t *call, void *val, int type, int bits) {
  int bytes = bits/8;

  if (type == TYPE_INT) {
    if (bytes <= 8) {
      memcpy((char *)call->reg + RAX, val, bytes);
    } else {
      memcpy((char *)call->reg + RAX, val, 8);
      memcpy((char *)call->reg + RDX, (char *)val+8, bytes-8);
    }
  } else if (type == TYPE_FLOAT) {
    if (bytes <= 16) {
      memcpy((char *)call->reg + XMM0, val, bytes);
    } else {
      memcpy((char *)call->reg + XMM0, val, 16);
      memcpy((char *)call->reg + XMM1, (char *)val+16, bytes-16);
    }
  }
}
