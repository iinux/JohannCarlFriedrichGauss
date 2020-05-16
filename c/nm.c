#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <time.h>
// #include <sys/timeb.h>
#include <sys/time.h>

#include <errno.h>
#include <ctype.h>
#include <netdb.h>

// https://blog.csdn.net/witxjp/article/details/8079751
// http://blog.zhangjikai.com/2016/03/05/%E3%80%90C%E3%80%91%E8%A7%A3%E6%9E%90%E5%91%BD%E4%BB%A4%E8%A1%8C%E5%8F%82%E6%95%B0--getopt%E5%92%8Cgetopt_long/

/*
long long getSystemTime()
{
    struct timeb t;
    ftime(&t);
    return 1000 * t.time + t.millitm;
}
*/

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
    char *host = "127.0.0.1";
    int port = 8888;
    int interval = 1;
    int isv6 = 0;

    int ch;
    while ((ch = getopt(argc, argv, "h:p:i:m:6")) != -1)
    {
        switch (ch)
        {
        case 'h':
            host = optarg;
            break;
        case 'p':
            port = atoi(optarg);
            break;
        case 'i':
            interval = atoi(optarg);
            break;
        case 'm':
            break;
        case '6':
            isv6 = 1;
            break;
        case '?': // 输入未定义的选项, 都会将该选项的值变为 ?
            printf("unknown option \n");
            break;
        default:
            printf("default \n");
        }
    }

    printf("host:%s port:%d interval:%d\n", host, port, interval);

    struct sockaddr_in serv_addr;
    memset(&serv_addr, 0, sizeof(serv_addr)); //每个字节都用0填充
    struct sockaddr_in6 serv_addr6;
    memset(&serv_addr6, 0, sizeof(serv_addr6)); //每个字节都用0填充

    int sock;
    int n;
    if (isv6 == 1)
    {
        struct addrinfo hints, *servinfo, *p;

        memset(&hints, 0, sizeof hints);
        // hints.ai_family = AF_UNSPEC; // 設定 AF_INET6 表示強迫使用 IPv6
        hints.ai_family = AF_INET6;
        hints.ai_socktype = SOCK_STREAM;
        char buf[8];
        sprintf((char *)&buf, "%d", port);
        n = getaddrinfo(host, buf, &hints, &servinfo);
        if (n != 0)
        {
            fprintf(stderr, "getaddrinfo: %s\n", gai_strerror(n));
            exit(1);
        }
        for (p = servinfo; p != NULL; p = p->ai_next)
        {
            char buf[64];
            struct sockaddr_in6 *clnt_addr = (struct sockaddr_in6 *)p->ai_addr;
            inet_ntop(p->ai_family, &clnt_addr->sin6_addr, buf, sizeof(buf));
            printf("name:%s, length:%d, addrtype:%d, ip:%s\n", p->ai_canonname, p->ai_addrlen, p->ai_family, buf);

            sock = socket(p->ai_family, p->ai_socktype, p->ai_protocol);
            if (sock == -1)
            {
                perror("socket");
                continue;
            }

            n = connect(sock, p->ai_addr, p->ai_addrlen);
            if (n == -1)
            {
                close(sock);
                perror("connect");
                continue;
            }
            // printf("%d\n", __LINE__);

            break; // if we get here, we must have connected successfully
        }
    }
    else
    {
        extern int h_errno;
        struct hostent *h;
        struct in_addr in;

        h = gethostbyname(host);
        // printf("%d\n",__LINE__);
        if (h == NULL)
        {
            perror("gethostbyname");
            // printf("gethostbyname %s\n", hstrerror(h_errno));
            exit(1);
        }

        memcpy(&serv_addr.sin_addr.s_addr, h->h_addr, 4);
        in.s_addr = serv_addr.sin_addr.s_addr;
        printf("name:%s, length:%d, addrtype:%d, ip:%s\n", h->h_name, h->h_length, h->h_addrtype, inet_ntoa(in));

        //创建套接字
        sock = socket(AF_INET, SOCK_STREAM, 0);
        //向服务器（特定的IP和端口）发起请求
        serv_addr.sin_family = AF_INET; //使用IPv4地址
        // serv_addr.sin_addr.s_addr = inet_addr(host);
        serv_addr.sin_port = htons(port); //端口
        n = connect(sock, (struct sockaddr *)&serv_addr, sizeof(serv_addr));
    }

    if (n < 0)
    {
        perror("connect");
        exit(1);
    }

    while (1)
    {
        sleep(interval);
        // printf("%d\n", __LINE__);

        /*
        // sys/timeb.h
        long long start = getSystemTime();
        */

        // sys/time.h
        struct timeval start2, end2;
        gettimeofday(&start2, NULL);

        char buffer[40] = "";
        n = write(sock, token, sizeof(token) - 1);
        if (n < 0)
        {
            perror("write");
            break;
        }
        // printf("%d\n", __LINE__);
        n = read(sock, buffer, sizeof(buffer) - 1);
        if (n < 0)
        {
            perror("read");
            break;
        }
        // printf("%d\n", __LINE__);
        // printf("Message form server: %s\n", buffer);

        /*
        long long end = getSystemTime();
        printf("%s  time: %lld ms\n", getTimeStr(), end - start);
        */

        gettimeofday(&end2, NULL);
        int timeuse = 1000000 * (end2.tv_sec - start2.tv_sec) + end2.tv_usec - start2.tv_usec;
        printf("%s  time: %d ms\n", getTimeStr(), timeuse / 1000);
    }

    //关闭套接字
    n = close(sock);
    if (n < 0)
    {
        perror("close");
        exit(1);
    }
    return 0;
}
