#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define PAGE_SZ (1<<12)

int main() {
    int i;
    int gb = 1; //以GB为单位分配内存大小

    for (i = 0; i < ((unsigned long)gb<<30)/PAGE_SZ ; ++i) {
        void *m = malloc(PAGE_SZ);
        if (!m)
            break;
        memset(m, 0, 1);
    }
    printf("allocated %lu MB\n", ((unsigned long)i*PAGE_SZ)>>20);
    getchar();
    return 0;
}

/*

oom
https://blog.51cto.com/laoxu/1267097

pgrep -f "/usr/sbin/sshd" | while read PID;do echo -17 > /proc/$PID/oom_adj;done

echo -17 > /proc/$(pidof sshd)/oom_adj

/include/uapi/linux/oom.h

# sysctl -w vm.panic_on_oom=1
vm.panic_on_oom = 1   //1表示关闭，默认为0表示开启OOM
# sysctl -p

top -e m

echo f > /proc/sysrq-trigger   // 'f' - Will call oom_kill to kill a memory hog process.

*/