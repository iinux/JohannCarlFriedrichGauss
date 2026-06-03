; ============================================================
; 红黑树 (Red-Black Tree) - x86-64 Linux, NASM Intel 语法
; 算法依据: CLRS《算法导论》第 13 章
; 接口:
;   rb_init(tree*)              初始化
;   rb_insert(tree*, int key)   插入
;   rb_search(tree*, int key)   查找, 返回节点指针或 0
;   rb_delete(tree*, int key)   按键删除
;   rb_inorder(tree*)           中序遍历打印
;   rb_destroy(tree*)           释放所有节点
;
; 节点布局 (NODE_SIZE = 40 字节, 8 字节对齐):
;   +0  parent (qword)
;   +8  left   (qword)
;   +16 right  (qword)
;   +24 key    (dword, 32-bit signed int)
;   +28 color  (byte: 0=RED, 1=BLACK)
;   +29..+39 padding
;
; 树结构 (TREE_SIZE = 16 字节):
;   +0  root (qword)
;   +8  nil  (qword)  ; 哨兵节点, 始终 BLACK
; ============================================================

%define NODE_SIZE   40
%define N_PARENT    0
%define N_LEFT      8
%define N_RIGHT     16
%define N_KEY       24
%define N_COLOR     28

%define TREE_ROOT   0
%define TREE_NIL    8

%define RED         0
%define BLACK       1

extern  malloc, free, printf

global  rb_init, rb_insert, rb_search, rb_delete
global  rb_inorder, rb_destroy
global  main

section .rodata
fmt_int:    db "%d ", 0
fmt_nl:     db 10, 0
msg_search: db "search %d => %s", 10, 0
str_found:  db "FOUND", 0
str_miss:   db "NULL", 0

section .text

; ============================================================
; main - 演示 / 测试
; ============================================================
main:
    push    rbp
    mov     rbp, rsp
    sub     rsp, 32             ; 预留 16 字节 tree 结构 + 16 对齐

    lea     rdi, [rsp]
    call    rb_init

    ; 插入一组键
    lea     rdi, [rsp]
    mov     esi, 20
    call    rb_insert
    lea     rdi, [rsp]
    mov     esi, 10
    call    rb_insert
    lea     rdi, [rsp]
    mov     esi, 30
    call    rb_insert
    lea     rdi, [rsp]
    mov     esi, 5
    call    rb_insert
    lea     rdi, [rsp]
    mov     esi, 15
    call    rb_insert
    lea     rdi, [rsp]
    mov     esi, 25
    call    rb_insert
    lea     rdi, [rsp]
    mov     esi, 35
    call    rb_insert
    lea     rdi, [rsp]
    mov     esi, 1
    call    rb_insert
    lea     rdi, [rsp]
    mov     esi, 12
    call    rb_insert
    lea     rdi, [rsp]
    mov     esi, 28
    call    rb_insert

    ; 中序遍历
    lea     rdi, [rsp]
    call    rb_inorder
    lea     rdi, [rel fmt_nl]
    xor     eax, eax
    call    printf

    ; 删除几个
    lea     rdi, [rsp]
    mov     esi, 10
    call    rb_delete
    lea     rdi, [rsp]
    mov     esi, 30
    call    rb_delete
    lea     rdi, [rsp]
    mov     esi, 1
    call    rb_delete

    lea     rdi, [rsp]
    call    rb_inorder
    lea     rdi, [rel fmt_nl]
    xor     eax, eax
    call    printf

    ; 释放
    lea     rdi, [rsp]
    call    rb_destroy

    xor     eax, eax
    leave
    ret

; ============================================================
; rb_init(tree*)  - 初始化红黑树
; rdi = tree*
; ============================================================
rb_init:
    push    rbx
    mov     rbx, rdi

    mov     edi, NODE_SIZE
    call    malloc
    test    rax, rax
    jz      .done

    mov     qword [rax + N_PARENT], 0
    mov     qword [rax + N_LEFT],   0
    mov     qword [rax + N_RIGHT],  0
    mov     dword [rax + N_KEY],    0
    mov     byte  [rax + N_COLOR],  BLACK

    mov     [rbx + TREE_ROOT], rax
    mov     [rbx + TREE_NIL],  rax
.done:
    pop     rbx
    ret

; ============================================================
; _left_rotate(tree* in r12, node* in rax)
; 内部使用, 调用约定: r12=tree, rax=x; 破坏 rax/rcx/rdx/r8/r9
; ============================================================
_left_rotate:
    ; y = x.right
    mov     rcx, [rax + N_RIGHT]            ; rcx = y
    ; x.right = y.left
    mov     rdx, [rcx + N_LEFT]
    mov     [rax + N_RIGHT], rdx
    ; if y.left != nil: y.left.p = x
    mov     r8,  [r12 + TREE_NIL]
    cmp     rdx, r8
    je      .lr1
    mov     [rdx + N_PARENT], rax
.lr1:
    ; y.p = x.p
    mov     r9, [rax + N_PARENT]
    mov     [rcx + N_PARENT], r9
    ; if x.p == nil: T.root = y
    cmp     r9, r8
    jne     .lr_check_side
    mov     [r12 + TREE_ROOT], rcx
    jmp     .lr_attach
.lr_check_side:
    ; elif x == x.p.left: x.p.left = y else x.p.right = y
    mov     rdx, [r9 + N_LEFT]
    cmp     rax, rdx
    jne     .lr_right
    mov     [r9 + N_LEFT], rcx
    jmp     .lr_attach
.lr_right:
    mov     [r9 + N_RIGHT], rcx
.lr_attach:
    ; y.left = x; x.p = y
    mov     [rcx + N_LEFT],   rax
    mov     [rax + N_PARENT], rcx
    ret

; ============================================================
; _right_rotate(tree* in r12, node* in rax)
; ============================================================
_right_rotate:
    ; y = x.left
    mov     rcx, [rax + N_LEFT]             ; rcx = y
    ; x.left = y.right
    mov     rdx, [rcx + N_RIGHT]
    mov     [rax + N_LEFT], rdx
    mov     r8,  [r12 + TREE_NIL]
    cmp     rdx, r8
    je      .rr1
    mov     [rdx + N_PARENT], rax
.rr1:
    mov     r9, [rax + N_PARENT]
    mov     [rcx + N_PARENT], r9
    cmp     r9, r8
    jne     .rr_check_side
    mov     [r12 + TREE_ROOT], rcx
    jmp     .rr_attach
.rr_check_side:
    mov     rdx, [r9 + N_RIGHT]
    cmp     rax, rdx
    jne     .rr_left
    mov     [r9 + N_RIGHT], rcx
    jmp     .rr_attach
.rr_left:
    mov     [r9 + N_LEFT], rcx
.rr_attach:
    mov     [rcx + N_RIGHT],  rax
    mov     [rax + N_PARENT], rcx
    ret

; ============================================================
; rb_insert(tree* rdi, int esi)
; ============================================================
rb_insert:
    push    rbp
    mov     rbp, rsp
    push    rbx
    push    r12
    push    r13
    push    r14
    push    r15
    sub     rsp, 8                  ; 16 字节对齐

    mov     r12, rdi                ; r12 = tree
    mov     r13d, esi               ; r13d = key

    ; z = malloc(NODE_SIZE)
    mov     edi, NODE_SIZE
    call    malloc
    test    rax, rax
    jz      .ins_exit
    mov     r15, rax                ; r15 = z

    mov     r14, [r12 + TREE_NIL]   ; r14 = nil

    mov     [r15 + N_PARENT], r14
    mov     [r15 + N_LEFT],   r14
    mov     [r15 + N_RIGHT],  r14
    mov     dword [r15 + N_KEY],  r13d
    mov     byte  [r15 + N_COLOR], RED

    ; BST 下降
    mov     rax, [r12 + TREE_ROOT]  ; x
    mov     rbx, r14                ; y = nil
.ins_loop:
    cmp     rax, r14
    je      .ins_attach
    mov     rbx, rax
    mov     edx, [rax + N_KEY]
    cmp     r13d, edx
    jl      .ins_go_left
    mov     rax, [rax + N_RIGHT]
    jmp     .ins_loop
.ins_go_left:
    mov     rax, [rax + N_LEFT]
    jmp     .ins_loop
.ins_attach:
    mov     [r15 + N_PARENT], rbx
    cmp     rbx, r14
    jne     .ins_has_parent
    mov     [r12 + TREE_ROOT], r15
    jmp     .ins_fixup_call
.ins_has_parent:
    mov     edx, [rbx + N_KEY]
    cmp     r13d, edx
    jl      .ins_set_left
    mov     [rbx + N_RIGHT], r15
    jmp     .ins_fixup_call
.ins_set_left:
    mov     [rbx + N_LEFT], r15

.ins_fixup_call:
    ; ----- INSERT-FIXUP -----
    ; rbx = z (复用), 我们用 rbx 作 z, r12=tree, r14=nil
    mov     rbx, r15
.ifx_loop:
    mov     rcx, [rbx + N_PARENT]   ; rcx = z.p
    cmp     rcx, r14
    je      .ifx_done
    cmp     byte [rcx + N_COLOR], RED
    jne     .ifx_done

    mov     rdx, [rcx + N_PARENT]   ; rdx = z.p.p
    mov     r8,  [rdx + N_LEFT]
    cmp     rcx, r8
    jne     .ifx_p_right

    ; case A: z.p == z.p.p.left
    mov     r9, [rdx + N_RIGHT]     ; r9 = uncle
    cmp     byte [r9 + N_COLOR], RED
    jne     .ifx_a_else
    ; case 1
    mov     byte [rcx + N_COLOR], BLACK
    mov     byte [r9  + N_COLOR], BLACK
    mov     byte [rdx + N_COLOR], RED
    mov     rbx, rdx
    jmp     .ifx_loop
.ifx_a_else:
    ; case 2: z == z.p.right -> left-rotate(z.p); z = z.p
    mov     r8, [rcx + N_RIGHT]
    cmp     rbx, r8
    jne     .ifx_a_case3
    mov     rbx, rcx
    mov     rax, rbx
    call    _left_rotate
    mov     rcx, [rbx + N_PARENT]
    mov     rdx, [rcx + N_PARENT]
.ifx_a_case3:
    mov     byte [rcx + N_COLOR], BLACK
    mov     byte [rdx + N_COLOR], RED
    mov     rax, rdx
    call    _right_rotate
    jmp     .ifx_loop

.ifx_p_right:
    ; case B: z.p == z.p.p.right (镜像)
    mov     r9, [rdx + N_LEFT]      ; uncle
    cmp     byte [r9 + N_COLOR], RED
    jne     .ifx_b_else
    mov     byte [rcx + N_COLOR], BLACK
    mov     byte [r9  + N_COLOR], BLACK
    mov     byte [rdx + N_COLOR], RED
    mov     rbx, rdx
    jmp     .ifx_loop
.ifx_b_else:
    mov     r8, [rcx + N_LEFT]
    cmp     rbx, r8
    jne     .ifx_b_case3
    mov     rbx, rcx
    mov     rax, rbx
    call    _right_rotate
    mov     rcx, [rbx + N_PARENT]
    mov     rdx, [rcx + N_PARENT]
.ifx_b_case3:
    mov     byte [rcx + N_COLOR], BLACK
    mov     byte [rdx + N_COLOR], RED
    mov     rax, rdx
    call    _left_rotate
    jmp     .ifx_loop

.ifx_done:
    mov     rax, [r12 + TREE_ROOT]
    mov     byte [rax + N_COLOR], BLACK

.ins_exit:
    add     rsp, 8
    pop     r15
    pop     r14
    pop     r13
    pop     r12
    pop     rbx
    pop     rbp
    ret

; ============================================================
; rb_search(tree* rdi, int esi) -> rax (节点或 0)
; ============================================================
rb_search:
    mov     rax, [rdi + TREE_ROOT]
    mov     rcx, [rdi + TREE_NIL]
.s_loop:
    cmp     rax, rcx
    je      .s_miss
    mov     edx, [rax + N_KEY]
    cmp     esi, edx
    je      .s_done
    jl      .s_left
    mov     rax, [rax + N_RIGHT]
    jmp     .s_loop
.s_left:
    mov     rax, [rax + N_LEFT]
    jmp     .s_loop
.s_miss:
    xor     eax, eax
.s_done:
    ret

; ============================================================
; _transplant(tree* r12, node* u in rdi, node* v in rsi)
; ============================================================
_transplant:
    mov     r8, [r12 + TREE_NIL]
    mov     r9, [rdi + N_PARENT]
    cmp     r9, r8
    jne     .tp_check
    mov     [r12 + TREE_ROOT], rsi
    jmp     .tp_set_p
.tp_check:
    mov     rcx, [r9 + N_LEFT]
    cmp     rdi, rcx
    jne     .tp_right
    mov     [r9 + N_LEFT], rsi
    jmp     .tp_set_p
.tp_right:
    mov     [r9 + N_RIGHT], rsi
.tp_set_p:
    mov     [rsi + N_PARENT], r9
    ret

; ============================================================
; _tree_minimum(node* in rax, nil in rcx) -> rax
; ============================================================
_tree_minimum:
.tm_loop:
    mov     rdx, [rax + N_LEFT]
    cmp     rdx, rcx
    je      .tm_done
    mov     rax, rdx
    jmp     .tm_loop
.tm_done:
    ret

; ============================================================
; rb_delete(tree* rdi, int esi)
; ============================================================
rb_delete:
    push    rbp
    mov     rbp, rsp
    push    rbx
    push    r12
    push    r13
    push    r14
    push    r15
    sub     rsp, 8

    mov     r12, rdi                ; tree

    ; z = rb_search(tree, key)
    call    rb_search
    test    rax, rax
    jz      .del_exit
    mov     r13, rax                ; r13 = z

    mov     r14, [r12 + TREE_NIL]   ; r14 = nil

    mov     rbx, r13                ; y = z
    mov     r15b, [rbx + N_COLOR]   ; y_orig_color (低 8 位)

    mov     rcx, [r13 + N_LEFT]
    cmp     rcx, r14
    jne     .del_check_right

    ; z.left == nil: x = z.right; transplant(z, z.right)
    mov     rdx, [r13 + N_RIGHT]    ; x
    push    rdx                     ; 保存 x
    mov     rdi, r13
    mov     rsi, rdx
    call    _transplant
    pop     rdx
    mov     rbx, rdx                ; 复用 rbx 作为 x 传给 fixup
    jmp     .del_after

.del_check_right:
    mov     rdx, [r13 + N_RIGHT]
    cmp     rdx, r14
    jne     .del_two_children

    ; z.right == nil: x = z.left; transplant(z, z.left)
    mov     rdx, [r13 + N_LEFT]
    push    rdx
    mov     rdi, r13
    mov     rsi, rdx
    call    _transplant
    pop     rdx
    mov     rbx, rdx
    jmp     .del_after

.del_two_children:
    ; y = TREE-MINIMUM(z.right)
    mov     rax, [r13 + N_RIGHT]
    mov     rcx, r14
    call    _tree_minimum
    mov     rbx, rax                ; y
    mov     r15b, [rbx + N_COLOR]   ; y_orig_color
    mov     r8, [rbx + N_RIGHT]     ; x

    mov     r9, [rbx + N_PARENT]
    cmp     r9, r13
    jne     .del_y_far

    ; y.p == z: x.p = y
    mov     [r8 + N_PARENT], rbx
    jmp     .del_splice
.del_y_far:
    ; transplant(y, y.right); y.right = z.right; y.right.p = y
    push    r8
    mov     rdi, rbx
    mov     rsi, r8
    call    _transplant
    pop     r8
    mov     r9, [r13 + N_RIGHT]
    mov     [rbx + N_RIGHT], r9
    mov     [r9  + N_PARENT], rbx

.del_splice:
    ; transplant(z, y)
    push    r8
    mov     rdi, r13
    mov     rsi, rbx
    call    _transplant
    pop     r8
    ; y.left = z.left; y.left.p = y; y.color = z.color
    mov     r9, [r13 + N_LEFT]
    mov     [rbx + N_LEFT], r9
    mov     [r9  + N_PARENT], rbx
    mov     dl, [r13 + N_COLOR]
    mov     [rbx + N_COLOR], dl

    mov     rbx, r8                 ; rbx = x

.del_after:
    ; 释放 z 节点
    mov     rdi, r13
    push    rbx
    call    free
    pop     rbx

    ; 若 y_orig_color == BLACK 则 fixup(x in rbx)
    cmp     r15b, BLACK
    jne     .del_exit

    ; ---------- DELETE-FIXUP ----------
    ; rbx = x, r12 = tree, r14 = nil
.dfx_loop:
    mov     rax, [r12 + TREE_ROOT]
    cmp     rbx, rax
    je      .dfx_done
    cmp     byte [rbx + N_COLOR], BLACK
    jne     .dfx_done

    mov     rcx, [rbx + N_PARENT]   ; p
    mov     rdx, [rcx + N_LEFT]
    cmp     rbx, rdx
    jne     .dfx_x_right

    ; x 是左孩子
    mov     r8, [rcx + N_RIGHT]     ; w = sibling
    cmp     byte [r8 + N_COLOR], RED
    jne     .dfx_l_case234
    ; case 1
    mov     byte [r8 + N_COLOR], BLACK
    mov     byte [rcx + N_COLOR], RED
    mov     rax, rcx
    call    _left_rotate
    mov     rcx, [rbx + N_PARENT]
    mov     r8,  [rcx + N_RIGHT]

.dfx_l_case234:
    mov     r9,  [r8 + N_LEFT]
    mov     r10, [r8 + N_RIGHT]
    cmp     byte [r9  + N_COLOR], BLACK
    jne     .dfx_l_case34
    cmp     byte [r10 + N_COLOR], BLACK
    jne     .dfx_l_case34
    ; case 2
    mov     byte [r8 + N_COLOR], RED
    mov     rbx, rcx
    jmp     .dfx_loop
.dfx_l_case34:
    cmp     byte [r10 + N_COLOR], RED
    je      .dfx_l_case4
    ; case 3: w.right is BLACK, w.left is RED
    mov     byte [r9 + N_COLOR], BLACK
    mov     byte [r8 + N_COLOR], RED
    mov     rax, r8
    call    _right_rotate
    mov     rcx, [rbx + N_PARENT]
    mov     r8,  [rcx + N_RIGHT]
    mov     r10, [r8  + N_RIGHT]
.dfx_l_case4:
    mov     dl, [rcx + N_COLOR]
    mov     [r8 + N_COLOR], dl
    mov     byte [rcx + N_COLOR], BLACK
    mov     byte [r10 + N_COLOR], BLACK
    mov     rax, rcx
    call    _left_rotate
    mov     rbx, [r12 + TREE_ROOT]
    jmp     .dfx_done

.dfx_x_right:
    ; 镜像
    mov     r8, [rcx + N_LEFT]      ; w
    cmp     byte [r8 + N_COLOR], RED
    jne     .dfx_r_case234
    mov     byte [r8 + N_COLOR], BLACK
    mov     byte [rcx + N_COLOR], RED
    mov     rax, rcx
    call    _right_rotate
    mov     rcx, [rbx + N_PARENT]
    mov     r8,  [rcx + N_LEFT]

.dfx_r_case234:
    mov     r9,  [r8 + N_RIGHT]
    mov     r10, [r8 + N_LEFT]
    cmp     byte [r9  + N_COLOR], BLACK
    jne     .dfx_r_case34
    cmp     byte [r10 + N_COLOR], BLACK
    jne     .dfx_r_case34
    mov     byte [r8 + N_COLOR], RED
    mov     rbx, rcx
    jmp     .dfx_loop
.dfx_r_case34:
    cmp     byte [r10 + N_COLOR], RED
    je      .dfx_r_case4
    mov     byte [r9 + N_COLOR], BLACK
    mov     byte [r8 + N_COLOR], RED
    mov     rax, r8
    call    _left_rotate
    mov     rcx, [rbx + N_PARENT]
    mov     r8,  [rcx + N_LEFT]
    mov     r10, [r8  + N_LEFT]
.dfx_r_case4:
    mov     dl, [rcx + N_COLOR]
    mov     [r8 + N_COLOR], dl
    mov     byte [rcx + N_COLOR], BLACK
    mov     byte [r10 + N_COLOR], BLACK
    mov     rax, rcx
    call    _right_rotate
    mov     rbx, [r12 + TREE_ROOT]

.dfx_done:
    mov     byte [rbx + N_COLOR], BLACK

.del_exit:
    add     rsp, 8
    pop     r15
    pop     r14
    pop     r13
    pop     r12
    pop     rbx
    pop     rbp
    ret

; ============================================================
; rb_inorder(tree* rdi)
; ============================================================
rb_inorder:
    push    rbx
    push    r12
    push    r13
    sub     rsp, 8

    mov     r12, rdi
    mov     r13, [r12 + TREE_NIL]
    mov     rbx, [r12 + TREE_ROOT]
    call    _inorder_rec

    add     rsp, 8
    pop     r13
    pop     r12
    pop     rbx
    ret

; rbx = node, r13 = nil
_inorder_rec:
    cmp     rbx, r13
    je      .ior_ret
    push    rbx
    mov     rbx, [rbx + N_LEFT]
    call    _inorder_rec
    pop     rbx

    mov     esi, [rbx + N_KEY]
    lea     rdi, [rel fmt_int]
    xor     eax, eax
    call    printf

    push    rbx
    mov     rbx, [rbx + N_RIGHT]
    call    _inorder_rec
    pop     rbx
.ior_ret:
    ret

; ============================================================
; rb_destroy(tree* rdi) - 释放全部节点(含 nil)
; ============================================================
rb_destroy:
    push    rbx
    push    r12
    push    r13
    sub     rsp, 8

    mov     r12, rdi
    mov     r13, [r12 + TREE_NIL]
    mov     rbx, [r12 + TREE_ROOT]
    call    _free_rec

    mov     rdi, r13
    call    free
    mov     qword [r12 + TREE_ROOT], 0
    mov     qword [r12 + TREE_NIL],  0

    add     rsp, 8
    pop     r13
    pop     r12
    pop     rbx
    ret

; rbx = node, r13 = nil
_free_rec:
    cmp     rbx, r13
    je      .fr_ret
    push    rbx
    mov     rbx, [rbx + N_LEFT]
    call    _free_rec
    pop     rbx
    push    rbx
    mov     rbx, [rbx + N_RIGHT]
    call    _free_rec
    pop     rbx
    mov     rdi, rbx
    call    free
.fr_ret:
    ret
