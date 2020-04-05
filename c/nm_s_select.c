#include <sys/types.h>
#include <sys/socket.h>
#include <stdio.h>
#include <netinet/in.h>
#include <sys/time.h>
#include <sys/ioctl.h>
#include <unistd.h>
#include <stdlib.h>
#include <string.h>
#include <arpa/inet.h>
#include <time.h>

// https://blog.csdn.net/zhengfl/article/details/21977831
// https://www.cnblogs.com/skyfsm/p/7079458.html

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
    int port = 8888;
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

    int server_sockfd, client_sockfd;
    int server_len;
    socklen_t client_len;
    struct sockaddr_in server_address;
    struct sockaddr_in client_address;
    int result;
    fd_set readfds, testfds;

    server_sockfd = socket(AF_INET, SOCK_STREAM, 0); //建立服务器端socket
    server_address.sin_family = AF_INET;
    server_address.sin_addr.s_addr = htonl(INADDR_ANY);
    server_address.sin_port = htons(port);
    server_len = sizeof(server_address);
    n = bind(server_sockfd, (struct sockaddr *)&server_address, server_len);
    if (n < 0) {
        perror("bind");
        exit(1);
    }
    n = listen(server_sockfd, 20); //监听队列最多容纳5个
    if (n < 0) {
        perror("listen");
        exit(1);
    }
    FD_ZERO(&readfds);
    FD_SET(server_sockfd, &readfds); //将服务器端socket加入到集合中

    while (1)
    {
        char ch;
        int fd;
        int nread;
        testfds = readfds; //将需要监视的描述符集copy到select查询队列中，select会对其修改，所以一定要分开使用变量

        /*无限期阻塞，并测试文件描述符变动 */
        result = select(FD_SETSIZE, &testfds, (fd_set *)0, (fd_set *)0, (struct timeval *)0); //FD_SETSIZE：系统默认的最大文件描述符
        if (result < 1)
        {
            perror("server5");
            exit(1);
        }

        /*扫描所有的文件描述符*/
        for (fd = 0; fd < FD_SETSIZE; fd++)
        {
            /*找到相关文件描述符*/
            if (FD_ISSET(fd, &testfds))
            {
                /*判断是否为服务器套接字，是则表示为客户请求连接。*/
                if (fd == server_sockfd)
                {
                    client_len = sizeof(client_address);
                    client_sockfd = accept(server_sockfd,
                                           (struct sockaddr *)&client_address, &client_len);
                    FD_SET(client_sockfd, &readfds); //将客户端socket加入到集合中
                    struct in_addr in;
                    in.s_addr = client_address.sin_addr.s_addr;
                    printf("%s  adding client fd: %d to select set, %s:%d\n", getTimeStr(), client_sockfd, inet_ntoa(in), ntohs(client_address.sin_port));
                }
                /*客户端socket中有数据请求时*/
                else
                {
                    ioctl(fd, FIONREAD, &nread); //取得数据量交给nread

                    /*客户数据请求完毕，关闭套接字，从集合中清除相应描述符 */
                    if (nread == 0)
                    {
                        struct sockaddr_in sa;
                        socklen_t len = sizeof(sa);
                        if (!getpeername(fd, (struct sockaddr *)&sa, &len))
                        {
                            printf("%s  removing client fd: %d from select set, %s:%d\n", getTimeStr(), fd, inet_ntoa(sa.sin_addr), ntohs(sa.sin_port));
                        }

                        close(fd);
                        FD_CLR(fd, &readfds); //去掉关闭的fd
                    }
                    /*处理客户数据请求*/
                    else
                    {
                        char buffer[40] = "";

                        n = read(fd, buffer, sizeof(buffer) - 1);
                        // printf("n=%d %d\n", n, __LINE__);
                        if (n < 0)
                        {
                            perror("read");
                        }
                        else if (n == 0)
                        {
                        }
                        else
                        {
                            n = write(fd, buffer, n);
                            if (n < 0)
                            {
                                perror("write");
                            }
                        }
                    }
                }
            }
        }
    }

    return 0;
}
