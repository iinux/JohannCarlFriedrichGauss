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

#define SERVER_HOST_V6 "va.google.party"
#define SERVER_HOST_V4 "va4.google.party"
#define SERVER_PORT_V6 5800
#define SERVER_PORT_V4 5801
#define BUFFER_SIZE 1024

int get_local_ipv6_addresses(char *buffer, size_t buffer_size) {
    struct ifaddrs *ifaddrs_ptr, *ifa;
    char ip_str[INET6_ADDRSTRLEN];
    int offset = 0;

    if (getifaddrs(&ifaddrs_ptr) == -1) {
        perror("getifaddrs");
        return -1;
    }

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

int get_local_ipv4_addresses(char *buffer, size_t buffer_size) {
    struct ifaddrs *ifaddrs_ptr, *ifa;
    char ip_str[INET_ADDRSTRLEN];
    int offset = 0;

    if (getifaddrs(&ifaddrs_ptr) == -1) {
        perror("getifaddrs");
        return -1;
    }

    for (ifa = ifaddrs_ptr; ifa != NULL; ifa = ifa->ifa_next) {
        if (ifa->ifa_addr == NULL)
            continue;

        // 只处理IPv4地址且不是回环地址
        if (ifa->ifa_addr->sa_family == AF_INET) {
            struct sockaddr_in *sin = (struct sockaddr_in *)ifa->ifa_addr;

            // 跳过回环地址 (127.0.0.0/8)
            if (ntohl(sin->sin_addr.s_addr) >> 24 == 127)
                continue;

            inet_ntop(AF_INET, &(sin->sin_addr), ip_str, INET_ADDRSTRLEN);
            
            // 只保留192.168.1.x开头的IPv4地址
            if (strncmp(ip_str, "192.168.1.", 10) == 0) {
                offset += snprintf(buffer + offset, buffer_size - offset, "%s", ip_str);
                break;
            }
        }
    }

    freeifaddrs(ifaddrs_ptr);
    return offset;
}

int send_udp_message_v6(const char *message, const char *host, int port) {
    int sockfd;
    struct sockaddr_in6 server_addr;
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
    server_addr.sin6_port = htons(port);

    struct addrinfo hints, *res;

    memset(&hints, 0, sizeof(hints));
    hints.ai_family = AF_INET6;
    hints.ai_socktype = SOCK_DGRAM;

    status = getaddrinfo(host, NULL, &hints, &res);
    if (status != 0) {
        fprintf(stderr, "getaddrinfo for %s: %s\n", host, gai_strerror(status));
        close(sockfd);
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
        printf("Sent %d bytes to %s:%d\n", status, host, port);
    }

    close(sockfd);
    return status;
}

int send_udp_message_v4(const char *message, const char *host, int port) {
    int sockfd;
    struct sockaddr_in server_addr;
    int status;

    // 创建UDP socket
    sockfd = socket(AF_INET, SOCK_DGRAM, 0);
    if (sockfd < 0) {
        perror("socket creation failed");
        return -1;
    }

    memset(&server_addr, 0, sizeof(server_addr));
    server_addr.sin_family = AF_INET;
    server_addr.sin_port = htons(port);

    struct addrinfo hints, *res;

    memset(&hints, 0, sizeof(hints));
    hints.ai_family = AF_INET;
    hints.ai_socktype = SOCK_DGRAM;

    status = getaddrinfo(host, NULL, &hints, &res);
    if (status != 0) {
        fprintf(stderr, "getaddrinfo for %s: %s\n", host, gai_strerror(status));
        close(sockfd);
        return -1;
    }

    struct sockaddr_in* ipv4 = (struct sockaddr_in*)res->ai_addr;
    server_addr.sin_addr = ipv4->sin_addr;

    freeaddrinfo(res);

    // 发送数据
    status = sendto(sockfd, message, strlen(message), 0,
                   (struct sockaddr *)&server_addr, sizeof(server_addr));

    if (status < 0) {
        perror("sendto failed");
    } else {
        printf("Sent %d bytes to %s:%d\n", status, host, port);
    }

    close(sockfd);
    return status;
}

int main() {
    char buffer_v6[BUFFER_SIZE];
    char buffer_v4[BUFFER_SIZE];
    int len_v6, len_v4;

    printf("Starting IP address monitor...\n");
    printf("Sending IPv6 addresses to %s:%d and IPv4 addresses to %s:%d every minute\n", 
           SERVER_HOST_V6, SERVER_PORT_V6, SERVER_HOST_V4, SERVER_PORT_V4);

    while (1) {
        int success = 0;
        
        // 获取本地IPv6地址
        len_v6 = get_local_ipv6_addresses(buffer_v6, sizeof(buffer_v6));
        if (len_v6 > 0) {
            printf("Collected IPv6 address: %s\n", buffer_v6);

            if (send_udp_message_v6(buffer_v6, SERVER_HOST_V6, SERVER_PORT_V6) < 0) {
                fprintf(stderr, "Failed to send IPv6 UDP packet\n");
            } else {
                success = 1;
            }
        } else {
            fprintf(stderr, "No matching IPv6 addresses found\n");
        }

        // 获取本地IPv4地址
        len_v4 = get_local_ipv4_addresses(buffer_v4, sizeof(buffer_v4));
        if (len_v4 > 0) {
            printf("Collected IPv4 address: %s\n", buffer_v4);

            if (send_udp_message_v4(buffer_v4, SERVER_HOST_V4, SERVER_PORT_V4) < 0) {
                fprintf(stderr, "Failed to send IPv4 UDP packet\n");
            } else {
                success = 1;
            }
        } else {
            fprintf(stderr, "No matching IPv4 addresses found\n");
        }

        if (!success) {
            fprintf(stderr, "Failed to collect any addresses\n");
        }

        // 等待60秒
        sleep(60);
    }

    return 0;
}
