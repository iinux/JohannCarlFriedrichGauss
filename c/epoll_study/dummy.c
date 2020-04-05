#include <stdio.h>
#include <stdlib.h>
#include <signal.h>
#include <string.h>

// https://www.cnblogs.com/muhe221/articles/4467512.html
// https://blog.csdn.net/YEYUANGEN/article/details/6822004

void dump(int signo)
{
    char buf[1024];
    char cmd[1024];
    FILE *fh;

    snprintf(buf, sizeof(buf), "/proc/%d/cmdline", getpid());
    if (!(fh = fopen(buf, "r")))
        exit(0);
    if (!fgets(buf, sizeof(buf), fh))
        exit(0);
    fclose(fh);
    if (buf[strlen(buf) - 1] == ' ')
        buf[strlen(buf) - 1] = '\0';
    snprintf(cmd, sizeof(cmd), "gdb %s %d", buf, getpid());
    system(cmd);

    exit(0);
}

void dummy_function(void)
{
    unsigned char *ptr = 0x00;
    *ptr = 0x00;
}

int main(void)
{
    signal(SIGSEGV, &dump);
    dummy_function();

    return 0;
}
