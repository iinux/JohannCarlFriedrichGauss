#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// gcc -fsanitize=address -fno-omit-frame-pointer -o heap_ovf_test heap_ovf_test.c
// yum install libasan
// ASAN_OPTIONS='stack_trace_format="[frame=%n, function=%f, location=%S]"'
// ASAN_OPTIONS=help=1
// https://mp.weixin.qq.com/s?__biz=MzI1MTIzMzI2MA==&mid=2650577970&idx=2&sn=7a06b610a31d867f4f5062eddea15305

int main()
{
        char *heap_buf = (char *)malloc(32 * sizeof(char));
        memcpy(heap_buf + 30, "overflow", 8); //在heap_buf的第30个字节开始，拷贝8个字符

        free(heap_buf);

        return 0;
}
