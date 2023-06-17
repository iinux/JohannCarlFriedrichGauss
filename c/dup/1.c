#include <stdio.h>
#include <unistd.h>
#include <fcntl.h>

int main()
{
    // 打开一个文件，获取它的文件描述符
    // 3
    int fd = open("output.txt", O_WRONLY | O_CREAT | O_TRUNC, 0666);
    if (fd == -1) {
        perror("open");
        return -1;
    }

    // 用dup函数复制标准输出的文件描述符，返回一个新的文件描述符
    // 1->4
    int stdout_copy = dup(STDOUT_FILENO);
    if (stdout_copy == -1) {
        perror("dup");
        return -1;
    }

    // 用close函数关闭标准输出的文件描述符
    // close 1
    if (close(STDOUT_FILENO) == -1) {
        perror("close");
        return -1;
    }

    // 用dup函数复制打开的文件的文件描述符，返回一个新的文件描述符，并赋值给标准输出的文件描述符
    // 3->1
    int new_fd = dup(fd);
    if (new_fd == -1) {
        perror("dup");
        return -1;
    }

    // 用printf函数输出一些内容，这些内容会被重定向到文件中
    printf("Hello, world!\n");
    printf("This is a test of dup function.\n");
    return 0;
}

