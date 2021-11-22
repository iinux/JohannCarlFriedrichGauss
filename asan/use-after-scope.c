#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main()
{
       int *p;
       {
              int num = 10;
              p = &num;
       }
       printf("%d/n", *p);

       return 0;
}
