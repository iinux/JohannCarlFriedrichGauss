# Example of loading an eBPF program to a Linux kernel tracepoint
# requires bpfcc-tools
# To run: sudo python3 listen.py

from bcc import BPF
from bcc.utils import printb

BPF_SOURCE_CODE = r"""
TRACEPOINT_PROBE(syscalls, sys_enter_mkdir) {
   bpf_trace_printk("Hello. Creating new directory: %s\n", args->pathname);
   return 0;
}
"""

BPF_SOURCE_CODE = r"""
TRACEPOINT_PROBE(syscalls, sys_enter_execve) {
   bpf_trace_printk("Hello. Creating new program: %s\n", args->filename);
   return 0;
}
"""

bpf = BPF(text = BPF_SOURCE_CODE)

print("To test, open another shell window and create a directory, e.g.")
print(" mkdir frodo")
print("CTRL-C to exit")

while True:
    try:
        (task, pid, cpu, flags, ts, msg) = bpf.trace_fields()
        printb(b"%-6d %s" % (pid, msg))
    except ValueError:
        continue
    except KeyboardInterrupt:
        break

