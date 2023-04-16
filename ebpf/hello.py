from bcc import BPF

b = BPF(src_file="hello.c")
b.attach_kprobe(event="do_sys_open", fn_name="hello_world")
b.trace_print()
