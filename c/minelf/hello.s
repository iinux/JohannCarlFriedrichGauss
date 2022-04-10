# hello.s
#
# $ as --32 -o hello.o hello.s
# $ ld -melf_i386 --oformat=binary -o hello hello.o
# $ export PATH=./:$PATH
# $ hello 0 0 0
# hello
#
    .file "hello.s"
    .global _start, _load
    .equ   LOAD_ADDR, 0x00010000   # Page aligned load addr, here 64k
    .equ   E_ENTRY, LOAD_ADDR + (_start - _load)
    .equ   P_MEM_SZ, E_ENTRY
    .equ   P_FILE_SZ, P_MEM_SZ
_load:
    .byte  0x7F
    .ascii "ELF"              # e_ident, Magic Number
    .long  1                                      # p_type, loadable seg
    .long  0                                      # p_offset
    .long  LOAD_ADDR                              # p_vaddr
    .word  2                  # e_type, exec  # p_paddr
    .word  3                  # e_machine, Intel 386 target
    .long  P_FILE_SZ          # e_version     # p_filesz
    .long  E_ENTRY            # e_entry       # p_memsz
    .long  4                  # e_phoff       # p_flags, read(exec)
    .text
_start:
    popl   %eax    # argc     # e_shoff       # p_align
                   # 4 args, eax = 4, sys_write(fd, addr, len) : ebx, ecx, edx
                   # set 2nd eax = random addr to trigger bad syscall for exit
    popl   %ecx    # argv[0]
    mov    $5, %dl # str len  # e_flags
    int    $0x80
    loop   _start  # loop to popup a random addr as a bad syscall number
    .word  0x34               # e_ehsize = 52
    .word  0x20               # e_phentsize = 32
    .byte  1                  # e_phnum = 1, remove trailing 7 bytes with 0 value
                              # e_shentsize
                              # e_shnum
                              # e_shstrndx
