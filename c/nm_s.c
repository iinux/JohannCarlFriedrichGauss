#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <time.h>
#include <signal.h>
#include <sys/wait.h>

int listen4(int);
int listen6(int);

char *getTimeStr()
{
    time_t timer;
    struct tm *Now;
    time(&timer);
    Now = localtime(&timer);
    return asctime(Now);
}

void sig_chld(int signo)
{
    int stat;
    wait(&stat);
    // printf("wait %d\n", stat);
}

int main(int argc, char **argv)
{
    int port = 8888;
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

    listen6(port);
    return 0;
    int n = fork();
    if (n < 0)
    {
        perror("fork");
        exit(1);
    }
    else if (n == 0)
    {
        listen4(port);
    }
    else
    {
        listen6(port);
    }
}

int listen4(int port)
{
    char token[] = "hello";
    int n;

    //创建套接字
    int serv_sock = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);
    //将套接字和IP、端口绑定
    struct sockaddr_in serv_addr;
    memset(&serv_addr, 0, sizeof(serv_addr));         //每个字节都用0填充
    serv_addr.sin_family = AF_INET;                   //使用IPv4地址
    serv_addr.sin_addr.s_addr = inet_addr("0.0.0.0"); //具体的IP地址
    serv_addr.sin_port = htons(port);                 //端口
    n = bind(serv_sock, (struct sockaddr *)&serv_addr, sizeof(serv_addr));
    if (n < 0)
    {
        perror("bind");
        exit(1);
    }
    //进入监听状态，等待用户发起请求
    n = listen(serv_sock, 20);
    if (n < 0)
    {
        perror("listen");
        exit(1);
    }

    signal(SIGCHLD, &sig_chld);

    while (1)
    {
        //接收客户端请求
        struct sockaddr_in clnt_addr;
        socklen_t clnt_addr_size = sizeof(clnt_addr);
        sleep(100000);
        int clnt_sock = accept(serv_sock, (struct sockaddr *)&clnt_addr, &clnt_addr_size);

        n = fork();
        if (n < 0)
        {
            perror("fork");
            exit(1);
        }
        else if (n == 0)
        {
            close(serv_sock);

            while (1)
            {
                char buffer[40] = "";
                struct timeval tv;
                tv.tv_sec = 10;
                tv.tv_usec = 0;
                setsockopt(clnt_sock, SOL_SOCKET, SO_RCVTIMEO, (const char *)&tv, sizeof tv);
                // printf("%d\n", __LINE__);
                n = read(clnt_sock, buffer, sizeof(buffer) - 1);
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

                n = write(clnt_sock, buffer, n);
                if (n < 0)
                {
                    perror("write");
                    break;
                }
            }

            //关闭套接字
            struct in_addr in;
            in.s_addr = clnt_addr.sin_addr.s_addr;
            printf("%s  pid:%d close %s:%d\n", getTimeStr(), getpid(), inet_ntoa(in), ntohs(clnt_addr.sin_port));
            close(clnt_sock);
            exit(0);
        }
        else
        {
            struct in_addr in;
            in.s_addr = clnt_addr.sin_addr.s_addr;
            printf("%s  %d forked %d, %s:%d\n", getTimeStr(), getpid(), n, inet_ntoa(in), ntohs(clnt_addr.sin_port));

            close(clnt_sock);

            continue;
        }
    }

    close(serv_sock);
    return 0;
}

int listen6(int port)
{
    char token[] = "hello";
    int n;

    //创建套接字
    int serv_sock = socket(AF_INET6, SOCK_STREAM, IPPROTO_TCP);
    //将套接字和IP、端口绑定
    struct sockaddr_in6 serv_addr;
    memset(&serv_addr, 0, sizeof(serv_addr)); //每个字节都用0填充
    serv_addr.sin6_family = AF_INET6;         //使用IPv6地址
    serv_addr.sin6_addr = in6addr_any;        //具体的IP地址
    serv_addr.sin6_port = htons(port);        //端口
    bind(serv_sock, (struct sockaddr *)&serv_addr, sizeof(serv_addr));
    //进入监听状态，等待用户发起请求
    n = listen(serv_sock, 20);
    if (n < 0)
    {
        perror("listen");
        exit(1);
    }

    signal(SIGCHLD, &sig_chld);

    while (1)
    {
        //接收客户端请求
        struct sockaddr_in6 clnt_addr;
        socklen_t clnt_addr_size = sizeof(clnt_addr);
        sleep(100000);
        int clnt_sock = accept(serv_sock, (struct sockaddr *)&clnt_addr, &clnt_addr_size);

        n = fork();
        if (n < 0)
        {
            perror("fork");
            exit(1);
        }
        else if (n == 0)
        {
            close(serv_sock);

            while (1)
            {
                char buffer[40] = "";
                struct timeval tv;
                tv.tv_sec = 10;
                tv.tv_usec = 0;
                setsockopt(clnt_sock, SOL_SOCKET, SO_RCVTIMEO, (const char *)&tv, sizeof tv);
                // printf("%d\n", __LINE__);
                n = read(clnt_sock, buffer, sizeof(buffer) - 1);
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

                n = write(clnt_sock, buffer, n);
                if (n < 0)
                {
                    perror("write");
                    break;
                }
            }

            //关闭套接字
            char buf[64];
            inet_ntop(AF_INET6, (void *)&clnt_addr.sin6_addr, buf, sizeof(buf));
            printf("%s  pid:%d close %s:%d\n", getTimeStr(), getpid(), buf, ntohs(clnt_addr.sin6_port));
            close(clnt_sock);
            exit(0);
        }
        else
        {
            char buf[64];
            inet_ntop(AF_INET6, (void *)&clnt_addr.sin6_addr, buf, sizeof(buf));
            printf("%s  %d forked %d, %s:%d\n", getTimeStr(), getpid(), n, buf, ntohs(clnt_addr.sin6_port));

            close(clnt_sock);

            continue;
        }
    }

    close(serv_sock);
    return 0;
}
