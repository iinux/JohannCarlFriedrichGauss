#include <signal.h>           /* Definition of SIG* constants */
#include <sys/syscall.h>      /* Definition of SYS_* constants */
#include <unistd.h>

int main() {
    syscall(SYS_tkill, 3300511,15);
}
