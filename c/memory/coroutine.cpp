#include <stdio.h>
#include <stdlib.h>

#define STACK_SIZE 1024

typedef void(*coro_start)();

class coroutine {
public:
    long* stack_pointer;
    char* stack;

    coroutine(coro_start entry) {
        if (entry == NULL) {
            stack = NULL;
            stack_pointer = NULL;
            return;
        }

        stack = (char*)malloc(STACK_SIZE);
        char* base = stack + STACK_SIZE;
        stack_pointer = (long*) base;
        stack_pointer -= 1;
        *stack_pointer = (long) entry;
        stack_pointer -= 1;
        *stack_pointer = (long) base;
    }

    ~coroutine() {
        if (!stack)
            return;
        free(stack);
        stack = NULL;
    }
};

coroutine* co_a, * co_b;

void yield_to(coroutine* old_co, coroutine* co) {
    __asm__ (
        "movq %%rsp, %0\n\t"
        "movq %%rax, %%rsp\n\t"
        :"=m"(old_co->stack_pointer):"a"(co->stack_pointer):);
}

void start_b() {
    printf("B");
    yield_to(co_b, co_a);
    printf("D");
    yield_to(co_b, co_a);
}

int main() {
    printf("A");
    co_b = new coroutine(start_b);
    co_a = new coroutine(NULL);
    yield_to(co_a, co_b);
    printf("C");
    yield_to(co_a, co_b);
    printf("E\n");
    delete co_a;
    delete co_b;
    return 0;
}

// g++ -g -o co -O0 coroutine.cpp
