#include <sys/wait.h>
#include <unistd.h>     /* For fork() */
#include <sys/ptrace.h>
#include <sys/wait.h>
#include <sys/reg.h>   /* For constants ORIG_RAX etc */
#include <sys/user.h>
#include <sys/syscall.h> /* SYS_write */
#include <stdio.h>
int main() {
    pid_t child;
    long orig_rax;
    int status;
    int iscalling = 0;
    struct user_regs_struct regs;
 
    child = fork();
    if(child == 0)
    {
        ptrace(PTRACE_TRACEME, 0, NULL, NULL);
        execl("/bin/ls", "ls", "-l", "-h", NULL);
    }
    else
    {
        while(1)
        {
            wait(&status);
            if(WIFEXITED(status))
            {
                break;
            }
            orig_rax = ptrace(PTRACE_PEEKUSER,
                              child, 8 * ORIG_RAX,
                              NULL);
            if(orig_rax == SYS_write)
            {
                ptrace(PTRACE_GETREGS, child, NULL, &regs);         //获取寄存器参数
                if(!iscalling)          //进入系统调用
                {
                    iscalling = 1;         
                    printf("[Enter SYS_write call] with regs.rdi [%ld], regs.rsi[%ld], regs.rdx[%ld], regs.rax[%ld], regs.orig_rax[%ld]\n",
                            regs.rdi, regs.rsi, regs.rdx,regs.rax, regs.orig_rax);
                }
                else            //离开此次系统调用
                {
                    printf("[Leave SYS_write call] return regs.rax [%ld], regs.orig_rax [%ld]\n", regs.rax, regs.orig_rax);
                    iscalling = 0;
                }
            }
            ptrace(PTRACE_SYSCALL, child, NULL, NULL);
        }
    }
    return 0;
}
