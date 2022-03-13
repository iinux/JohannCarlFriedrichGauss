#include <pthread.h>
#include <unistd.h>

#include <stdio.h>

void* func(void* arg)
{
    while (true)
    {
        sleep(1);
        printf("child loops\n");
    }
    return NULL;
}

int main(int argc, char* argv[])
{
    pthread_t main_tid = pthread_self();
    pthread_t tid = 0;
    pthread_create(&tid, NULL, func, &main_tid);
    sleep(1);
    printf("main exit\n");
    return 0;
}
