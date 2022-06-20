
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
    * -d
    * -P
    * -p
        * ip:hostPort:containerPort
        * ip::containerPort
        * hostPort:containerPort
        * containerPort
    * -v
    * --privileged=true
    * --volumes-from
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

# remote docker

server

/etc/sysconfig/docker

OPTIONS add -H tcp://0.0.0.0:2375

client

export DOCKER_HOST=tcp://127.0.0.1:2375

