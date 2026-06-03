; ============================================================
; rbtree.asm — 红黑树完整实现 (x86-64 NASM, Linux)
; 编译: nasm -f elf64 -o rbtree.o rbtree.asm
;        gcc -no-pie -o rbtree rbtree.o
; ============================================================
;
; 红黑树性质:
;   1. 每个节点是红色或黑色
;   2. 根节点是黑色
;   3. 所有叶子 (哨兵 NIL) 是黑色
;   4. 红色节点的子节点必须是黑色
;   5. 从任意节点到其叶子的所有路径包含相同数量的黑色节点
;
; 节点结构 (40 字节):
;   offset  size   field
;   0       8      key      (键值)
;   8       8      color    (0=RED, 1=BLACK)
;   16      8      left     (左子指针)
;   24      8      right    (右子指针)
;   32      8      parent   (父指针)
;
; 树结构 (8 字节):
;   offset  size   field
;   0       8      root     (根节点指针)

extern printf
extern malloc
extern free

; ============================================================
; 常量
; ============================================================
RED     equ 0
BLACK   equ 1

KEY     equ 0
COLOR   equ 8
LEFT    equ 16
RIGHT   equ 24
PARENT  equ 32
NODE_SIZE equ 40

ROOT    equ 0

; ============================================================
; 数据段
; ============================================================
section .data

; 全局哨兵节点 — 所有叶子指向此处，根父节点也指向此处
; 访问哨兵的 color 字段始终得到 BLACK，无需空指针检查
align 8
sentinel:
    dq 0            ; key = 0 (未使用)
    dq BLACK        ; color = BLACK
    dq sentinel     ; left = 自身
    dq sentinel     ; right = 自身
    dq sentinel     ; parent = 自身

; 格式字符串
fmt_key:      db "%d ", 0
fmt_newline:  db 10, 0
fmt_insert:   db "插入 %d...", 10, 0
fmt_delete:   db "删除 %d...", 10, 0
fmt_search:   db "搜索 %d: ", 0
fmt_found:    db "找到", 10, 0
fmt_notfound: db "未找到", 10, 0
fmt_inorder:  db "中序遍历: ", 0
fmt_verify:   db "验证红黑树性质... ", 0
fmt_valid:    db "有效", 10, 0
fmt_invalid:  db "无效!", 10, 0
fmt_destroy:  db "销毁树...", 10, 0
fmt_treeinfo: db "=== 红黑树测试 ===", 10, 0

; 测试数据
test_data:     dq 10, 5, 15, 3, 7, 12, 18, 1, 4, 6, 8
test_data_cnt equ ($ - test_data) / 8

; ============================================================
; 代码段
; ============================================================
section .text

; ============================================================
; rb_create_tree() — 创建空红黑树
; 返回: rax = 树指针 (0 表示分配失败)
; ============================================================
global rb_create_tree
rb_create_tree:
    push rbp
    mov rbp, rsp

    mov rdi, 8              ; 分配 8 字节树结构
    call malloc
    test rax, rax
    jz  .create_fail

    mov qword [rax + ROOT], sentinel   ; 根指向哨兵

    pop rbp
    ret

.create_fail:
    xor eax, eax
    pop rbp
    ret

; ============================================================
; rb_make_node(key) — 创建新节点 (颜色为红色)
; 参数: rdi = key
; 返回: rax = 节点指针 (0 表示分配失败)
; ============================================================
rb_make_node:
    push rbp
    mov rbp, rsp
    push rbx

    mov rbx, rdi            ; 保存 key

    mov rdi, NODE_SIZE
    call malloc
    test rax, rax
    jz  .node_fail

    mov [rax + KEY],    rbx
    mov qword [rax + COLOR], RED
    mov qword [rax + LEFT],   sentinel
    mov qword [rax + RIGHT],  sentinel
    mov qword [rax + PARENT], sentinel

    pop rbx
    pop rbp
    ret

.node_fail:
    xor eax, eax
    pop rbx
    pop rbp
    ret

; ============================================================
; rb_rotate_left(tree, x) — 左旋
;   设 y = x.right
;   使 y 取代 x 的位置，x 成为 y 的左子
; 参数: rdi = tree, rsi = x
; ============================================================
rb_rotate_left:
    push rbp
    mov rbp, rsp
    push rbx
    push r12
    push r13

    mov r12, rdi            ; tree
    mov r13, rsi            ; x

    mov rbx, [r13 + RIGHT]  ; rbx = y = x.right

    ; x.right = y.left
    mov rax, [rbx + LEFT]
    mov [r13 + RIGHT], rax

    ; if y.left != sentinel: y.left.parent = x
    cmp qword [rbx + LEFT], sentinel
    je  .ll_parent_done
    mov rax, [rbx + LEFT]
    mov [rax + PARENT], r13
.ll_parent_done:

    ; y.parent = x.parent
    mov rax, [r13 + PARENT]
    mov [rbx + PARENT], rax

    ; 将 y 连接到 x.parent
    cmp rax, sentinel
    jne .ll_check_left
    mov [r12 + ROOT], rbx
    jmp .ll_link

.ll_check_left:
    mov rcx, [rax + LEFT]
    cmp rcx, r13
    jne .ll_right_child
    mov [rax + LEFT], rbx
    jmp .ll_link

.ll_right_child:
    mov [rax + RIGHT], rbx

.ll_link:
    mov [rbx + LEFT], r13
    mov [r13 + PARENT], rbx

    pop r13
    pop r12
    pop rbx
    pop rbp
    ret

; ============================================================
; rb_rotate_right(tree, x) — 右旋 (左旋对称)
;   设 y = x.left
;   使 y 取代 x 的位置，x 成为 y 的右子
; 参数: rdi = tree, rsi = x
; ============================================================
rb_rotate_right:
    push rbp
    mov rbp, rsp
    push rbx
    push r12
    push r13

    mov r12, rdi            ; tree
    mov r13, rsi            ; x

    mov rbx, [r13 + LEFT]   ; rbx = y = x.left

    ; x.left = y.right
    mov rax, [rbx + RIGHT]
    mov [r13 + LEFT], rax

    ; if y.right != sentinel: y.right.parent = x
    cmp qword [rbx + RIGHT], sentinel
    je  .lr_parent_done
    mov rax, [rbx + RIGHT]
    mov [rax + PARENT], r13
.lr_parent_done:

    ; y.parent = x.parent
    mov rax, [r13 + PARENT]
    mov [rbx + PARENT], rax

    ; 将 y 连接到 x.parent
    cmp rax, sentinel
    jne .lr_check_right
    mov [r12 + ROOT], rbx
    jmp .lr_link

.lr_check_right:
    mov rcx, [rax + RIGHT]
    cmp rcx, r13
    jne .lr_left_child
    mov [rax + RIGHT], rbx
    jmp .lr_link

.lr_left_child:
    mov [rax + LEFT], rbx

.lr_link:
    mov [rbx + RIGHT], r13
    mov [r13 + PARENT], rbx

    pop r13
    pop r12
    pop rbx
    pop rbp
    ret

; ============================================================
; rb_transplant(tree, u, v) — 用子树 v 替换子树 u
;   只修改父指针连接，不处理子指针
; 参数: rdi = tree, rsi = u, rdx = v
; ============================================================
rb_transplant:
    push rbp
    mov rbp, rsp

    mov rax, [rsi + PARENT] ; rax = u.parent

    cmp rax, sentinel
    jne .tp_check_left

    mov [rdi + ROOT], rdx   ; u 是根
    jmp .tp_set_parent

.tp_check_left:
    mov rcx, [rax + LEFT]
    cmp rcx, rsi
    jne .tp_right

    mov [rax + LEFT], rdx
    jmp .tp_set_parent

.tp_right:
    mov [rax + RIGHT], rdx

.tp_set_parent:
    mov [rdx + PARENT], rax

    pop rbp
    ret

; ============================================================
; rb_search(tree, key) — 搜索键值
; 参数: rdi = tree, rsi = key
; 返回: rax = 节点指针 (哨兵表示未找到)
; ============================================================
global rb_search
rb_search:
    push rbp
    mov rbp, rsp

    mov rax, [rdi + ROOT]   ; cur = root

.srch_loop:
    cmp rax, sentinel
    je  .srch_not_found

    cmp rsi, [rax + KEY]
    je  .srch_done
    jl  .srch_left
    ; jg .srch_right
.srch_right:
    mov rax, [rax + RIGHT]
    jmp .srch_loop

.srch_left:
    mov rax, [rax + LEFT]
    jmp .srch_loop

.srch_not_found:
    mov rax, sentinel

.srch_done:
    pop rbp
    ret

; ============================================================
; rb_minimum(node) — 子树最小节点
; 参数: rdi = node
; 返回: rax = 最小节点 (哨兵表示空子树)
; ============================================================
rb_minimum:
    push rbp
    mov rbp, rsp

    mov rax, rdi

.min_loop:
    cmp qword [rax + LEFT], sentinel
    je  .min_done
    mov rax, [rax + LEFT]
    jmp .min_loop

.min_done:
    pop rbp
    ret

; ============================================================
; rb_maximum(node) — 子树最大节点
; 参数: rdi = node
; 返回: rax = 最大节点 (哨兵表示空子树)
; ============================================================
rb_maximum:
    push rbp
    mov rbp, rsp

    mov rax, rdi

.max_loop:
    cmp qword [rax + RIGHT], sentinel
    je  .max_done
    mov rax, [rax + RIGHT]
    jmp .max_loop

.max_done:
    pop rbp
    ret

; ============================================================
; rb_successor(node) — 中序后继
; 参数: rdi = node
; 返回: rax = 后继节点 (哨兵表示无后继)
; ============================================================
rb_successor:
    push rbp
    mov rbp, rsp
    push rbx

    mov rbx, rdi

    cmp qword [rbx + RIGHT], sentinel
    je  .succ_no_right

    mov rdi, [rbx + RIGHT]
    call rb_minimum
    pop rbx
    pop rbp
    ret

.succ_no_right:
    mov rax, [rbx + PARENT]

.succ_loop:
    cmp rax, sentinel
    je  .succ_done

    mov rcx, [rax + LEFT]
    cmp rcx, rbx
    je  .succ_done

    mov rbx, rax
    mov rax, [rbx + PARENT]
    jmp .succ_loop

.succ_done:
    pop rbx
    pop rbp
    ret

; ============================================================
; rb_predecessor(node) — 中序前驱
; 参数: rdi = node
; 返回: rax = 前驱节点 (哨兵表示无前驱)
; ============================================================
rb_predecessor:
    push rbp
    mov rbp, rsp
    push rbx

    mov rbx, rdi

    cmp qword [rbx + LEFT], sentinel
    je  .pred_no_left

    mov rdi, [rbx + LEFT]
    call rb_maximum
    pop rbx
    pop rbp
    ret

.pred_no_left:
    mov rax, [rbx + PARENT]

.pred_loop:
    cmp rax, sentinel
    je  .pred_done

    mov rcx, [rax + RIGHT]
    cmp rcx, rbx
    je  .pred_done

    mov rbx, rax
    mov rax, [rbx + PARENT]
    jmp .pred_loop

.pred_done:
    pop rbx
    pop rbp
    ret

; ============================================================
; rb_insert_fixup(tree, z) — 插入后修复红黑性质
; 参数: rdi = tree, rsi = z (新插入的红色节点)
; ============================================================
rb_insert_fixup:
    push rbp
    mov rbp, rsp
    push rbx
    push r12
    push r13
    push r14

    mov r12, rdi            ; tree
    mov r13, rsi            ; z

.if_loop:
    ; while z.parent.color == RED
    mov rax, [r13 + PARENT]     ; rax = z.parent
    cmp qword [rax + COLOR], RED
    jne .if_done

    ; if z.parent == z.parent.parent.left
    mov r14, [rax + PARENT]     ; r14 = grandparent
    cmp rax, [r14 + LEFT]
    jne .if_uncle_left

    ; 叔节点 = grandparent.right
    mov rbx, [r14 + RIGHT]      ; rbx = uncle (y)

    cmp qword [rbx + COLOR], RED
    jne .if_case2_left

    ; === Case 1: 叔节点为红色 ===
    ; 父和叔变黑，祖父变红，z 上移两层
    mov qword [rax + COLOR], BLACK
    mov qword [rbx + COLOR], BLACK
    mov qword [r14 + COLOR], RED
    mov r13, r14
    jmp .if_loop

.if_case2_left:
    ; === Case 2: z 是右孩子 ===
    cmp r13, [rax + RIGHT]
    jne .if_case3_left

    mov r13, rax                ; z = z.parent
    mov rdi, r12
    mov rsi, r13
    call rb_rotate_left
    mov rax, [r13 + PARENT]     ; 刷新 parent

.if_case3_left:
    ; === Case 3: z 是左孩子 ===
    mov qword [rax + COLOR], BLACK
    mov r14, [rax + PARENT]
    mov qword [r14 + COLOR], RED
    mov rdi, r12
    mov rsi, r14
    call rb_rotate_right
    jmp .if_done

    ; ======= 对称情况 (parent 是右孩子) =======
.if_uncle_left:
    mov rbx, [r14 + LEFT]       ; rbx = uncle

    cmp qword [rbx + COLOR], RED
    jne .if_case2_right

    ; Case 1 对称
    mov qword [rax + COLOR], BLACK
    mov qword [rbx + COLOR], BLACK
    mov qword [r14 + COLOR], RED
    mov r13, r14
    jmp .if_loop

.if_case2_right:
    cmp r13, [rax + LEFT]
    jne .if_case3_right

    mov r13, rax
    mov rdi, r12
    mov rsi, r13
    call rb_rotate_right
    mov rax, [r13 + PARENT]

.if_case3_right:
    mov qword [rax + COLOR], BLACK
    mov r14, [rax + PARENT]
    mov qword [r14 + COLOR], RED
    mov rdi, r12
    mov rsi, r14
    call rb_rotate_left

.if_done:
    ; 根节点保持黑色
    mov rax, [r12 + ROOT]
    mov qword [rax + COLOR], BLACK

    pop r14
    pop r13
    pop r12
    pop rbx
    pop rbp
    ret

; ============================================================
; rb_insert(tree, key) — 插入键值
; 参数: rdi = tree, rsi = key
; 返回: rax = 0 成功, 1 键已存在, 2 内存不足
; ============================================================
global rb_insert
rb_insert:
    push rbp
    mov rbp, rsp
    push rbx
    push r12
    push r13
    push r14

    mov r12, rdi            ; tree
    mov r13, rsi            ; key

    ; 创建新节点 (红色)
    mov rdi, r13
    call rb_make_node
    test rax, rax
    jz  .ins_oom
    mov r14, rax            ; r14 = z (新节点)

    ; BST 查找插入位置
    mov rbx, [r12 + ROOT]   ; cur
    xor rcx, rcx
    mov rcx, sentinel       ; parent = NIL

    cmp rbx, sentinel
    je  .ins_set_root       ; 空树

.ins_search:
    cmp rbx, sentinel
    je  .ins_attach

    mov rcx, rbx            ; parent = cur
    mov rax, [rbx + KEY]

    cmp r13, rax
    je  .ins_key_exists
    jl  .ins_go_left
    jg  .ins_go_right

.ins_go_left:
    mov rbx, [rbx + LEFT]
    jmp .ins_search

.ins_go_right:
    mov rbx, [rbx + RIGHT]
    jmp .ins_search

.ins_set_root:
    mov [r12 + ROOT], r14
    mov qword [r14 + PARENT], sentinel
    jmp .ins_fixup

.ins_attach:
    mov [r14 + PARENT], rcx

    cmp r13, [rcx + KEY]
    jl  .ins_attach_left
    jg  .ins_attach_right

.ins_attach_left:
    mov [rcx + LEFT], r14
    jmp .ins_fixup

.ins_attach_right:
    mov [rcx + RIGHT], r14

.ins_fixup:
    mov rdi, r12
    mov rsi, r14
    call rb_insert_fixup

    xor eax, eax
    pop r14
    pop r13
    pop r12
    pop rbx
    pop rbp
    ret

.ins_key_exists:
    mov rdi, r14
    call free
    mov eax, 1
    pop r14
    pop r13
    pop r12
    pop rbx
    pop rbp
    ret

.ins_oom:
    mov eax, 2
    pop r14
    pop r13
    pop r12
    pop rbx
    pop rbp
    ret

; ============================================================
; rb_delete_fixup(tree, x) — 删除后修复红黑性质
; 参数: rdi = tree, rsi = x
; 算法: CLRS 第三版 §13.4
; ============================================================
rb_delete_fixup:
    push rbp
    mov rbp, rsp
    push rbx
    push r12
    push r13
    push r14
    push r15

    mov r12, rdi            ; tree
    mov r13, rsi            ; x

.df_loop:
    ; while x != root AND x.color == BLACK
    mov rax, [r12 + ROOT]
    cmp r13, rax
    je  .df_done

    cmp qword [r13 + COLOR], BLACK
    jne .df_done

    mov rbx, [r13 + PARENT] ; rbx = x.parent

    ; ===== x 是左孩子 =====
    cmp r13, [rbx + LEFT]
    jne .df_symmetric

    mov r14, [rbx + RIGHT]  ; r14 = w (兄弟节点)

    ; Case 1: w 为红色 → 转为 Case 2/3/4
    cmp qword [r14 + COLOR], RED
    jne .df_case2

    mov qword [r14 + COLOR], BLACK
    mov qword [rbx + COLOR], RED
    mov rdi, r12
    mov rsi, rbx
    call rb_rotate_left
    mov rbx, [r13 + PARENT]
    mov r14, [rbx + RIGHT]

.df_case2:
    ; Case 2: w 的两个孩子都是黑色
    mov rax, [r14 + LEFT]
    cmp qword [rax + COLOR], BLACK
    jne .df_case3

    mov rax, [r14 + RIGHT]
    cmp qword [rax + COLOR], BLACK
    jne .df_case3

    mov qword [r14 + COLOR], RED
    mov r13, rbx
    jmp .df_loop

.df_case3:
    ; Case 3: w.right 为黑色, w.left 为红色
    mov rax, [r14 + RIGHT]
    cmp qword [rax + COLOR], BLACK
    jne .df_case4

    mov rax, [r14 + LEFT]
    mov qword [rax + COLOR], BLACK
    mov qword [r14 + COLOR], RED
    mov rdi, r12
    mov rsi, r14
    call rb_rotate_right
    mov rbx, [r13 + PARENT]
    mov r14, [rbx + RIGHT]

.df_case4:
    ; Case 4: w.right 为红色
    mov rax, [rbx + COLOR]
    mov [r14 + COLOR], rax
    mov qword [rbx + COLOR], BLACK
    mov rax, [r14 + RIGHT]
    mov qword [rax + COLOR], BLACK
    mov rdi, r12
    mov rsi, rbx
    call rb_rotate_left
    mov r13, [r12 + ROOT]
    jmp .df_loop

    ; ===== x 是右孩子 (对称) =====
.df_symmetric:
    mov r14, [rbx + LEFT]   ; r14 = w (兄弟)

    cmp qword [r14 + COLOR], RED
    jne .df_sym_case2

    mov qword [r14 + COLOR], BLACK
    mov qword [rbx + COLOR], RED
    mov rdi, r12
    mov rsi, rbx
    call rb_rotate_right
    mov rbx, [r13 + PARENT]
    mov r14, [rbx + LEFT]

.df_sym_case2:
    mov rax, [r14 + LEFT]
    cmp qword [rax + COLOR], BLACK
    jne .df_sym_case3

    mov rax, [r14 + RIGHT]
    cmp qword [rax + COLOR], BLACK
    jne .df_sym_case3

    mov qword [r14 + COLOR], RED
    mov r13, rbx
    jmp .df_loop

.df_sym_case3:
    mov rax, [r14 + LEFT]
    cmp qword [rax + COLOR], BLACK
    jne .df_sym_case4

    mov rax, [r14 + RIGHT]
    mov qword [rax + COLOR], BLACK
    mov qword [r14 + COLOR], RED
    mov rdi, r12
    mov rsi, r14
    call rb_rotate_left
    mov rbx, [r13 + PARENT]
    mov r14, [rbx + LEFT]

.df_sym_case4:
    mov rax, [rbx + COLOR]
    mov [r14 + COLOR], rax
    mov qword [rbx + COLOR], BLACK
    mov rax, [r14 + LEFT]
    mov qword [rax + COLOR], BLACK
    mov rdi, r12
    mov rsi, rbx
    call rb_rotate_right
    mov r13, [r12 + ROOT]
    jmp .df_loop

.df_done:
    mov qword [r13 + COLOR], BLACK

    pop r15
    pop r14
    pop r13
    pop r12
    pop rbx
    pop rbp
    ret

; ============================================================
; rb_delete(tree, z) — 删除节点 z
; 参数: rdi = tree, rsi = z (节点指针, 不能为哨兵)
; 返回值: 无
; ============================================================
global rb_delete
rb_delete:
    push rbp
    mov rbp, rsp
    push rbx
    push r12
    push r13
    push r14
    push r15

    mov r12, rdi            ; tree
    mov r13, rsi            ; z (要删除的节点)

    cmp r13, sentinel
    je  .del_done

    mov r14, r13            ; y = z
    mov r15, [r14 + COLOR]  ; y_original_color

    ; ===== 情况 A: z 没有左孩子 =====
    cmp qword [r13 + LEFT], sentinel
    jne .del_check_right

    mov rbx, [r13 + RIGHT]  ; x = z.right
    mov rdi, r12
    mov rsi, r13
    mov rdx, rbx
    call rb_transplant

    jmp .del_check_color

    ; ===== 情况 B: z 没有右孩子 =====
.del_check_right:
    cmp qword [r13 + RIGHT], sentinel
    jne .del_both

    mov rbx, [r13 + LEFT]   ; x = z.left
    mov rdi, r12
    mov rsi, r13
    mov rdx, rbx
    call rb_transplant

    jmp .del_check_color

    ; ===== 情况 C: z 有两个孩子 =====
.del_both:
    mov rdi, [r13 + RIGHT]
    call rb_minimum
    mov r14, rax            ; y = z.right 的最小节点

    mov r15, [r14 + COLOR]  ; y_original_color
    mov rbx, [r14 + RIGHT]  ; x = y.right

    cmp qword [r14 + PARENT], r13
    je  .del_y_parent_z

    ; transplant(y, y.right)
    mov rdi, r12
    mov rsi, r14
    mov rdx, rbx
    call rb_transplant

    ; y.right = z.right
    mov rax, [r13 + RIGHT]
    mov [r14 + RIGHT], rax
    mov [rax + PARENT], r14

.del_y_parent_z:
    ; x.parent = y (CLRS: 即使 x 是哨兵也要设置)
    mov [rbx + PARENT], r14

    ; transplant(z, y)
    mov rdi, r12
    mov rsi, r13
    mov rdx, r14
    call rb_transplant

    ; y.left = z.left
    mov rax, [r13 + LEFT]
    mov [r14 + LEFT], rax
    mov [rax + PARENT], r14

    ; y.color = z.color
    mov rax, [r13 + COLOR]
    mov [r14 + COLOR], rax

.del_check_color:
    ; 如果 y 原为黑色，需要修复
    cmp r15, BLACK
    jne .del_skip_fixup

    mov rdi, r12
    mov rsi, rbx
    call rb_delete_fixup

.del_skip_fixup:
    mov rdi, r13
    call free               ; 释放原节点

.del_done:
    pop r15
    pop r14
    pop r13
    pop r12
    pop rbx
    pop rbp
    ret

; ============================================================
; rb_delete_key(tree, key) — 按键值删除
; 参数: rdi = tree, rsi = key
; 返回: rax = 0 成功, 1 未找到
; ============================================================
global rb_delete_key
rb_delete_key:
    push rbp
    mov rbp, rsp
    push rbx
    push r12

    mov r12, rdi            ; tree
    mov rbx, rsi            ; key

    mov rsi, rbx
    call rb_search

    cmp rax, sentinel
    je  .delk_not_found

    mov rdi, r12
    mov rsi, rax
    call rb_delete

    xor eax, eax
    pop r12
    pop rbx
    pop rbp
    ret

.delk_not_found:
    mov eax, 1
    pop r12
    pop rbx
    pop rbp
    ret

; ============================================================
; rb_print_inorder_helper(node) — 递归中序遍历打印
; 参数: rdi = node
; ============================================================
rb_print_inorder_helper:
    push rbp
    mov rbp, rsp
    push rbx

    mov rbx, rdi
    cmp rbx, sentinel
    je  .ph_done

    mov rdi, [rbx + LEFT]
    call rb_print_inorder_helper

    mov rdi, fmt_key
    mov rsi, [rbx + KEY]
    xor eax, eax
    call printf

    mov rdi, [rbx + RIGHT]
    call rb_print_inorder_helper

.ph_done:
    pop rbx
    pop rbp
    ret

; ============================================================
; rb_print_inorder(tree) — 中序打印整棵树
; 参数: rdi = tree
; ============================================================
global rb_print_inorder
rb_print_inorder:
    push rbp
    mov rbp, rsp
    push rbx

    mov rbx, rdi            ; rbx = tree (callee-saved, printf 不破坏)

    mov rdi, fmt_inorder
    xor eax, eax
    call printf

    mov rdi, [rbx + ROOT]   ; tree->root
    call rb_print_inorder_helper

    mov rdi, fmt_newline
    xor eax, eax
    call printf

    pop rbx
    pop rbp
    ret

; ============================================================
; rb_verify_helper(node) — 递归验证红黑性质
;   检查: 无连续红节点、黑色高度一致
; 参数: rdi = node
; 返回: rax = 黑色高度 (-1 表示无效)
; ============================================================
rb_verify_helper:
    push rbp
    mov rbp, rsp
    push rbx
    push r12
    push r13

    mov rbx, rdi

    cmp rbx, sentinel
    jne .vh_node

    mov eax, 1              ; 哨兵自身是黑色，黑高 = 1
    pop r13
    pop r12
    pop rbx
    pop rbp
    ret

.vh_node:
    ; 红色节点的子节点不能是红色
    cmp qword [rbx + COLOR], RED
    jne .vh_skip_red_check

    mov rax, [rbx + LEFT]
    cmp qword [rax + COLOR], RED
    je  .vh_invalid

    mov rax, [rbx + RIGHT]
    cmp qword [rax + COLOR], RED
    je  .vh_invalid

.vh_skip_red_check:
    ; 递归左子树
    mov rdi, [rbx + LEFT]
    call rb_verify_helper
    cmp eax, -1
    je  .vh_invalid
    mov r12d, eax           ; left_bh

    ; 递归右子树
    mov rdi, [rbx + RIGHT]
    call rb_verify_helper
    cmp eax, -1
    je  .vh_invalid
    mov r13d, eax           ; right_bh

    ; 左右黑高必须相等
    cmp r12d, r13d
    jne .vh_invalid

    ; 计算当前节点黑高
    mov eax, r12d
    cmp qword [rbx + COLOR], BLACK
    jne .vh_return
    inc eax

.vh_return:
    pop r13
    pop r12
    pop rbx
    pop rbp
    ret

.vh_invalid:
    mov eax, -1
    pop r13
    pop r12
    pop rbx
    pop rbp
    ret

; ============================================================
; rb_verify(tree) — 验证整棵红黑树
; 参数: rdi = tree
; 返回: rax = 1 有效, 0 无效
; ============================================================
global rb_verify
rb_verify:
    push rbp
    mov rbp, rsp

    mov rax, [rdi + ROOT]

    cmp rax, sentinel
    je  .vfy_valid          ; 空树有效

    ; 根必须为黑色
    cmp qword [rax + COLOR], BLACK
    jne .vfy_invalid

    mov rdi, rax
    call rb_verify_helper
    cmp eax, -1
    je  .vfy_invalid

    mov eax, 1
    pop rbp
    ret

.vfy_invalid:
    xor eax, eax
    pop rbp
    ret

.vfy_valid:
    mov eax, 1
    pop rbp
    ret

; ============================================================
; rb_destroy_helper(node) — 后序遍历释放所有节点
; 参数: rdi = node
; ============================================================
rb_destroy_helper:
    push rbp
    mov rbp, rsp
    push rbx

    mov rbx, rdi
    cmp rbx, sentinel
    je  .dh_done

    mov rdi, [rbx + LEFT]
    call rb_destroy_helper

    mov rdi, [rbx + RIGHT]
    call rb_destroy_helper

    mov rdi, rbx
    call free

.dh_done:
    pop rbx
    pop rbp
    ret

; ============================================================
; rb_destroy_tree(tree) — 销毁整棵树
; 参数: rdi = tree
; ============================================================
global rb_destroy_tree
rb_destroy_tree:
    push rbp
    mov rbp, rsp
    push rbx

    mov rbx, rdi

    mov rdi, [rbx + ROOT]
    call rb_destroy_helper

    mov rdi, rbx
    call free

    pop rbx
    pop rbp
    ret

; ============================================================
; main — 测试入口
; ============================================================
global main
main:
    push rbp
    mov rbp, rsp
    push rbx
    push r12

    mov rdi, fmt_treeinfo
    xor eax, eax
    call printf

    ; === 创建树 ===
    call rb_create_tree
    test rax, rax
    jz  .exit_fail
    mov r12, rax            ; r12 = tree

    ; === 插入测试数据 ===
    xor ebx, ebx            ; i = 0
.ins_loop:
    cmp rbx, test_data_cnt
    jae .ins_done

    mov rdi, fmt_insert
    mov rsi, [test_data + rbx*8]
    xor eax, eax
    call printf

    mov rdi, r12
    mov rsi, [test_data + rbx*8]
    call rb_insert

    inc rbx
    jmp .ins_loop
.ins_done:

    ; === 中序遍历 ===
    mov rdi, r12
    call rb_print_inorder

    ; === 验证 ===
    mov rdi, fmt_verify
    xor eax, eax
    call printf

    mov rdi, r12
    call rb_verify
    test eax, eax
    jz  .invalid

    mov rdi, fmt_valid
    xor eax, eax
    call printf
    jmp .test_search

.invalid:
    mov rdi, fmt_invalid
    xor eax, eax
    call printf

    ; === 搜索测试 ===
.test_search:
    mov rdi, fmt_search
    mov rsi, 7
    xor eax, eax
    call printf

    mov rdi, r12
    mov rsi, 7
    call rb_search
    cmp rax, sentinel
    je  .notfound7

    mov rdi, fmt_found
    xor eax, eax
    call printf
    jmp .test_delete

.notfound7:
    mov rdi, fmt_notfound
    xor eax, eax
    call printf

    ; === 删除测试 ===
.test_delete:
    mov rdi, fmt_delete
    mov rsi, 5
    xor eax, eax
    call printf

    mov rdi, r12
    mov rsi, 5
    call rb_delete_key

    mov rdi, r12
    call rb_print_inorder

    mov rdi, fmt_verify
    xor eax, eax
    call printf

    mov rdi, r12
    call rb_verify
    test eax, eax
    jz  .invalid2

    mov rdi, fmt_valid
    xor eax, eax
    call printf
    jmp .del_more

.invalid2:
    mov rdi, fmt_invalid
    xor eax, eax
    call printf

.del_more:
    mov rdi, fmt_delete
    mov rsi, 10
    xor eax, eax
    call printf

    mov rdi, r12
    mov rsi, 10
    call rb_delete_key

    mov rdi, r12
    call rb_print_inorder

    mov rdi, fmt_verify
    xor eax, eax
    call printf

    mov rdi, r12
    call rb_verify
    test eax, eax
    jz  .invalid3

    mov rdi, fmt_valid
    xor eax, eax
    call printf
    jmp .destroy

.invalid3:
    mov rdi, fmt_invalid
    xor eax, eax
    call printf

    ; === 清理 ===
.destroy:
    mov rdi, fmt_destroy
    xor eax, eax
    call printf

    mov rdi, r12
    call rb_destroy_tree

    xor eax, eax
    pop r12
    pop rbx
    pop rbp
    ret

.exit_fail:
    mov eax, 1
    pop r12
    pop rbx
    pop rbp
    ret
