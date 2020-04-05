#include <stdio.h>
#include <unistd.h>
#include <sys/epoll.h>

int main(void)
{
    int epfd, nfds, i;
    struct epoll_event ev, events[5]; //ev用于注册事件，数组用于返回要处理的事件
    epfd = epoll_create(1);           //只需要监听一个描述符——标准输入
    ev.data.fd = STDOUT_FILENO;
    ev.events = EPOLLOUT | EPOLLET;                     //监听读状态同时设置ET模式
    epoll_ctl(epfd, EPOLL_CTL_ADD, STDOUT_FILENO, &ev); //注册epoll事件
    for (;;)
    {
        nfds = epoll_wait(epfd, events, 5, -1);
        for (i = 0; i < nfds; i++)
        {
            if (events[i].data.fd == STDOUT_FILENO)
            {
                printf("welcome to epoll's word!"); // 只是将输出语句的printf的换行符移除。我们看到程序成挂起状态。
                // 因为第一次epoll_wait返回写就绪后，程序向标准输出的buffer中写入“welcome to epoll's world！”，
                // 但是因为没有输出换行，所以buffer中的内容一直存在，下次epoll_wait的时候，虽然有写空间但是ET模式下不再返回写就绪。
                // 回忆第一节关于ET的实现，这种情况原因就是第一次buffer为空，导致epitem加入rdlist，返回一次就绪后移除此epitem，
                // 之后虽然buffer仍然可写，但是由于对应epitem已经不再rdlist中，就不会对其就绪fd的events的在检测了。
            }
        }
    }
}