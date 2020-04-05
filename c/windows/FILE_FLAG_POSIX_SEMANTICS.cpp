#include <windows.h>
#include <iostream>
#include <string>
#include <atlbase.h>
#include "ConsoleApplication2.h"

/*
[HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\kernel]
"obcaseinsensitive" = dword:00000000
https://stackoverflow.com/questions/8044506/how-to-convert-from-lpcstr-to-lpcwstr-in-c
https://blog.yuwu.me/?p=3942
https://docs.microsoft.com/zh-cn/windows/win32/api/fileapi/nf-fileapi-createfilea
*/

using namespace std;

std::string get_content(HANDLE handle)
{
    const int BUFSIZE = 100;
    char buffer[BUFSIZE];
    DWORD nRead = 0;
    memset(buffer, 0, BUFSIZE);
    ReadFile(handle, buffer, BUFSIZE - 1, &nRead, NULL);
    return buffer;
}

std::string set_content(HANDLE handle, const char* buffer)
{
    DWORD nWrite = 0;
    WriteFile(handle, buffer, strlen(buffer), &nWrite, NULL);
    return buffer;
}

int main(int argc, char* argv[])
{
    if (argc < 3) {
        cout << "not enough argv" << endl;
        return 1;
    }
    USES_CONVERSION;
    char* method = argv[1];
    LPCWSTR file = A2W(argv[2]);
    
    if (strcmp(method, "get") == 0) {
        HANDLE h1 = CreateFile(
            file,
            GENERIC_READ | GENERIC_WRITE,
            FILE_SHARE_READ,
            NULL,
            OPEN_EXISTING,
            FILE_FLAG_POSIX_SEMANTICS,
            NULL);
        cout << get_content(h1) << endl;

        CloseHandle(h1);
    }
    else if (strcmp(method, "set") == 0) {
        HANDLE h1 = CreateFile(
            file,
            GENERIC_READ | GENERIC_WRITE,
            FILE_SHARE_READ,
            NULL,
            CREATE_NEW,
            FILE_FLAG_POSIX_SEMANTICS,
            NULL);
        if (argc < 4) {
            cout << "not enough argv" << endl;
            return 1;
        }
        char* content = argv[3];
        set_content(h1, content);

        CloseHandle(h1);
    }
    else {
        cout << "error method" << endl;
        return 1;
    }
    
    return 0;
}
