// go run mkasm.go darwin amd64
// Code generated by the command above; DO NOT EDIT.

//go:build go1.13
// +build go1.13

#include "textflag.h"

TEXT libc_fdopendir_trampoline<>(SB),NOSPLIT,$0-0
	JMP	libc_fdopendir(SB)

GLOBL	·libc_fdopendir_trampoline_addr(SB), RODATA, $8
DATA	·libc_fdopendir_trampoline_addr(SB)/8, $libc_fdopendir_trampoline<>(SB)

TEXT libc_closedir_trampoline<>(SB),NOSPLIT,$0-0
	JMP	libc_closedir(SB)

GLOBL	·libc_closedir_trampoline_addr(SB), RODATA, $8
DATA	·libc_closedir_trampoline_addr(SB)/8, $libc_closedir_trampoline<>(SB)

TEXT libc_readdir_r_trampoline<>(SB),NOSPLIT,$0-0
	JMP	libc_readdir_r(SB)

GLOBL	·libc_readdir_r_trampoline_addr(SB), RODATA, $8
DATA	·libc_readdir_r_trampoline_addr(SB)/8, $libc_readdir_r_trampoline<>(SB)
