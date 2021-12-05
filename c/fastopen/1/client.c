#include <sys/socket.h>
#include <netinet/in.h>
#include <netdb.h>
#include <errno.h>
int main()
{
    struct sockaddr_in serv_addr;
    struct hostent *server;

    char *data = "Hello, tcp fast open";
    int data_len = strlen(data);

    int sfd = socket(AF_INET, SOCK_STREAM, 0);
    server = gethostbyname("hw.iinux.cn");

    bzero((char *)&serv_addr, sizeof(serv_addr));
    serv_addr.sin_family = AF_INET;
    bcopy((char *)server->h_addr,
          (char *)&serv_addr.sin_addr.s_addr,
          server->h_length);
    serv_addr.sin_port = htons(6666);

    // /usr/src/linux-headers-4.4.0-22/include/linux/socket.h:#define MSG_FASTOPEN	0x20000000	/* Send data in TCP SYN */

    //  int len = sendto(sfd, data, data_len, 0x20000000/*MSG_FASTOPEN*/,

    int len = sendto(sfd, data, data_len, MSG_FASTOPEN /*MSG_FASTOPEN*/,
                     (struct sockaddr *)&serv_addr, sizeof(serv_addr));
    if (errno != 0)
    {
        printf("error: %s\n", strerror(errno));
    }
    close(sfd);
}