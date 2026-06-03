# Red-Black Tree (x86-64 Linux, NASM)

NASM=nasm
CC=gcc
ASMFLAGS=-f elf64 -g -F dwarf
LDFLAGS=-no-pie

TARGET=rbtree
OBJ=rbtree.o

all: $(TARGET)

$(OBJ): rbtree.asm
	$(NASM) $(ASMFLAGS) $< -o $@

$(TARGET): $(OBJ)
	$(CC) $(LDFLAGS) $< -o $@

run: $(TARGET)
	./$(TARGET)

clean:
	rm -f $(OBJ) $(TARGET)

.PHONY: all run clean
