#include <unistd.h>
#include <string.h>

#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>

#include <netinet/tcp.h>

#ifndef MSG_FASTOPEN
#define MSG_FASTOPEN 0x20000000
#endif

// gcc -DFAST_OPEN  client.c -o client

int main()
{
    int clientSock;
    struct sockaddr_in addr;

    clientSock = socket(AF_INET, SOCK_STREAM, 0);

    if (clientSock == -1)
    {
        write(STDERR_FILENO, "create socket failed!\n", 8);
        return 1;
    }

#ifdef DEFER_ACCEPT
    int soValue = 1;
    if (setsockopt(clientSock, IPPROTO_TCP, TCP_DEFER_ACCEPT, &soValue, sizeof(soValue)) < 0)
    {
        write(STDERR_FILENO, "failed!\n", 8);
        return 10;
    }
#endif

    memset(&addr, 0, sizeof(addr));
    addr.sin_family = AF_INET;
    addr.sin_port = htons(7890);
    addr.sin_addr.s_addr = inet_addr("1.1.1.1");

#ifdef FAST_OPEN
    char* str = "HelloFastOpen\n";
    int ret = sendto(clientSock, str, strlen(str), MSG_FASTOPEN,
                     (struct sockaddr *)&addr, sizeof(addr));

    if (ret < 0)
    {
        write(STDERR_FILENO, "sendto failed!\n", 8);
        return 11;
    }

#else
    if (connect(clientSock, (struct sockaddr *)&addr, sizeof(addr)) < 0)
    {
        write(STDERR_FILENO, "connect failed!\n", 8);
        return 2;
    }
    char* str = "HelloNoFastOpen\n";
    if (send(clientSock, str, strlen(str), 0) < 0)
    {
        write(STDERR_FILENO, "send failed!\n", 8);
        return 3;
    }
#endif

    close(clientSock);

    return 0;
}
