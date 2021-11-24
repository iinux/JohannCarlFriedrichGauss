#include <sys/ptrace.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>
#include <sys/user.h>

int main()
{
    pid_t child;
    long orig_eax;
    child = fork();
    if(child == 0) {
        ptrace(PTRACE_TRACEME, 0, NULL, NULL);
        execl("/bin/ls", "ls", NULL);
    } else {
        wait(NULL);

        struct user_regs_struct regs;
        ptrace(PTRACE_GETREGS, child, NULL, &regs);
        printf("The child made a system call %ldn", regs.orig_rax);

        ptrace(PTRACE_CONT, child, NULL, NULL);
    }
    return 0;
}
