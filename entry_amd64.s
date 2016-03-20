.text
.global cgo_callback_asm_entry
.type cgo_callback_asm_entry, @function

// Second, "assembly" entry point for amd64.
// It's responsibility is to provide C code with ability to read raw stack and
// to read and write registers. Every implemented architecture mush have its own
// version of this function.
cgo_callback_asm_entry:
	push %rbp
	mov	%rsp,	%rbp

	// Store all registers into memory. This register map will be used by higher
	// levels to load register values as they were at the moment of the call;
	// and to store register values for them to adopt at the moment of return.
	//
	// Some of these registers are callee-saved on different calling conventions,
	// some of them are used to pass arguments or return values.
	// It's easier to store all of them to support every possible conv.
	sub $0x80, %rsp
	movdqu %xmm7, 0x70(%rsp)
	movdqu %xmm6, 0x60(%rsp)
	movdqu %xmm5, 0x50(%rsp)
	movdqu %xmm4, 0x40(%rsp)
	movdqu %xmm3, 0x30(%rsp)
	movdqu %xmm2, 0x20(%rsp)
	movdqu %xmm1, 0x10(%rsp)
	movdqu %xmm0, 0x00(%rsp)
	push %r15
	push %r14
	push %r13
	push %r12
	push %r11
	push %r10
	push %r9
	push %r8
	push %rdi
	push %rsi
	push %rdx
	push %rcx
	push %rbx
	push %rax

	// Pass ptr to register map as second arg and ptr to the call stack
	// as first arg to the C entry point.
	mov %rsp, %rsi
	// Skip 8 bytes, containing old value of RBP.
	lea 8(%rbp), %rdi
  call cgo_callback_c_entry

	// Pop'em all.
	pop %rax
	pop %rbx
	pop %rcx
	pop %rdx
	pop %rsi
	pop %rdi
	pop %r8
	pop %r9
	pop %r10
	pop %r11
	pop %r12
	pop %r13
	pop %r14
	pop %r15
	movdqu 0x00(%rsp), %xmm0
	movdqu 0x10(%rsp), %xmm1
	movdqu 0x20(%rsp), %xmm2
	movdqu 0x30(%rsp), %xmm3
	movdqu 0x40(%rsp), %xmm4
	movdqu 0x50(%rsp), %xmm5
	movdqu 0x60(%rsp), %xmm6
	movdqu 0x70(%rsp), %xmm7
	add $0x80, %rsp

	pop %rbp
	// Discard 8 bytes from the stack, containing address of the port instruction.
	add $8, %rsp
  ret
