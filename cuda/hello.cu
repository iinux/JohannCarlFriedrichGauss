// hello.cu
#include <iostream>
__global__ void helloCUDA() {
    printf("Hello from CUDA!\n");
}
int main() {
    helloCUDA<<<1, 1>>>();
    cudaDeviceSynchronize();
    return 0;
}

// nvcc hello.cu -o hello
