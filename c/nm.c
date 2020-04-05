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

    int ch;
    while ((ch = getopt(argc, argv, "h:p:i:m:")) != -1)
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
        case '?': // 输入未定义的选项, 都会将该选项的值变为 ?
            printf("unknown option \n");
            break;
        default:
            printf("default \n");
        }
    }

    printf("host:%s port:%d interval:%d\n", host, port, interval);

    extern int h_errno;
    struct hostent *h;
    struct in_addr in;
    struct sockaddr_in serv_addr;
    memset(&serv_addr, 0, sizeof(serv_addr)); //每个字节都用0填充
    h = gethostbyname(host);
    if (h == NULL)
    {
        printf("%s\n", hstrerror(h_errno));
        exit(1);
    }
    else
    {
        memcpy(&serv_addr.sin_addr.s_addr, h->h_addr, 4);
        in.s_addr = serv_addr.sin_addr.s_addr;
        printf("name:%s, length:%d, addrtype:%d, ip:%s\n", h->h_name, h->h_length, h->h_addrtype, inet_ntoa(in));
    }

    //创建套接字
    int sock = socket(AF_INET, SOCK_STREAM, 0);
    //向服务器（特定的IP和端口）发起请求
    serv_addr.sin_family = AF_INET; //使用IPv4地址
    // serv_addr.sin_addr.s_addr = inet_addr(host);
    serv_addr.sin_port = htons(port); //端口
    int n;
    n = connect(sock, (struct sockaddr *)&serv_addr, sizeof(serv_addr));
    if (n < 0)
    {
        perror("connect");
        exit(1);
    }

    while (1)
    {
        sleep(interval);

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
        n = read(sock, buffer, sizeof(buffer) - 1);
        if (n < 0)
        {
            perror("read");
            break;
        }
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
