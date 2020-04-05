#include <stdio.h>
#include <time.h>
#include <unistd.h>

#define BST (+1)
#define CCT (+8)

int main()
{
    clock_t start_t, finish_t;
    double total_t = 0;
    int i = 0;
    start_t = clock();

    time_t seconds;
    seconds = time(NULL);
    printf("timestamp = %ld\n", seconds);

    struct tm t; //更多情况下是通过localtime函数及gmtime函数获得tm结构
    t.tm_sec = 10;
    t.tm_min = 10;
    t.tm_hour = 6;
    t.tm_mday = 25;
    t.tm_mon = 2;
    t.tm_year = 89;
    t.tm_wday = 6;
    printf("%s", asctime(&t));

    time_t timer;
    struct tm *Now;
    time(&timer);
    Now = localtime(&timer);
    printf("当前的本地时间和日期：%s", asctime(Now));

    time_t curtime;
    time(&curtime);
    printf("当前时间 = %s", ctime(&curtime));

    time_t first, second;
    time(&first);
    sleep(2);
    time(&second);
    printf("The difference is: %f seconds\n", difftime(second, first));

    time_t rawtime;
    struct tm *info;
    time(&rawtime);
    // 获取 GMT 时间
    info = gmtime(&rawtime);
    printf("当前的世界时钟：\n");
    printf("伦敦：%2d:%02d\n", (info->tm_hour + BST) % 24, info->tm_min);
    printf("中国：%2d:%02d\n", (info->tm_hour + CCT) % 24, info->tm_min);

    int ret;
    struct tm info2;
    char buffer1[80];
    info2.tm_year = 2001 - 1900;
    info2.tm_mon = 7 - 1;
    info2.tm_mday = 4;
    info2.tm_hour = 0;
    info2.tm_min = 0;
    info2.tm_sec = 1;
    info2.tm_isdst = -1;
    ret = mktime(&info2);
    if (ret == -1)
    {
        printf("错误：不能使用 mktime 转换时间。\n");
    }
    else
    {
        strftime(buffer1, sizeof(buffer1), "%c", &info2);
        printf("%s\n", buffer1);
    }

    time_t rawtime1;
    struct tm *info1;
    char buffer2[80];
    time(&rawtime1);
    info1 = localtime(&rawtime1);
    strftime(buffer2, 80, "%Y%m%e_%H%M%S", info1); //以年月日_时分秒的形式表示当前时间
    printf("%s\n", buffer2);

    finish_t = clock();
    total_t = (double)(finish_t - start_t) / CLOCKS_PER_SEC; //将时间转换为秒
    printf("CPU 占用的总时间：%f\n", total_t);

    return 0;
}