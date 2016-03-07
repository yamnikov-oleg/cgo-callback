.text
.global cgo_callback_asm_entry
.type cgo_callback_asm_entry, @function

cgo_callback_asm_entry:
	push %ebp
	mov	%esp,	%ebp

	sub $8, %esp
	fstpl (%esp)
	push %edi
	push %esi
	push %edx
	push %ecx
	push %ebx
	push %eax

  mov %esp, %eax
  sub $8, %esp
	mov %eax, 4(%esp)
	lea 4(%ebp), %eax
  mov %eax, (%esp)
  call cgo_callback_c_entry
  add $8, %esp

	pop %eax
	pop %ebx
	pop %ecx
	pop %edx
	pop %esi
	pop %edi
	fldl (%esp)
	add $8, %esp

	pop %ebp
	add $4, %esp
  ret $0
