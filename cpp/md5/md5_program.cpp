#include <iostream>
#include <openssl/md5.h>

std::string calculateMD5(const std::string& input) {
    unsigned char digest[MD5_DIGEST_LENGTH];
    MD5(reinterpret_cast<const unsigned char*>(input.c_str()), input.length(), digest);

    char md5String[2 * MD5_DIGEST_LENGTH + 1];
    for (int i = 0; i < MD5_DIGEST_LENGTH; ++i) {
        sprintf(&md5String[i * 2], "%02x", static_cast<unsigned int>(digest[i]));
    }

    return md5String;
}

int main() {
    std::string input;
    std::cout << "Enter the string to calculate MD5: ";
    std::cin >> input;

    std::string md5 = calculateMD5(input);
    std::cout << "MD5: " << md5 << std::endl;

    return 0;
}

