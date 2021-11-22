#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// 运行时，启用ASAN_OPTIONS=detect_stack_use_after_return=1标志，才能检测此种内存错误使用的情况。

int *ptr;
void get_pointer()
{
       int local[10];
       ptr = &local[0];
       return;
}

int main()
{
       get_pointer();

       printf("%d\n", *ptr);

       return 0;
}
