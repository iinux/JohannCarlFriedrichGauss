#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
#include <arpa/inet.h>

#define GROUP_IP "224.0.1.0"
//#define GROUP_IP "239.0.1.10"

int main()
{
    // 1. 创建通信的套接字
    int fd = socket(AF_INET, SOCK_DGRAM, 0);
    if(fd == -1)
    {
        perror("socket");
        exit(0);
    }

    // 2. 设置组播属性 (经测试可以不设置发送端组播属性也能正常发送)
    struct in_addr opt;
    // 将组播地址初始化到这个结构体成员中即可
    inet_pton(AF_INET, GROUP_IP, &opt.s_addr);
    setsockopt(fd, IPPROTO_IP, IP_MULTICAST_IF, &opt, sizeof(opt));

    char buf[1024];
	char sendaddrbuf[64];
	
	socklen_t len = sizeof(struct sockaddr_in);
	struct sockaddr_in sendaddr;
	
    struct sockaddr_in cliaddr;
    cliaddr.sin_family = AF_INET;
    cliaddr.sin_port = htons(9999); // 接收端需要绑定9999端口
    // 发送组播消息, 需要使用组播地址, 和设置组播属性使用的组播地址一致就可以
    inet_pton(AF_INET, GROUP_IP, &cliaddr.sin_addr.s_addr);
						
    // 3. 通信
    int num = 0;
    while(1)
    {
        if (num > 10) {
            break;
        }
		memset(buf, 0, sizeof(buf));
        sprintf(buf, "hello, client...%d\n", num++);
        // 数据广播
        sendto(fd, buf, strlen(buf)+1, 0, (struct sockaddr*)&cliaddr, len);
        printf("发送的组播的数据: %s\n", buf);
		memset(buf, 0, sizeof(buf));
		recvfrom(fd, buf, sizeof(buf), 0, (struct sockaddr *)&sendaddr, &len);
		printf("sendaddr:%s, port:%d\n", inet_ntop(AF_INET, &sendaddr.sin_addr.s_addr,  sendaddrbuf, sizeof(sendaddrbuf)),  sendaddr.sin_port);
		printf("接收到的组播消息: %s\n", buf);
    }
    close(fd);
    return 0;
}

