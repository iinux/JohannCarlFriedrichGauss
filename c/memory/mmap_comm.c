#include <sys/mman.h>
#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
int main() {
    pid_t pid;
    char* shm = (char*)mmap(0, 4096, PROT_READ | PROT_WRITE, MAP_SHARED | MAP_ANONYMOUS, -1, 0);
    if (!(pid = fork())){
        sleep(1);
        printf("child got a message: %s\n", shm);
        sprintf(shm, "%s", "hello, father.");
        exit(0);
    }
    sprintf(shm, "%s", "hello, my child");
    sleep(2);
    printf("parent got a message: %s\n", shm);
    return 0;
}
