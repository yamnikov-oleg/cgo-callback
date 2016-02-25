void void_void(void *ptr) {
  ((void (*)(void))ptr)();
}

void void_int(void *ptr, int arg1) {
  ((void (*)(int))ptr)(arg1);
}

void void_uint(void *ptr, unsigned arg1) {
  ((void (*)(unsigned))ptr)(arg1);
}

void void_int_int(void *ptr, int arg1, int arg2) {
  ((void (*)(int, int))ptr)(arg1, arg2);
}
