
https://bowers.github.io/eBPF-Hello-World/

more /boot/config-$(uname -r) | grep CONFIG_BPF
sudo strace -e bpf python3 listen.py

---

https://tonybai.com/2022/07/05/develop-hello-world-ebpf-program-in-c-from-scratch/

git clone https://github.com/libbpf/libbpf-bootstrap.git
git submodule update --init --recursive

libbpf_bootstrap/examples/c/ add

```c
// helloworld.bpf.c 

#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>

SEC("tracepoint/syscalls/sys_enter_execve")

int bpf_prog(void *ctx) {
  char msg[] = "Hello, World!";
  bpf_printk("invoke bpf_prog: %s\n", msg);
  return 0;
}

char LICENSE[] SEC("license") = "Dual BSD/GPL";

// helloworld.c

#include <stdio.h>
#include <unistd.h>
#include <sys/resource.h>
#include <bpf/libbpf.h>
#include "helloworld.skel.h"

static int libbpf_print_fn(enum libbpf_print_level level, const char *format, va_list args)
{
    return vfprintf(stderr, format, args);
}

int main(int argc, char **argv)
{
    struct helloworld_bpf *skel;
    int err;

    libbpf_set_strict_mode(LIBBPF_STRICT_ALL);
    /* Set up libbpf errors and debug info callback */
    libbpf_set_print(libbpf_print_fn);

    /* Open BPF application */
    skel = helloworld_bpf__open();
    if (!skel) {
        fprintf(stderr, "Failed to open BPF skeleton\n");
        return 1;
    }   

    /* Load & verify BPF programs */
    err = helloworld_bpf__load(skel);
    if (err) {
        fprintf(stderr, "Failed to load and verify BPF skeleton\n");
        goto cleanup;
    }

    /* Attach tracepoint handler */
    err = helloworld_bpf__attach(skel);
    if (err) {
        fprintf(stderr, "Failed to attach BPF skeleton\n");
        goto cleanup;
    }

    printf("Successfully started! Please run `sudo cat /sys/kernel/debug/tracing/trace_pipe` "
           "to see output of the BPF programs.\n");

    for (;;) {
        /* trigger our BPF program */
        fprintf(stderr, ".");
        sleep(1);
    }

cleanup:
    helloworld_bpf__destroy(skel);
    return -err;
}
```

libbpf_bootstrap/examples/c/Makefile add
APPS = helloworld minimal minimal_legacy bootstrap uprobe kprobe fentry

---

sudo bpftrace -l | grep open
sudo bpftool prog list
sudo bpftool prog dump xlated id 89
sudo bpftool prog dump jited id 89
sudo bpftool version -p
bpftool feature probe
bpftool btf dump file /sys/kernel/btf/vmlinux format c >> vmlinux.h
bpftool map dump id 386
bpftool map dump name mapname
bpftool perf