#include <seccomp.h>
#include <sys/socket.h>
#include <errno.h>
#include <stdio.h>
#include <unistd.h>

int main(int argc, char *argv[]) {
    scmp_filter_ctx ctx = seccomp_init(SCMP_ACT_ALLOW);
    seccomp_rule_add(ctx, SCMP_ACT_ERRNO(EAFNOSUPPORT), SCMP_SYS(socket), 1, SCMP_A0(SCMP_CMP_EQ, AF_INET6));
    seccomp_load(ctx);
    seccomp_release(ctx);
    execvp(argv[1], &argv[1]);
    return 1;
}
