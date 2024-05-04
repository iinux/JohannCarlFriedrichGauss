#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <netinet/ip.h>
#include <arpa/inet.h>

int main(int argc, char *argv[]) {
    printf("Number of command line arguments: %d\n", argc);

    // 打印每个命令行参数
    for (int i = 0; i < argc; i++) {
        printf("Argument %d: %s\n", i, argv[i]);
    }

    for (int i = IPPROTO_IP; i <= 255; i++) {
        if (i == 0) continue;
        if (i == 66) continue;
        if (i == 255) continue;
        printf("%d\n", i);
        int sockfd;
        struct sockaddr_in dest_addr;

        // 创建原始套接字
        if ((sockfd = socket(AF_INET, SOCK_RAW, i)) < 0) {
            perror("socket");
            exit(1);
        }

        // 设置目标地址
        memset(&dest_addr, 0, sizeof(dest_addr));
        dest_addr.sin_family = AF_INET;
        if (inet_pton(AF_INET, argv[1], &dest_addr.sin_addr) <= 0) {
            perror("inet_pton");
            exit(1);
        }

        // 构造数据报
        const char *message = "Hello, raw socket!";
        if (sendto(sockfd, message, strlen(message), 0, (struct sockaddr *)&dest_addr, sizeof(dest_addr)) < 0) {
            perror("sendto");
            exit(1);
        }

        printf("Raw packet sent successfully!\n");
    }

    // close(sockfd);
    return 0;
}

