void sayHello() {
    const char* s = "hello\n";
    __asm__("int $0x80\n\r"
    ::"a"(4), "b"(1), "c"(s), "d"(6):);
}
int main() {
    sayHello();
    return 0;
}
