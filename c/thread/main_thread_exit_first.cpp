#include <pthread.h>
#include <unistd.h>

#include <stdio.h>

void* func(void* arg)
{
    pthread_t main_tid = *static_cast<pthread_t*>(arg);
    pthread_cancel(main_tid);
    while (true)
    {
        //printf("child loops\n");
    }
    return NULL;
}

int main(int argc, char* argv[])
{
    pthread_t main_tid = pthread_self();
    pthread_t tid = 0;
    pthread_create(&tid, NULL, func, &main_tid);
    while (true)
    {
        printf("main loops\n");
    }
    sleep(1);
    printf("main exit\n");
    return 0;
}