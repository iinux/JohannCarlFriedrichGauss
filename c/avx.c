// 一个简单的数组相加函数，编译器可能将其自动向量化
void add_arrays(float* a, float* b, float* c, int n) {
    for (int i = 0; i < n; i++) {
        c[i] = a[i] + b[i];
    }
}

int main() {
    return 0;
}
