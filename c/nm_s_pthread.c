#include <stdio.h>
#include <pthread.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <time.h>

struct thread_arg
{
    struct sockaddr_in clnt_addr;
    socklen_t clnt_addr_size;
    int clnt_sock;
    struct in_addr in;
    pthread_t id;
    int ret;
};

char *getTimeStr();

void thread(void *arg)
{
    struct thread_arg *a = (struct thread_arg *)arg;
    int n;
    while (1)
    {
        char buffer[40] = "";
        struct timeval tv;
        tv.tv_sec = 10;
        tv.tv_usec = 0;
        setsockopt(a->clnt_sock, SOL_SOCKET, SO_RCVTIMEO, (const char *)&tv, sizeof tv);
        // printf("%d\n", __LINE__);
        n = read(a->clnt_sock, buffer, sizeof(buffer) - 1);
        // printf("n=%d %d\n", n, __LINE__);
        if (n < 0)
        {
            perror("read");
            break;
        }
        else if (n == 0)
        {
            break;
        }

        n = write(a->clnt_sock, buffer, n);
        if (n < 0)
        {
            perror("write");
            break;
        }
    }

    //关闭套接字
    struct in_addr in;
    in.s_addr = a->clnt_addr.sin_addr.s_addr;
    printf("%s  disconnected: %s:%d\n", getTimeStr(), inet_ntoa(in), ntohs(a->clnt_addr.sin_port));
    close(a->clnt_sock);

    free(a);
}

char *getTimeStr()
{
    time_t timer;
    struct tm *Now;
    time(&timer);
    Now = localtime(&timer);
    return asctime(Now);
}

int main(int argc, char **argv)
{
    char token[] = "hello";
    int port = 8888;
    int n;
    int ch;

    while ((ch = getopt(argc, argv, "p:")) != -1)
    {
        switch (ch)
        {
        case 'p':
            port = atoi(optarg);
            break;
        case '?': // 输入未定义的选项, 都会将该选项的值变为 ?
            printf("unknown option \n");
            break;
        default:
            printf("default \n");
        }
    }

    printf("listen port:%d\n", port);

    //创建套接字
    int serv_sock = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);
    //将套接字和IP、端口绑定
    struct sockaddr_in serv_addr;
    memset(&serv_addr, 0, sizeof(serv_addr));         //每个字节都用0填充
    serv_addr.sin_family = AF_INET;                   //使用IPv4地址
    serv_addr.sin_addr.s_addr = inet_addr("0.0.0.0"); //具体的IP地址
    serv_addr.sin_port = htons(port);                 //端口
    bind(serv_sock, (struct sockaddr *)&serv_addr, sizeof(serv_addr));
    //进入监听状态，等待用户发起请求
    n = listen(serv_sock, 20);
    if (n < 0)
    {
        perror("listen");
        exit(1);
    }

    while (1)
    {
        //接收客户端请求
        struct thread_arg* arg;
        arg = malloc(sizeof(struct thread_arg));
        arg->clnt_addr_size = sizeof(arg->clnt_addr);
        arg->clnt_sock = accept(serv_sock, (struct sockaddr *)&arg->clnt_addr, &arg->clnt_addr_size);

        arg->in.s_addr = arg->clnt_addr.sin_addr.s_addr;
        printf("%s  connected: %s:%d\n", getTimeStr(), inet_ntoa(arg->in), ntohs(arg->clnt_addr.sin_port));

        // 成功返回0，错误返回错误编号
        arg->ret = pthread_create(&arg->id, NULL, (void *)thread, arg);
        if (arg->ret != 0)
        {
            printf("Create pthread error!\n");
            break;
        }

        // pthread_join(id, NULL);
    }

    close(serv_sock);
    return 0;
}
