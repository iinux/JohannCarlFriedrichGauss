#include <linux/init.h>
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/version.h>
#include <linux/err.h>
#include <linux/time.h>
#include <linux/skbuff.h>
#include <linux/sched.h>
#include <net/tcp.h>
#include <net/inet_common.h>
#include <linux/uaccess.h>
#include <linux/netdevice.h>
#include <net/net_namespace.h>
#include <linux/mm.h>
#include <linux/kallsyms.h>
#include <net/ipv6.h>
#include <net/transp_v6.h>

// http://just4coding.com/2021/07/21/socket-to-pid/

unsigned long sk_data_ready_addr = 0;

static int
inet_stream_connect_tcpconn(struct socket *sock, struct sockaddr *uaddr,
        int addr_len, int flags)
{
    int retval = 0;
    char *pathname, *p;
    struct mm_struct *mm;
    struct sockaddr_in *in;
    struct inet_sock *inet;

    wait_queue_head_t *q;
    struct task_struct  *t;

    retval = inet_stream_connect(sock, uaddr, addr_len, flags);

    printk(KERN_INFO "inet_stream_connect_tcpconn called\n");

    if (sock && sock->sk) {
        inet = inet_sk(sock->sk);

        in = (struct sockaddr_in *)uaddr;
        if (in && inet) {
            printk(KERN_INFO "CPU [%u] CONN: %08x:%d->%08x:%d\n",
                    smp_processor_id(),
                    ntohl(inet->inet_saddr),
                    ntohs(inet->inet_sport),
                    ntohl(in->sin_addr.s_addr),
                    ntohs(in->sin_port));
        }
    }

    mm = get_task_mm(current);
    if (!mm) {
        goto out;
    }

    down_read(&mm->mmap_sem);
    if (mm->exe_file) {
        pathname = kmalloc(PATH_MAX, GFP_ATOMIC);
        if (pathname) {
            p = d_path(&mm->exe_file->f_path, pathname, PATH_MAX);

            printk(KERN_INFO "CPU [%u], FILE: %s, COMM: %s, PID: %d\n",
                    smp_processor_id(),
                    p, current->comm, current->pid);

            kfree(pathname);
        }
    }

    up_read(&mm->mmap_sem);

    mmput(mm);

out:
    return retval;
}

static inline int
hook_tcpconn_functions(void)
{
    unsigned int level;
    pte_t *pte;

    struct proto_ops *inet_stream_ops_p =
            (struct proto_ops *)&inet_stream_ops;

    pte = lookup_address((unsigned long)inet_stream_ops_p, &level);
    if (pte == NULL) {
        return 1;
    }

    if (pte->pte & ~_PAGE_RW) {
        pte->pte |= _PAGE_RW;
    }

    inet_stream_ops_p->connect = inet_stream_connect_tcpconn;
    printk(KERN_INFO "CPU [%u] hooked inet_stream_connect <%p> --> <%p>\n",
        smp_processor_id(), inet_stream_connect, inet_stream_ops_p->connect);

    return 0;
}

static int
unhook_tcpconn_functions(void)
{
    unsigned int level;
    pte_t *pte;

    struct proto_ops *inet_stream_ops_p =
            (struct proto_ops *)&inet_stream_ops;

    inet_stream_ops_p->connect = inet_stream_connect;
    printk(KERN_INFO "CPU [%u] unhooked inet_stream_connect\n",
        smp_processor_id());

    pte = lookup_address((unsigned long)inet_stream_ops_p, &level);
    if (pte == NULL) {
        return 1;
    }

    pte->pte |= pte->pte & ~_PAGE_RW;

    return 0;
}

static int __init
tcpconn_init(void)
{
    printk(KERN_INFO "loading tcpconn\n");

    sk_data_ready_addr = kallsyms_lookup_name("sock_def_readable");
    printk(KERN_INFO "CPU [%u] sk_data_ready_addr = "
        "kallsyms_lookup_name(sock_def_readable) = %lu\n",
         smp_processor_id(), sk_data_ready_addr);
    if (0 == sk_data_ready_addr) {
        printk(KERN_INFO "cannot find sock_def_readable.\n");
        goto err;
    }

    hook_tcpconn_functions();

    printk(KERN_INFO "tcpconn loaded\n");
    return 0;

err:
    return 1;
}

static void __exit
tcpconn_exit(void)
{
    unhook_tcpconn_functions();
    synchronize_net();

    printk(KERN_INFO "tcpconn unloaded\n");
}

module_init(tcpconn_init);
module_exit(tcpconn_exit);
MODULE_LICENSE("GPL");
