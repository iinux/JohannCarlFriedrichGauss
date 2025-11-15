#include <netdb.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <ifaddrs.h>
#include <net/if.h>

#define SERVER_HOST "va.google.party"
#define SERVER_PORT 5800
#define BUFFER_SIZE 1024

int get_local_ipv6_addresses(char *buffer, size_t buffer_size) {
    struct ifaddrs *ifaddrs_ptr, *ifa;
    char ip_str[INET6_ADDRSTRLEN];
    int offset = 0;
    
    if (getifaddrs(&ifaddrs_ptr) == -1) {
        perror("getifaddrs");
        return -1;
    }
    
    //offset += snprintf(buffer + offset, buffer_size - offset, "Local IPv6 addresses:\n");
    
    for (ifa = ifaddrs_ptr; ifa != NULL; ifa = ifa->ifa_next) {
        if (ifa->ifa_addr == NULL)
            continue;
            
        // 只处理IPv6地址且不是回环地址
        if (ifa->ifa_addr->sa_family == AF_INET6) {
            struct sockaddr_in6 *sin6 = (struct sockaddr_in6 *)ifa->ifa_addr;
            
            // 跳过回环地址
            if (IN6_IS_ADDR_LOOPBACK(&(sin6->sin6_addr)))
                continue;
                
            // 跳过链路本地地址
            if (IN6_IS_ADDR_LINKLOCAL(&(sin6->sin6_addr)))
                continue;
            
            inet_ntop(AF_INET6, &(sin6->sin6_addr), ip_str, INET6_ADDRSTRLEN);
            if (strncmp(ip_str, "24", 2) == 0) {
                //offset += snprintf(buffer + offset, buffer_size - offset, "%s: %s\n", ifa->ifa_name, ip_str);
                offset += snprintf(buffer + offset, buffer_size - offset, "%s", ip_str);
                break;
            }
        }
    }
    
    freeifaddrs(ifaddrs_ptr);
    return offset;
}

int send_udp_message(const char *message) {
    int sockfd;
    struct sockaddr_in6 server_addr;
    struct in6_addr server_ip;
    int status;
    
    // 创建UDP socket
    sockfd = socket(AF_INET6, SOCK_DGRAM, 0);
    if (sockfd < 0) {
        perror("socket creation failed");
        return -1;
    }
    
    // 解析主机名获取IPv6地址
    memset(&server_addr, 0, sizeof(server_addr));
    server_addr.sin6_family = AF_INET6;
    server_addr.sin6_port = htons(SERVER_PORT);
    
    struct addrinfo hints, *res;

    memset(&hints, 0, sizeof(hints));
    hints.ai_family = AF_INET6;
    hints.ai_socktype = SOCK_DGRAM;

    status = getaddrinfo(SERVER_HOST, NULL, &hints, &res);
    if (status != 0) {
        fprintf(stderr, "getaddrinfo: %s\n", gai_strerror(status));
        return -1;
    }

    struct sockaddr_in6* ipv6 = (struct sockaddr_in6*)res->ai_addr;
    server_addr.sin6_addr = ipv6->sin6_addr;

    freeaddrinfo(res);
    
    // 发送数据
    status = sendto(sockfd, message, strlen(message), 0,
                   (struct sockaddr *)&server_addr, sizeof(server_addr));
                   
    if (status < 0) {
        perror("sendto failed");
    } else {
        printf("Sent %d bytes to %s:%d\n", status, SERVER_HOST, SERVER_PORT);
    }
    
    close(sockfd);
    return status;
}

int main() {
    char buffer[BUFFER_SIZE];
    int len;
    
    printf("Starting IPv6 address monitor...\n");
    printf("Sending IPv6 addresses to %s:%d every minute\n", SERVER_HOST, SERVER_PORT);
    
    while (1) {
        // 获取本地IPv6地址
        len = get_local_ipv6_addresses(buffer, sizeof(buffer));
        if (len > 0) {
            printf("Collected IPv6 information:\n%s\n", buffer);
            
            // 发送UDP数据包
            if (send_udp_message(buffer) < 0) {
                fprintf(stderr, "Failed to send UDP packet\n");
            }
        } else {
            fprintf(stderr, "Failed to collect IPv6 addresses\n");
        }
        
        // 等待60秒
        sleep(60);
    }
    
    return 0;
}
