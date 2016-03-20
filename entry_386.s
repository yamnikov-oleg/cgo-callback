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

	// Use %ebp, because its value is last to be restored.
	mov %eax, %ebp

	pop %eax
	pop %ebx
	pop %ecx
	pop %edx
	pop %esi
	pop %edi
	fldl (%esp)
	add $8, %esp

	push %eax
	// Stack at this point:
	// (%esp) old eax
	// (%esp+$4) old ebp
	// (%esp+$8) port address
	// (%esp+$12) ret address
	// Gotta move it %ebp bytes up:
	// (%esp+%ebp+4) old eax
	// (%esp+%ebp+8) old ebp
	// (%esp+%ebp+12) ret address
	mov 12(%esp), %eax
	mov %eax, 12(%esp,%ebp,1)
	mov 4(%esp), %eax
	mov %eax, 8(%esp,%ebp,1)
	mov (%esp), %eax
	mov %eax, 4(%esp,%ebp,1)
	add %ebp, %esp
	add $4, %esp

	pop %eax
	pop %ebp
  ret
