

# online platform
play-with-k8s.com

kubectl version
 v1.20.1
kubectl get ns

# istio

curl -L https://istio.io/downloadIstio > is
ISTIO_VERSION=1.5.1 sh is



# 容器编排
Container Orchestration

# CNCF
Cloud Native Computing Foundation，云原生基金会


# kind与minikube

# install minikube


`curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64`

or 

`curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-arm64`

`sudo install minikube /usr/local/bin/`

```bash
minikube version

minikube kubectl
# 会把与当前 Kubernetes 版本匹配的 kubectl 下载下来，存放在内部目录（例如 .minikube/cache/linux/arm64/v1.23.3）

minikube start --kubernetes-version=v1.23.3
minikube status
minikube node list
minikube ssh

kubectl version
kubectl version --short
kubectl version --client

minikube kubectl -- version 
alias kubectl="minikube kubectl --"
source <(kubectl completion bash)

kubectl run ngx --image=nginx:alpine
kubectl get pod
kubectl get node
kubectl get pod -n kube-system
```


# 解决国内网络环境
minikube 提供了特殊的启动参数 --image-mirror-country=cn  --registry-mirror=xxx --image-repository=xxx 等
https://minikube.sigs.k8s.io/docs/handbook/vpn_and_proxy/

# component and addon

* Kubernetes 的节点内部也具有复杂的结构，是由很多的模块构成的，这些模块又可以分成组件（Component）和插件（Addon）两类。 
* 组件实现了 Kubernetes 的核心功能特性，没有这些组件 Kubernetes 就无法启动，
* 而插件则是 Kubernetes 的一些附加功能，属于“锦上添花”，不安装也不会影响 Kubernetes 的正常运行。

# Master 里有 4 个组件，分别是 apiserver、etcd、scheduler、controller-manager。

* apiserver 是 Master 节点——同时也是整个 Kubernetes 系统的**唯一入口**，它对外公开了一系列的 **RESTful API**，并且加上了验证、授权等功能，所有其他组件都只能和它直接通信，可以说是 Kubernetes 里的联络员。 
* etcd 是一个**高可用**的分布式 **Key-Value 数据库**，用来持久化存储系统里的各种资源对象和状态，相当于 Kubernetes 里的配置管理员。注意它只与 apiserver 有直接联系，也就是说任何其他组件想要读写 etcd 里的数据都**必须经过 apiserver**。 
* scheduler 负责容器的编排工作，检查节点的资源状态，把 Pod 调度到最适合的节点上运行，相当于部署人员。因为节点状态和 Pod 信息都存储在 etcd 里，所以 scheduler 必须通过 apiserver 才能获得。 **调度Pod运行**
* controller-manager 负责维护容器和节点等资源的状态，实现故障检测、服务迁移、应用伸缩等功能，相当于**监控运维**人员。同样地，它也必须通过 apiserver 获得存储在 etcd 里的信息，才能够实现对资源的各种操作。

这 4 个组件也都被容器化了，运行在集群的 Pod 里，我们可以用 kubectl 来查看它们的状态，使用命令：
kubectl get pod -n kube-system

# Node 里的 3 个组件了，分别是 kubelet、kube-proxy、container-runtime。
* kubelet 是 **Node 的代理**，负责管理 Node 相关的绝大部分操作，Node 上只有它能够**与 apiserver 通信**，实现状态报告、命令下发、启停容器等功能，相当于是 Node 上的一个“小管家”。 
* kube-proxy 的作用有点特别，它是 Node 的网络代理，只负责管理容器的网络通信，简单来说就是为 Pod **转发 TCP/UDP 数据包**，相当于是专职的“小邮差”。 **实现反向代理**
* 第三个组件 container-runtime 我们就比较熟悉了，它是容器和镜像的实际使用者，在 kubelet 的指挥下创建容器，管理 Pod 的生命周期，是真正干活的“苦力”。**OCI标准**，我们一定要注意，因为 Kubernetes 的定位是容器编排平台，所以它没有限定 container-runtime 必须是 Docker，完全可以替换成任何符合标准的其他容器运行时，例如 containerd、CRI-O 等等，只不过在这里我们使用的是 Docker。

这 3 个组件中只有 kube-proxy 被容器化了，而 kubelet 因为必须要管理整个节点，容器化会限制它的能力，所以它必须在 container-runtime 之外运行。

minikube ssh
docker ps |grep kube-proxy

而 kubelet 用 docker ps 是找不到的，需要用操作系统的 ps 命令：
ps -ef|grep kubelet

# Kubernetes 的大致工作流程
* 每个 Node 上的 kubelet 会定期向 apiserver 上报节点状态，apiserver 再存到 etcd 里。
* 每个 Node 上的 kube-proxy 实现了 TCP/UDP 反向代理，让容器对外提供稳定的服务。
* scheduler 通过 apiserver 得到当前的节点状态，调度 Pod，然后 apiserver 下发命令给某个 Node 的 kubelet，kubelet 调用
* container-runtime 启动容器。 controller-manager 也通过 apiserver 得到实时的节点状态，监控可能的异常情况，再使用相应的手段去调节恢复。

# 插件（Addons）有哪些
使用命令 minikube addons list 就可以查看插件列表

minikube dashboard

# etcd
是由CoreOS公司开发，基于类 Paxos 的 Raft 算法实现数据一致性

# containerd
* dockershim
* cri-dockerd
* 如果 Kubernetes 直接使用 containerd 来操纵容器，那么它就是一个与 Docker 独立的工作环境，彼此都不能访问对方管理的容器和镜像。换句话说，使用命令 docker ps 就看不到在 Kubernetes 里运行的容器了。 这对有的人来说可能需要稍微习惯一下，改用新的工具 crictl，不过用来查看容器、镜像的子命令还是一样的，比如 ps、images 等等，适应起来难度不大（但如果我们一直用 kubectl 来管理 Kubernetes 的话，这就是没有任何影响了）。

# yaml
* YAML Ain't a Markup Language
* 命令式
* 声明式
* 两者不是对立的关系

kubectl api-resources
kubectl get pod --v=9
kubectl get pod
kubectl get job
kubectl get cj
kubectl get pod -o wide
kubectl get pod -o yaml
kubectl get pod -w
kubectl get cm
kubectl get secret


```yaml
apiVersion: v1
kind: Pod
metadata:
  name: ngx-pod
  labels:
    env: demo
    owner: chrono

spec:
  containers:
  - image: nginx:alpine
    name: ngx
    ports:
    - containerPort: 80

```

metadata 里有标签不能任意写，必须要符合域名规范（FQDN）

```yaml
spec:
  containers:
  - image: busybox:latest
    name: busy
    imagePullPolicy: IfNotPresent
    env:
      - name: os
        value: "ubuntu"
      - name: debug
        value: "on"
    command:
      - /bin/echo
    args:
      - "$(os), $(debug)"

```

“header”包含的是 API 对象的基本信息，有三个字段：apiVersion、kind、metadata。 
* apiVersion 表示操作这种资源的 API 版本号，由于 Kubernetes 的迭代速度很快，不同的版本创建的对象会有差异，为了区分这些版本就需要使用 apiVersion 这个字段，比如 v1、v1alpha1、v1beta1 等等。
* kind 表示资源对象的类型，这个应该很好理解，比如 Pod、Node、Job、Service 等等。
* metadata 这个字段顾名思义，表示的是资源的一些“元信息”，也就是用来标记对象，方便 Kubernetes 管理的一些信息。

使用 kubectl apply、kubectl delete，再加上参数 -f，你就可以使用这个 YAML 文件，创建或者删除对象了：
kubectl apply -f ngx-pod.yml
kubectl delete -f ngx-pod.yml

kubectl delete pod busy-pod
kubectl logs busy-pod

kubectl describe pod busy-pod

echo 'aaa' > a.txt
kubectl cp a.txt ngx-pod:/tmp

kubectl exec -it ngx-pod -- sh

kubectl explain pod
kubectl explain pod.metadata
kubectl explain pod.spec
kubectl explain pod.spec.containers

kubectl run ngx --image=nginx:alpine --dry-run=client -o yaml

export out="--dry-run=client -o yaml"
kubectl run ngx --image=nginx:alpine $out

kubectl create job echo-job --image=busybox $out
kubectl create cj echo-cj --image=busybox --schedule="" $out


Cron 来源于希腊语  Chronos 意思是时间

ttlSecondsAfterFinished
job运行结束后保留期限

https://crontab.guru/

successfulJobsHistoryLimit
保留最近N个job执行结果

kubectl create cm info $out
kubectl create cm info --from-literal=k=v $out
kubectl describe cm info

kubectl create secret generic user --from-literal=name=root $out
kubectl describe secret user

kubectl explain pod.spec.containers.env.valueFrom


如果已经存在的配置文件，可以使用 --from-file 从文件自动创建出ConfigMap 或 Secret

https://kubernetes.io/blog/2018/07/18/11-ways-not-to-get-hacked/


kubectl port-forward wp-pod 8080:80 &

```nginx
server {
  listen 80;
  default_type text/html;

  location / {
      proxy_http_version 1.1;
      proxy_set_header Host $host;
      proxy_pass http://127.0.0.1:8080;
  }
}
```

minikube dashboard
kubectl  edit