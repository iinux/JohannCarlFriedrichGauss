#include <unistd.h>
#include <string.h>

#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <netinet/tcp.h>

// gcc -DFAST_OPEN  server.c -o server

int main()
{
    int serverSock;
    struct sockaddr_in addr;

    int clientSock;
    struct sockaddr_in clientAddr;
    int addrLen;

    char buf[1024];
    int read;

    serverSock = socket(AF_INET, SOCK_STREAM, 0);

    if (serverSock == -1)
    {
        write(STDERR_FILENO, "failed!\n", 8);
        return 1;
    }

#ifdef DEFER_ACCEPT
    int soValue = 1;
    if (setsockopt(serverSock, IPPROTO_TCP, TCP_DEFER_ACCEPT, &soValue, sizeof(soValue)) < 0)
    {
        write(STDERR_FILENO, "failed!\n", 8);
        return 10;
    }
#endif

    memset(&addr, 0, sizeof(addr));
    addr.sin_family = AF_INET;
    addr.sin_port = htons(7890);
    addr.sin_addr.s_addr = inet_addr("0.0.0.0");

    if (bind(serverSock, (struct sockaddr *)&addr, sizeof(addr)) < 0)
    {
        write(STDERR_FILENO, "failed!\n", 8);
        return 2;
    }

#ifdef FAST_OPEN
    int qlen = 5;
    setsockopt(serverSock, SOL_TCP, TCP_FASTOPEN, &qlen, sizeof(qlen));
#endif

    if (listen(serverSock, 511) < 0)
    {
        write(STDERR_FILENO, "failed!\n", 8);
        return 3;
    }

    while (1)
    {
        addrLen = sizeof(clientAddr);
        clientSock = accept(serverSock, (struct sockaddr *)&clientAddr, &addrLen);
        if (clientSock < 0)
        {
            write(STDERR_FILENO, "failed!\n", 8);
            return 4;
        }

        read = recv(clientSock, buf, 1024, 0);
        write(STDOUT_FILENO, buf, read);
        close(clientSock);
    }

    return 0;
}
