
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

