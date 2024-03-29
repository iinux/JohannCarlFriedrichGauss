#include <sys/ptrace.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>
#include <linux/user.h>   /* For constants ORIG_EAX etc */

// https://www.cnblogs.com/mysky007/p/11047943.html
// https://www.linuxjournal.com/article/6100?page=0,0
// https://www.cnblogs.com/tangr206/articles/3094358.html
//

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

        orig_eax = ptrace(PTRACE_PEEKUSER, child, 4 * ORIG_EAX, NULL);
        printf("The child made a system call %ldn", orig_eax);

        ptrace(PTRACE_CONT, child, NULL, NULL);
    }

    return 0;
}
