; https_baidu.asm
; x86_64 Linux NASM program to fetch https://www.baidu.com via OpenSSL
;
; Build:
;   nasm -f elf64 https_baidu.asm -o https_baidu.o
;   gcc -no-pie https_baidu.o -o https_baidu -lssl -lcrypto

default rel

global main

extern socket, connect, gethostbyname, htons, close
extern write, perror, exit, printf
extern SSL_CTX_new, SSL_new, SSL_set_fd, SSL_connect
extern SSL_write, SSL_read, SSL_free, SSL_CTX_free
extern SSL_shutdown
extern TLS_client_method
extern ERR_print_errors_fp, stderr

section .data
    hostname:          db "www.baidu.com", 0

    get_request:
        db "GET / HTTP/1.1", 13, 10
        db "Host: www.baidu.com", 13, 10
        db "User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36", 13, 10
        db "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8", 13, 10
        db "Accept-Language: en-US,en;q=0.5", 13, 10
        db "Connection: close", 13, 10
        db 13, 10
    request_len: dq $ - get_request

    msg_connecting:    db "Connecting to www.baidu.com:443 ...", 10, 0
    msg_connected:     db "Connected, performing TLS handshake...", 10, 0
    msg_tls_ok:        db "TLS handshake OK, sending request...", 10, 0
    msg_response:      db 10, "=== Response from www.baidu.com ===", 10, 0
    msg_done:          db 10, "Done.", 10, 0

    err_socket:        db "socket() failed", 0
    err_connect_str:   db "connect() failed", 0
    err_dns:           db "gethostbyname() failed", 0
    err_ssl_ctx:       db "SSL_CTX_new() failed", 0
    err_ssl_conn:      db "SSL_connect() failed", 0
    err_ssl_write:     db "SSL_write() failed", 0
    err_ssl_read:      db "SSL_read() failed", 0

    ; offsets in struct hostent (x86_64):
    ; +0: h_name     (char *)       8 bytes
    ; +8: h_aliases  (char **)      8 bytes
    ; +16: h_addrtype (int)         4 bytes
    ; +20: h_length   (int)         4 bytes
    ; +24: h_addr_list (char **)    8 bytes
    HOSTENT_ADDRLIST equ 24

section .bss
    sockfd:     resq 1
    ssl_ctx:    resq 1
    ssl:        resq 1
    server_addr: resb 16    ; struct sockaddr_in
    hostent_ptr: resq 1
    buffer:     resb 8192

section .text
main:
    push rbp
    mov rbp, rsp

    ; ---- print status ----
    mov rdi, msg_connecting
    xor eax, eax
    call printf

    ; ---- 1. socket(AF_INET, SOCK_STREAM, 0) ----
    mov rdi, 2          ; AF_INET
    mov rsi, 1          ; SOCK_STREAM
    xor rdx, rdx        ; protocol = 0
    call socket
    cmp rax, -1
    je  .socket_fail
    mov [sockfd], rax

    ; ---- 2. gethostbyname("www.baidu.com") ----
    mov rdi, hostname
    call gethostbyname
    test rax, rax
    jz   .dns_fail
    mov [hostent_ptr], rax

    ; ---- 3. prepare sockaddr_in ----
    ; sin_family = AF_INET = 2
    mov word [server_addr], 2

    ; sin_port = htons(443)
    mov rdi, 443
    call htons
    mov word [server_addr + 2], ax

    ; sin_addr = *(hostent->h_addr_list[0])
    mov rax, [hostent_ptr]
    mov rax, [rax + HOSTENT_ADDRLIST]   ; rax = h_addr_list (char **)
    mov rax, [rax]                       ; rax = h_addr_list[0] (char *)
    mov ecx, [rax]                       ; ecx = 4-byte IP address
    mov [server_addr + 4], ecx

    ; sin_zero is already zeroed from .bss

    ; ---- 4. connect(sockfd, &server_addr, 16) ----
    mov rdi, [sockfd]
    lea rsi, [server_addr]
    mov rdx, 16
    call connect
    test rax, rax
    jnz  .connect_fail

    mov rdi, msg_connected
    xor eax, eax
    call printf

    ; ---- 5. TLS: get method and create SSL_CTX ----
    call TLS_client_method
    test rax, rax
    jz   .ssl_ctx_fail
    mov rdi, rax
    call SSL_CTX_new
    test rax, rax
    jz   .ssl_ctx_fail
    mov [ssl_ctx], rax

    ; ---- 6. SSL_new(ssl_ctx) ----
    mov rdi, [ssl_ctx]
    call SSL_new
    test rax, rax
    jz   .ssl_fail
    mov [ssl], rax

    ; ---- 7. SSL_set_fd(ssl, sockfd) ----
    mov rdi, [ssl]
    mov rsi, [sockfd]
    call SSL_set_fd

    ; ---- 8. SSL_connect(ssl) ----
    mov rdi, [ssl]
    call SSL_connect
    cmp rax, 1
    jne .ssl_conn_fail

    mov rdi, msg_tls_ok
    xor eax, eax
    call printf

    ; ---- 9. SSL_write(ssl, request, request_len) ----
    mov rdi, [ssl]
    mov rsi, get_request
    mov rdx, [request_len]
    call SSL_write
    test rax, rax
    jle .ssl_write_fail

    ; ---- 10. read & print response ----
    mov rdi, msg_response
    xor eax, eax
    call printf

.read_loop:
    mov rdi, [ssl]
    lea rsi, [buffer]
    mov rdx, 8192
    call SSL_read

    cmp rax, 0
    jle .read_done          ; 0 = EOF, negative = error

    ; write(1, buffer, rax)
    mov rdi, 1              ; stdout
    lea rsi, [buffer]
    mov rdx, rax
    call write
    jmp .read_loop

.read_done:
    js  .ssl_read_fail      ; negative means error

    ; ---- 11. cleanup ----
    mov rdi, [ssl]
    call SSL_shutdown
    mov rdi, [ssl]
    call SSL_free
    mov rdi, [ssl_ctx]
    call SSL_CTX_free
    mov rdi, [sockfd]
    call close

    mov rdi, msg_done
    xor eax, eax
    call printf

    xor eax, eax
    pop rbp
    ret

; ---------- error handlers ----------
.socket_fail:
    mov rdi, err_socket
    call perror
    mov rdi, 1
    call exit

.dns_fail:
    mov rdi, err_dns
    call perror
    mov rdi, 1
    call exit

.connect_fail:
    mov rdi, err_connect_str
    call perror
    mov rdi, 1
    call exit

.ssl_ctx_fail:
    mov rdi, err_ssl_ctx
    call perror
    mov rdi, 1
    call exit

.ssl_fail:
    mov rdi, err_ssl_conn
    call perror
    mov rdi, 1
    call exit

.ssl_conn_fail:
    mov rdi, [stderr]       ; FILE *
    mov rsi, 0              ; flags
    mov rdx, 0              ; unused
    call ERR_print_errors_fp
    jmp .cleanup_and_exit

.ssl_write_fail:
    mov rdi, err_ssl_write
    call perror
    jmp .cleanup_and_exit

.ssl_read_fail:
    mov rdi, err_ssl_read
    call perror
    jmp .cleanup_and_exit

.cleanup_and_exit:
    mov rdi, [ssl]
    test rdi, rdi
    jz   .skip_ssl_free
    call SSL_free
.skip_ssl_free:
    mov rdi, [ssl_ctx]
    test rdi, rdi
    jz   .skip_ctx_free
    call SSL_CTX_free
.skip_ctx_free:
    mov rdi, [sockfd]
    call close
    mov rdi, 1
    call exit
