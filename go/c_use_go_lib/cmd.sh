go build -v -x -buildmode=c-shared -o lib.so a.go 
# 生成静态库
# go build -v -x -buildmode=c-archive -o lib.a a.go
sudo mv lib.so /lib64/libgo.so
gcc b.c -o b -lgo
./b
