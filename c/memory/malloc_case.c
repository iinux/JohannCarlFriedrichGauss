#include <stdio.h>
#include <stdlib.h>

int main() {
    int* ptr = (int*)malloc(sizeof(int));

    if(ptr == NULL) {
        printf("内存分配失败\n");
        return 1;
    }

    *ptr = 10;

    printf("分配的内存地址：%p\n", ptr);
    printf("分配的内存值：%d\n", *ptr);

    free(ptr);

    return 0;
}
