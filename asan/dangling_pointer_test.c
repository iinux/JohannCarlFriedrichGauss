#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main()
{
        char *p = (char *)malloc(32 * sizeof(char));
        free(p);

        int a = p[1];

        return 0;
}
