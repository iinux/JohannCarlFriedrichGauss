// Red-Black Tree in ARM64 macOS assembly
// Build: clang -arch arm64 -o rbtree rbtree.s
//
// Node layout (40 bytes):
//   [+0]  key    (i64)
//   [+8]  color  (i64)  0 = RED, 1 = BLACK
//   [+16] left   (ptr)
//   [+24] right  (ptr)
//   [+32] parent (ptr)
//
// CLRS-style sentinel NIL is used so every node always has valid children/parent.

.equ KEY,    0
.equ COLOR,  8
.equ LEFT,   16
.equ RIGHT,  24
.equ PARENT, 32
.equ RED,    0
.equ BLACK,  1
.equ NSIZE,  40

.section __DATA,__data
.align 3
root:    .quad 0
nilnode: .space NSIZE

.section __TEXT,__cstring,cstring_literals
Lfmt:    .asciz "%ld "
Lnl:     .asciz "\n"

.section __TEXT,__text
.align 2

// ---------- helpers to load globals ----------
// load_root  -> x0 = root
// load_nil   -> x0 = &nilnode

.macro LOAD_ROOT reg
    adrp \reg, root@PAGE
    add  \reg, \reg, root@PAGEOFF
.endm

.macro LOAD_NIL reg
    adrp \reg, nilnode@PAGE
    add  \reg, \reg, nilnode@PAGEOFF
.endm

// ---------- init ----------
// Initializes NIL sentinel and sets root = NIL
.globl _rb_init
_rb_init:
    LOAD_NIL x0
    mov  x1, #BLACK
    str  x1, [x0, #COLOR]
    str  xzr, [x0, #KEY]
    str  x0, [x0, #LEFT]
    str  x0, [x0, #RIGHT]
    str  x0, [x0, #PARENT]
    LOAD_ROOT x1
    str  x0, [x1]
    ret

// ---------- new_node(key) -> x0 ----------
// allocates a node, color = RED, children/parent = NIL
_new_node:
    stp  x29, x30, [sp, #-32]!
    mov  x29, sp
    str  x19, [sp, #16]
    mov  x19, x0
    mov  x0, #NSIZE
    bl   _malloc
    str  x19, [x0, #KEY]
    mov  x2, #RED
    str  x2, [x0, #COLOR]
    LOAD_NIL x3
    str  x3, [x0, #LEFT]
    str  x3, [x0, #RIGHT]
    str  x3, [x0, #PARENT]
    ldr  x19, [sp, #16]
    ldp  x29, x30, [sp], #32
    ret

// ---------- left_rotate(x) ----------
// x in x0
// y = x.right; x.right = y.left; if y.left != NIL: y.left.parent = x
// y.parent = x.parent; fix x.parent's child link; y.left = x; x.parent = y
_left_rotate:
    // x = x0
    ldr  x1, [x0, #RIGHT]          // y = x.right
    ldr  x2, [x1, #LEFT]            // x.right = y.left
    str  x2, [x0, #RIGHT]
    LOAD_NIL x3
    cmp  x2, x3
    b.eq 1f
    str  x0, [x2, #PARENT]          // y.left.parent = x
1:
    ldr  x4, [x0, #PARENT]          // y.parent = x.parent
    str  x4, [x1, #PARENT]
    cmp  x4, x3                     // x.parent == NIL?
    b.ne 2f
    LOAD_ROOT x5
    str  x1, [x5]                   // root = y
    b    4f
2:
    ldr  x6, [x4, #LEFT]
    cmp  x0, x6                     // x == x.parent.left?
    b.ne 3f
    str  x1, [x4, #LEFT]
    b    4f
3:
    str  x1, [x4, #RIGHT]
4:
    str  x0, [x1, #LEFT]            // y.left = x
    str  x1, [x0, #PARENT]          // x.parent = y
    ret

// ---------- right_rotate(x) ----------
// mirror of left_rotate
_right_rotate:
    ldr  x1, [x0, #LEFT]            // y = x.left
    ldr  x2, [x1, #RIGHT]           // x.left = y.right
    str  x2, [x0, #LEFT]
    LOAD_NIL x3
    cmp  x2, x3
    b.eq 1f
    str  x0, [x2, #PARENT]
1:
    ldr  x4, [x0, #PARENT]
    str  x4, [x1, #PARENT]
    cmp  x4, x3
    b.ne 2f
    LOAD_ROOT x5
    str  x1, [x5]
    b    4f
2:
    ldr  x6, [x4, #RIGHT]
    cmp  x0, x6
    b.ne 3f
    str  x1, [x4, #RIGHT]
    b    4f
3:
    str  x1, [x4, #LEFT]
4:
    str  x0, [x1, #RIGHT]
    str  x1, [x0, #PARENT]
    ret

// ---------- rb_insert_fixup(z) ----------
// z in x0
// Walks up the tree fixing red-red violations.
_rb_insert_fixup:
    stp  x29, x30, [sp, #-32]!
    mov  x29, sp
    str  x19, [sp, #16]
    mov  x19, x0                    // x19 = z

Lfix_loop:
    ldr  x1, [x19, #PARENT]         // p = z.parent
    ldr  x2, [x1, #COLOR]
    cmp  x2, #RED
    b.ne Lfix_done                  // while p.color == RED

    ldr  x3, [x1, #PARENT]          // gp = p.parent
    ldr  x4, [x3, #LEFT]
    cmp  x1, x4
    b.ne Lfix_p_right

    // ----- p is gp.left -----
    ldr  x5, [x3, #RIGHT]           // y = gp.right (uncle)
    ldr  x6, [x5, #COLOR]
    cmp  x6, #RED
    b.ne Lcase2_left
    // Case 1: uncle red — recolor and move up
    mov  x7, #BLACK
    str  x7, [x1, #COLOR]
    str  x7, [x5, #COLOR]
    mov  x7, #RED
    str  x7, [x3, #COLOR]
    mov  x19, x3
    b    Lfix_loop
Lcase2_left:
    ldr  x6, [x1, #RIGHT]
    cmp  x19, x6                    // z == p.right?
    b.ne Lcase3_left
    // Case 2: rotate left around p, then z = p
    mov  x19, x1
    mov  x0, x19
    bl   _left_rotate
Lcase3_left:
    // Case 3: recolor p,gp and right_rotate(gp)
    ldr  x1, [x19, #PARENT]
    ldr  x3, [x1, #PARENT]
    mov  x7, #BLACK
    str  x7, [x1, #COLOR]
    mov  x7, #RED
    str  x7, [x3, #COLOR]
    mov  x0, x3
    bl   _right_rotate
    b    Lfix_loop

Lfix_p_right:
    // ----- p is gp.right (mirror) -----
    ldr  x5, [x3, #LEFT]            // y = gp.left
    ldr  x6, [x5, #COLOR]
    cmp  x6, #RED
    b.ne Lcase2_right
    mov  x7, #BLACK
    str  x7, [x1, #COLOR]
    str  x7, [x5, #COLOR]
    mov  x7, #RED
    str  x7, [x3, #COLOR]
    mov  x19, x3
    b    Lfix_loop
Lcase2_right:
    ldr  x6, [x1, #LEFT]
    cmp  x19, x6                    // z == p.left?
    b.ne Lcase3_right
    mov  x19, x1
    mov  x0, x19
    bl   _right_rotate
Lcase3_right:
    ldr  x1, [x19, #PARENT]
    ldr  x3, [x1, #PARENT]
    mov  x7, #BLACK
    str  x7, [x1, #COLOR]
    mov  x7, #RED
    str  x7, [x3, #COLOR]
    mov  x0, x3
    bl   _left_rotate
    b    Lfix_loop

Lfix_done:
    LOAD_ROOT x0
    ldr  x1, [x0]
    mov  x2, #BLACK
    str  x2, [x1, #COLOR]           // root.color = BLACK
    ldr  x19, [sp, #16]
    ldp  x29, x30, [sp], #32
    ret

// ---------- rb_insert(key) ----------
// x0 = key
.globl _rb_insert
_rb_insert:
    stp  x29, x30, [sp, #-48]!
    mov  x29, sp
    stp  x19, x20, [sp, #16]
    stp  x21, x22, [sp, #32]

    bl   _new_node                  // x0 = z
    mov  x19, x0                    // x19 = z
    ldr  x22, [x19, #KEY]           // x22 = key (cache)

    LOAD_NIL x21                    // x21 = NIL
    LOAD_ROOT x20
    ldr  x2, [x20]                  // x = root
    mov  x3, x21                    // y = NIL

Lins_loop:
    cmp  x2, x21
    b.eq Lins_done_walk
    mov  x3, x2                     // y = x
    ldr  x4, [x2, #KEY]
    cmp  x22, x4
    b.ge Lins_right
    ldr  x2, [x2, #LEFT]
    b    Lins_loop
Lins_right:
    ldr  x2, [x2, #RIGHT]
    b    Lins_loop

Lins_done_walk:
    str  x3, [x19, #PARENT]         // z.parent = y
    cmp  x3, x21
    b.ne 1f
    str  x19, [x20]                 // root = z
    b    2f
1:
    ldr  x4, [x3, #KEY]
    cmp  x22, x4
    b.ge 11f
    str  x19, [x3, #LEFT]
    b    2f
11:
    str  x19, [x3, #RIGHT]
2:
    // children & color already set by _new_node
    mov  x0, x19
    bl   _rb_insert_fixup

    ldp  x21, x22, [sp, #32]
    ldp  x19, x20, [sp, #16]
    ldp  x29, x30, [sp], #48
    ret

// ---------- inorder(node) ----------
// recursive: print left, self, right
_inorder:
    stp  x29, x30, [sp, #-32]!
    mov  x29, sp
    str  x19, [sp, #16]
    mov  x19, x0
    LOAD_NIL x1
    cmp  x19, x1
    b.eq Lin_ret

    ldr  x0, [x19, #LEFT]
    bl   _inorder

    // macOS arm64: variadic args go on the stack, not registers
    ldr  x1, [x19, #KEY]
    sub  sp, sp, #16
    str  x1, [sp]
    adrp x0, Lfmt@PAGE
    add  x0, x0, Lfmt@PAGEOFF
    bl   _printf
    add  sp, sp, #16

    ldr  x0, [x19, #RIGHT]
    bl   _inorder
Lin_ret:
    ldr  x19, [sp, #16]
    ldp  x29, x30, [sp], #32
    ret

// ---------- rb_print() ----------
.globl _rb_print
_rb_print:
    stp  x29, x30, [sp, #-16]!
    mov  x29, sp
    LOAD_ROOT x0
    ldr  x0, [x0]
    bl   _inorder
    adrp x0, Lnl@PAGE
    add  x0, x0, Lnl@PAGEOFF
    bl   _printf
    ldp  x29, x30, [sp], #16
    ret

// ---------- main ----------
// Inserts a sequence of keys, then prints them in sorted order.
.globl _main
_main:
    stp  x29, x30, [sp, #-32]!
    mov  x29, sp
    str  x19, [sp, #16]

    bl   _rb_init

    // keys: 20 15 25 10 5 1 30 35 40 50 45 7 12 17
    adrp x19, Lkeys@PAGE
    add  x19, x19, Lkeys@PAGEOFF
Lmain_loop:
    ldr  x0, [x19], #8
    cbz  x0, Lmain_done
    bl   _rb_insert
    b    Lmain_loop
Lmain_done:
    bl   _rb_print
    mov  w0, #0
    ldr  x19, [sp, #16]
    ldp  x29, x30, [sp], #32
    ret

.section __DATA,__data
.align 3
Lkeys:
    .quad 20, 15, 25, 10, 5, 1, 30, 35, 40, 50, 45, 7, 12, 17
    .quad 0
