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

    // 2. 通信的套接字和本地的IP与端口绑定
    struct sockaddr_in addr;
    addr.sin_family = AF_INET;
    addr.sin_port = htons(9999);    // 大端
    addr.sin_addr.s_addr = htonl(INADDR_ANY);  // 0.0.0.0
    int ret = bind(fd, (struct sockaddr*)&addr, sizeof(addr));
    if(ret == -1)
    {
        perror("bind");
        exit(0);
    }

    // 3. 加入到多播组
	#if 0 //使用struct ip_mreqn或者 struct ip_mreq 设置接收端组播属性都可以正常接收
    struct ip_mreqn opt;
    // 要加入到哪个多播组, 通过组播地址来区分
    inet_pton(AF_INET, GROUP_IP, &opt.imr_multiaddr.s_addr);
    opt.imr_address.s_addr = htonl(INADDR_ANY);
    opt.imr_ifindex = if_nametoindex("ens33");
    setsockopt(fd, IPPROTO_IP, IP_ADD_MEMBERSHIP, &opt, sizeof(opt));
	#else
	struct ip_mreq mreq; // 多播地址结构体
	mreq.imr_multiaddr.s_addr=inet_addr(GROUP_IP);
	mreq.imr_interface.s_addr = htonl(INADDR_ANY);	
	ret=setsockopt(fd,IPPROTO_IP,IP_ADD_MEMBERSHIP,&mreq,sizeof(mreq));
	#endif
	
    char buf[1024];
	char sendaddrbuf[64];
	socklen_t len = sizeof(struct sockaddr_in);
	struct sockaddr_in sendaddr;
	
    // 3. 通信
    while(1)
    {
        // 接收广播消息
        memset(buf, 0, sizeof(buf));
        // 阻塞等待数据达到
        recvfrom(fd, buf, sizeof(buf), 0, (struct sockaddr *)&sendaddr, &len);
		printf("sendaddr:%s, port:%d\n", inet_ntop(AF_INET, &sendaddr.sin_addr.s_addr,  sendaddrbuf, sizeof(sendaddrbuf)),  sendaddr.sin_port);
        printf("接收到的组播消息: %s\n", buf);
		sendto(fd, buf, strlen(buf)+1, 0, (struct sockaddr *)&sendaddr, len);
    }
    close(fd);
    return 0;
}

