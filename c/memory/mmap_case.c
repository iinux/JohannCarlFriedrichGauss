#include <stdio.h>
#include <stdlib.h>
#include <sys/mman.h>

int main() {
    size_t size = sizeof(int);

    int* ptr = (int*)mmap(NULL, size, PROT_READ | PROT_WRITE, MAP_PRIVATE | MAP_ANONYMOUS, -1, 0);

    if(ptr == MAP_FAILED) {
        printf("内存映射失败\n");
        return 1;
    }

    *ptr = 10;

    printf("分配的内存地址：%p\n", ptr);
    printf("分配的内存值：%d\n", *ptr);

    munmap(ptr, size);

    return 0;
}
