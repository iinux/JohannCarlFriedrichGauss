#include <stdio.h>
#include <unistd.h>
#include <string.h>
#include <stdlib.h>
#include <netinet/in.h>
#include <time.h>
#include <arpa/inet.h>
#include <sys/epoll.h>

#define MAXSIZE 1024
#define LISTENQ 5
#define FDSIZE 1000
#define EPOLLEVENTS 100

// https://www.cnblogs.com/lojunren/p/3856290.html
// https://segmentfault.com/a/1190000003063859

char *getTimeStr()
{
    time_t timer;
    struct tm *Now;
    time(&timer);
    Now = localtime(&timer);
    return asctime(Now);
}

//添加事件
static void add_event(int epollfd, int fd, int state)
{
    struct epoll_event ev;
    ev.events = state;
    ev.data.fd = fd;
    epoll_ctl(epollfd, EPOLL_CTL_ADD, fd, &ev);
}

//删除事件
static void delete_event(int epollfd, int fd, int state)
{
    struct epoll_event ev;
    ev.events = state;
    ev.data.fd = fd;
    epoll_ctl(epollfd, EPOLL_CTL_DEL, fd, &ev);
}

//修改事件
static void modify_event(int epollfd, int fd, int state)
{
    struct epoll_event ev;
    ev.events = state;
    ev.data.fd = fd;
    epoll_ctl(epollfd, EPOLL_CTL_MOD, fd, &ev);
}

//处理接收到的连接
static void handle_accpet(int epollfd, int listenfd)
{
    int clifd;
    struct sockaddr_in cliaddr;
    socklen_t cliaddrlen = sizeof(cliaddr);
    clifd = accept(listenfd, (struct sockaddr *)&cliaddr, &cliaddrlen);
    if (clifd == -1)
        perror("accpet error:");
    else
    {
        printf("%s  accept a new client: %s:%d\n", getTimeStr(), inet_ntoa(cliaddr.sin_addr), ntohs(cliaddr.sin_port)); //添加一个客户描述符和事件
        add_event(epollfd, clifd, EPOLLIN);
    }
}

//读处理
static void do_read(int epollfd, int fd, char *buf)
{
    int nread;
    nread = read(fd, buf, MAXSIZE);
    if (nread == -1)
    {
        perror("read error:");
        close(fd);                          //记住close fd
        delete_event(epollfd, fd, EPOLLIN); //删除监听
    }
    else if (nread == 0)
    {
        struct sockaddr_in sa;
        socklen_t len = sizeof(sa);
        if (!getpeername(fd, (struct sockaddr *)&sa, &len))
        {
            printf("%s  close client fd: %d %s:%d\n", getTimeStr(), fd, inet_ntoa(sa.sin_addr), ntohs(sa.sin_port));
        }
        close(fd);                          //记住close fd
        delete_event(epollfd, fd, EPOLLIN); //删除监听
    }
    else
    {
        // printf("read message is : %s", buf);
        //修改描述符对应的事件，由读改为写
        modify_event(epollfd, fd, EPOLLOUT);
    }
}

//写处理
static void do_write(int epollfd, int fd, char *buf)
{
    int nwrite;
    nwrite = write(fd, buf, strlen(buf));
    if (nwrite == -1)
    {
        perror("write error:");
        close(fd);                           //记住close fd
        delete_event(epollfd, fd, EPOLLOUT); //删除监听
    }
    else
    {
        modify_event(epollfd, fd, EPOLLIN);
    }
    // memset(buf, 0, MAXSIZE);
}

//事件处理函数
static void handle_events(int epollfd, struct epoll_event *events, int num, int listenfd, char *buf)
{
    int i;
    int fd;
    //进行遍历;这里只要遍历已经准备好的io事件。num并不是当初epoll_create时的FDSIZE。
    for (i = 0; i < num; i++)
    {
        fd = events[i].data.fd;
        //根据描述符的类型和事件类型进行处理
        if ((fd == listenfd) && (events[i].events & EPOLLIN))
            handle_accpet(epollfd, listenfd);
        else if (events[i].events & EPOLLIN)
            do_read(epollfd, fd, buf);
        else if (events[i].events & EPOLLOUT)
            do_write(epollfd, fd, buf);
    }
}

int main(int argc, char **argv)
{
    int listenfd, epollfd, ret;
    char buf[MAXSIZE];
    int port = 8888;
    int n;
    int ch;

    memset(buf, 0, MAXSIZE);

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

    struct epoll_event events[EPOLLEVENTS];

    //创建一个描述符
    epollfd = epoll_create(FDSIZE);

    //添加监听描述符事件
    add_event(epollfd, serv_sock, EPOLLIN);

    //循环等待
    for (;;)
    {
        //该函数返回已经准备好的描述符事件数目
        ret = epoll_wait(epollfd, events, EPOLLEVENTS, -1);
        //处理接收到的连接
        handle_events(epollfd, events, ret, serv_sock, buf);
    }
}