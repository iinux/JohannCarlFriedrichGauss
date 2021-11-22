#include <stdio.h>
#include <stdlib.h>
#include <string.h>

char *get_systeminfo()
{
        char *p_system = (char *)malloc(38 * sizeof(char));
        strcpy(p_system, "Linux version 4.18.0-147.el8.x86_64");
        return p_system;
}

int main()
{
        printf("system info:%s", get_systeminfo());
        return 0;
}
