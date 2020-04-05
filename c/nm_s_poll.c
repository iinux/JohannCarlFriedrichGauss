#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h> /* See NOTES */
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <netdb.h>
#include <errno.h>
#include <unistd.h>
#include <sys/time.h>
#include <sys/stat.h>
#include <fcntl.h>

#include <poll.h>

// https://blog.csdn.net/zhouzhenhe2008/article/details/75807720
// https://blog.csdn.net/superbfly/article/details/72782264
// https://blog.csdn.net/zhuxiaoping54532/article/details/51701549

extern int errno;

#define MaxConnectNum (5)

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
    int domain = AF_INET;
    int type = SOCK_STREAM;
    int protocol = 0;
    int ret = -1;
    int nListenFd = -1;
    int nNewClientFd = -1;
    short int port = 8888;
    struct sockaddr_in addr_in;
    int backlog = 128; // 默认是128
    int len = 0;
    char chBuffer[1024] = {0};
    int flags = 0;
    int nMaxFd = -1;
    int i = 0;
    static int s_nCountClient = 0;
    struct pollfd stuPollFd[MaxConnectNum + 1];
    int ch;
    int n;

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

    nListenFd = socket(domain, type, protocol);
    if (nListenFd < 0)
    {
        printf("\n socket failed ! errno[%d]  err[%s]\n", errno, strerror(errno));
        return -1;
    }

    memset(&addr_in, 0, sizeof(struct sockaddr_in));
    addr_in.sin_family = AF_INET;
    addr_in.sin_port = htons(port); //htons的返回值是16位的网络字节序整型数   htons尾的字母s代表short
    addr_in.sin_addr.s_addr = htonl(INADDR_ANY);

    ret = bind(nListenFd, (struct sockaddr *)(&addr_in), sizeof(struct sockaddr_in));
    if (ret < 0)
    {
        printf("\n bind failed ! errno[%d]  err[%s]\n", errno, strerror(errno));
        close(nListenFd); //避免资源泄漏
        return -1;
    }

    ret = listen(nListenFd, backlog);
    if (ret < 0)
    {
        printf("\n listen failed ! errno[%d]	err[%s]\n", errno, strerror(errno));
        close(nListenFd); //避免资源泄漏
        return -1;
    }

    nMaxFd = 1;
    memset(stuPollFd, 0, sizeof(stuPollFd));

    stuPollFd[0].fd = nListenFd;
    stuPollFd[0].events |= POLLIN;

    for (i = 1; i <= MaxConnectNum; i++)
    {
        stuPollFd[i].fd = -1;
    }

    while (1)
    {
        int time_out_ms = 3000;
        int num = 0;

        num = poll(stuPollFd, nMaxFd, time_out_ms);
        if (num > 0)
        {
            // printf("\n  num =%d\n", num);
            for (i = 1; i <= MaxConnectNum; i++)
            {
                if ((stuPollFd[i].fd != -1) && (POLLIN & stuPollFd[i].revents))
                {
                    char buffer[40] = "";

                    n = read(stuPollFd[i].fd, buffer, sizeof(buffer) - 1);
                    // printf("n=%d %d\n", n, __LINE__);
                    if (n < 0)
                    {
                        perror("read");
                    }
                    else if (n == 0)
                    {
                        struct sockaddr_in sa;
                        socklen_t len = sizeof(sa);
                        if (!getpeername(stuPollFd[i].fd, (struct sockaddr *)&sa, &len))
                        {
                            printf("%s  close fd: %d, %s:%d\n", getTimeStr(), stuPollFd[i].fd, inet_ntoa(sa.sin_addr), ntohs(sa.sin_port));
                        }

                        close(stuPollFd[i].fd);
                        stuPollFd[i].events = 0;
                        s_nCountClient--;
                        stuPollFd[i].fd = -1;
                    }
                    else
                    {
                        n = write(stuPollFd[i].fd, buffer, n);
                        if (n < 0)
                        {
                            perror("write");
                        }
                    }
                    /*
                    len = recv(stuPollFd[i].fd, chBuffer, sizeof(chBuffer), flags); //flags为0，阻塞模式
                    if (len <= 0)
                    {
                        printf("\n recv failed ! errno[%d]	err[%s] len[%d]\n", errno, strerror(errno), len);
                        //	close(nListenFd); //避免资源泄漏
                        close(stuPollFd[i].fd);
                        stuPollFd[i].events = 0;
                        s_nCountClient--;
                        stuPollFd[i].fd = -1;

                        //return -1;
                        continue;
                    }

                    printf("\n i[%d] fd[%d] chBuffer[%s] \n", i, stuPollFd[i].fd, chBuffer);
                    */
                }
            }

            if (POLLIN & stuPollFd[0].revents)
            {
                struct sockaddr_in clnt_addr;
                socklen_t clnt_addr_size = sizeof(clnt_addr);
                nNewClientFd = accept(nListenFd, (struct sockaddr *)&clnt_addr, &clnt_addr_size); //阻塞模式
                if (nNewClientFd < 0)
                {
                    printf("\n accept failed ! errno[%d]	err[%s]\n", errno, strerror(errno));
                    //close(nListenFd); //避免资源泄漏
                    break;
                }

                if (s_nCountClient >= MaxConnectNum)
                {
                    close(nNewClientFd);

                    printf("\n s_nCountClient >= MaxConnectNum \n");
                    continue;
                }

                printf("%s  client fd: %d, %s:%d\n", getTimeStr(), nNewClientFd, inet_ntoa(clnt_addr.sin_addr), ntohs(clnt_addr.sin_port));
                s_nCountClient++;

                for (i = 1; i <= MaxConnectNum; i++)
                {
                    if (stuPollFd[i].fd == -1)
                    {
                        stuPollFd[i].fd = nNewClientFd;
                        stuPollFd[i].events = 0;
                        stuPollFd[i].events |= POLLIN;

                        if (i >= nMaxFd)
                        {
                            nMaxFd = i + 1;
                        }
                        break;
                    }
                }
            }
        }
        else if (num == 0)
        {
        }
        else
        {
            perror("poll");
            break;
        }
    }

    for (i = 1; i <= MaxConnectNum; i++)
    {
        if (stuPollFd[i].fd != -1)
        {
            close(stuPollFd[i].fd);
        }
    }
    close(stuPollFd[0].fd);

    return 0;
}
