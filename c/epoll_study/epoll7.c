#include <stdio.h>
#include <unistd.h>
#include <sys/epoll.h>

int main(void)
{
    int epfd, nfds, i;
    struct epoll_event ev, events[5]; //ev用于注册事件，数组用于返回要处理的事件
    epfd = epoll_create(1);           //只需要监听一个描述符——标准输入
    ev.data.fd = STDOUT_FILENO;
    ev.events = EPOLLOUT | EPOLLET;                     //监听读状态同时设置LT模式
    epoll_ctl(epfd, EPOLL_CTL_ADD, STDOUT_FILENO, &ev); //注册epoll事件
    for (;;)
    {
        nfds = epoll_wait(epfd, events, 5, -1);
        for (i = 0; i < nfds; i++)
        {
            if (events[i].data.fd == STDOUT_FILENO)
            {
                printf("welcome to epoll's word!");
                ev.data.fd = STDOUT_FILENO;
                ev.events = EPOLLOUT | EPOLLET;                     //设置ET模式
                epoll_ctl(epfd, EPOLL_CTL_MOD, STDOUT_FILENO, &ev); //重置epoll事件（ADD无效）
            }
        }
    }
}