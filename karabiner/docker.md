
# command

* `docker version`
* `docker info`
* `docker --help`
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
