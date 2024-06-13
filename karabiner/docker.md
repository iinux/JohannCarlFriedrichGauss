
# basic
 * 容器技术中有三个核心概念：容器（Container）、镜像（Image），以及镜像仓库（Registry）
 * like docker: kata gVisor rkt podman
 * 原理：Linux 操作系统内核之中，为资源隔离提供了三种技术：namespace、cgroup、chroot
 * Alpine Linux
 * chroot->pivot_root
 * LXC(Linux Containers)->libcontainer->containerd+runc
 * UnionFS implement: aufs, btrfs, device-mapper, overlay2(current docker used)

# command

* `docker --help`
* version
* info
* images
    * `docker images -a`
    * `docker images -q`
    * `docker images --digests`
    * `docker images --digests --no-trunc`
* search
    * `docker search tomcat --no-trunc`
* run
    * -i
    * -t
    * --name
    * -h
    * -d
    * -P
    * -p
        * ip:hostPort:containerPort
        * ip::containerPort
        * hostPort:containerPort
        * containerPort
    * -v
        * default is rw
        * -v /hosttmp:/dockertmp:ro
    * --privileged=true
    * --volumes-from
    * --rm
    * --net=
        * null
        * host
        * bridge
    * --mount
    * --env
* ps
    * -a
    * -l
    * -n \<number>
    * -q
    * --no-trunc
* exit
    * exit
    * ctrl+p+q
* start
* restart
* stop
* kill
* rm
* rmi
* logs
    * -f
    * -t
    * --tail \<number>
* top
* inspect
* attach
* exec
* cp
* commit
* build
    * -f 
    * -t
* history
* login
* tag
* push
* pull
* save
* load

# dockerfile

* FROM
    * scratch
* ENV
* RUN
* ARG
* VOLUME
* COPY
* ENTRYPOINT
* EXPOSE
* CMD
* MAINTAINER
* ADD
* LABEL
* WORKDIR
* ONBUILD

# .dockerignore

# Registry
 * Docker Registry
 * CNCF Harbor

# /etc/docker/daemon.json

# public hub
 * hub.docker.com
 * quay.io
 * gcr.io
 * ghcr.io

# self hub
 * docker pull registry
 * docker run -d -p 5000:5000 registry
 * docker tag nginx:alpine 127.0.0.1:5000/nginx:alpine
 * docker push 127.0.0.1:5000/nginx:alpine
 * docker rmi  127.0.0.1:5000/nginx:alpine
 * docker pull 127.0.0.1:5000/nginx:alpine
 * curl 127.1:5000/v2/_catalog
 * curl 127.1:5000/v2/nginx/tags/list
 * /var/lib/registry

# remote docker

server

/etc/sysconfig/docker

OPTIONS add -H tcp://0.0.0.0:2375

client

export DOCKER_HOST=tcp://127.0.0.1:2375

# namespace

ip netns help
ip netns ls
ip netns add net1
ls /var/run/netns

ip netns exec net1 ip addr
ip netns exec net1 bash
ip netns exec net1 bash --rcfile <(echo "PS1=\"namespace ns1> \"")

ip netns exec net1 ip link set lo up

ip link add type veth
ip link add vethfoo type veth peer name vethbar
ip link

ip link set veth0 netns net0
ip link set veth1 netns net1

ip netns exec net0 ip link set veth0 up
ip netns exec net0 ip addr add 10.0.1.1/24 dev veth0

ip netns exec net0 ip route

ip netns exec net1 ip link set veth1 up
ip netns exec net1 ip addr add 10.0.1.2/24 dev veth1

---

ip link add br0 type bridge
ip link set dev br0 up

ip netns exec net0 ip link set dev veth1 name eth0
ip link set dev veth0 master br0

bridge link

lsns

Linux 中有三个系统调用可以直接操作命名空间，分别是 clone()、unshare() 和 setns()。

unshare -m
mkdir mount-dir
mount -n -o size=10m -t tmpfs tmpfs mount-dir
df mount-dir
touch mount-dir/{0,1,2}

cat /proc/{pid}/mountinfo
mount -l
df -h

---

unshare -fp --mount-proc

---

unshare -n

ip link show type veth

---

unshare -u

---

id -u

cat /proc/self/uid_map
cat /proc/self/gid_map

unshare -U

cat /proc/sys/kernel/overflow{g,u}id

cat /proc/self/setgroups

# chroot

yum install skopeo
skopeo copy docker://opensuse/tumbleweed:latest oci:tumbleweed:latest
umoci unpack --image tumbleweed:latest bundle
chroot bundle/rootfs


```c
#include <sys/stat.h>
#include <unistd.h>

int main(void) {
  mkdir(".out", 0755);
  chroot(".out");
  chdir("../../../../../");
  chroot(".");
  return execl("/bin/bash", "-i", NULL);
}
```

# cgroup

mkdir /sys/fs/cgroup/memory/demo
ls /sys/fs/cgroup/memory/demo

echo 100000000 > /sys/fs/cgroup/memory/demo/memory.limit_in_bytes
echo 0 > /sys/fs/cgroup/memory/demo/memory.swappiness
echo $$ > /sys/fs/cgroup/memory/demo/cgroup.procs


> # 设置只使用 2,3 两个核
> mkdir -p /sys/fs/cgroup/cpuset/demo
> echo "2,3" > /sys/fs/cgroup/cpuset/demo/cpuset.cpus
> echo "0" > /sys/fs/cgroup/cpuset/demo/cpuset.mems
> echo $$ > /sys/fs/cgroup/cpuset/demo/cgroup.procs

> # 设置CPU利用率限制为50%
> mkdir -p /sys/fs/cgroup/cpu/demo
> echo "100000" > /sys/fs/cgroup/cpu/demo/cpu.cfs_quota_us
> echo "100000" > /sys/fs/cgroup/cpu/demo/cpu.cfs_period_us
> echo $$ > /sys/fs/cgroup/cpu/demo/cgroup.procs

> yes > /dev/null&
> yes > /dev/null&

nsenter --pid=/proc/$PID/ns/pid unshare --mount-proc

# chroot namespace cgroup

```c
#define _GNU_SOURCE
#include <errno.h>
#include <sched.h>
#include <stdio.h>
#include <string.h>
#include <sys/mount.h>
#include <sys/msg.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>

#define STACKSIZE (1024 * 1024)

static char stack[STACKSIZE];
void print_err(char const *const reason) {
  fprintf(stderr, "Error %s: %s\n", reason, strerror(errno));
}
int exec(void *args) {
  // Remount proc
  if (mount("proc", "/proc", "proc", 0, "") != 0) {
    print_err("mounting proc");
    return 1;
  }
  // Set a new hostname
  char const *const hostname = "new-hostname";
  if (sethostname(hostname, strlen(hostname)) != 0) {
    print_err("setting hostname");
    return 1;
  }
  // Create a message queue
  key_t key = {0};
  if (msgget(key, IPC_CREAT) == -1) {
    print_err("creating message queue");
    return 1;
  }
  // Execute the given command
  char **const argv = args;
  if (execvp(argv[0], argv) != 0) {
    print_err("executing command");
    return 1;
  }
  return 0;
}
int main(int argc, char **argv) {
  // Provide some feedback about the usage
  if (argc < 2) {
    fprintf(stderr, "No command specified\n");
    return 1;
  }
  // Namespace flags
  const int flags = CLONE_NEWNET | CLONE_NEWUTS | CLONE_NEWNS | CLONE_NEWIPC |
                    CLONE_NEWPID | CLONE_NEWUSER | SIGCHLD;
  // Create a new child process
  pid_t pid = clone(exec, stack + STACKSIZE, flags, &argv[1]);
  if (pid < 0) {
    print_err("calling clone");
    return 1;
  }
  // Wait for the process to finish
  int status = 0;
  if (waitpid(pid, &status, 0) == -1) {
    print_err("waiting for pid");
    return 1;
  }
  // Return the exit code
  return WEXITSTATUS(status);
}
```

> gcc -o namespaces namespaces.c
> ./namespaces ip a
> ./namespaces ps aux
> ./namespaces whoami


runc run -b bundle container

https://panzhongxian.cn/cn/2023/10/demystifying-containers-part-i-kernel-space/