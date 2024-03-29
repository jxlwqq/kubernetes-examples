# 使用 kubeadm 创建集群

本地开发环境以 macOS 为例。

### 安装 VirtualBox 和 Vagrant

* 安装 [VirtualBox](https://www.virtualbox.org/)
* 安装 [Vagrant](https://www.vagrantup.com)

### 启动虚拟机

虚拟机配置了 3 台 CKA 考试平台对应的 Ubuntu 18.04 LTS，分别是 `k8s-1(192.168.205.10)`、`k8s-2(192.168.205.11)` 和 `k8s-3(192.168.205.12)`。配置详见 [Vagrantfile](./Vagrantfile)
这个文件。

虚拟机初始化的时候，已经帮助你安装了 Docker 环境，详见 [config.vm.provision "shell"](./Vagrantfile#L39) 中信息。Vagrant 是用 Ruby
写的，语法都是通用的，应该能看懂，看不懂也没关系。

```shell
git clone https://github.com/jxlwqq/kubernetes-examples.git # clone 仓库到本地
cd cka-training # 进入这个目录
vagrant box add ubuntu/bionic64 # 提前下载考试对应的操作系统镜像文件(ubuntu 18.04 LTS)，方便后续快速启动
vagrant up # 启动虚拟机
# vagrant halt # 关闭虚拟机
# vagrant destroy #销毁虚拟机 
```

### 登录虚拟机

开3个命令行窗口，分别登录这3台虚拟机：

```shell
cd cka-training # 一定要进入在 Vagrantfile 所在的目录
vagrant ssh k8s-1 # 这台作为 master
vagrant ssh k8s-2 # node
vagrant ssh k8s-3 # node
```

### 允许 iptables 检查桥接流量

> 提示：在虚拟机：`k8s-1`、`k8s-2` 和 `k8s-3`中执行命令。

确保 br_netfilter 模块被加载。这一操作可以通过运行 lsmod | grep br_netfilter 来完成。若要显式加载该模块，可执行 sudo modprobe br_netfilter。

为了让你的 Linux 节点上的 iptables 能够正确地查看桥接流量，你需要确保在你的 sysctl 配置中将 net.bridge.bridge-nf-call-iptables 设置为 1。如下：

```shell
cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf
br_netfilter
EOF

cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF
sudo sysctl --system
```

### 安装 kubeadm、kubelet 和 kubectl

> 提示：在虚拟机：`k8s-1`、`k8s-2` 和 `k8s-3`中执行命令。

使用阿里云的镜像进行安装：

```shell

sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates curl

sudo curl -s https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | sudo apt-key add -
echo "deb https://mirrors.aliyun.com/kubernetes/apt/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list

# 关闭 swap
swapoff -a

sudo apt-get update
sudo apt-get install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl

# 自动补全
echo "source <(kubectl completion bash)" >> ~/.bashrc
```

### 配置 cgroup 驱动程序

> 提示：在虚拟机：`k8s-1`、`k8s-2` 和 `k8s-3`中执行命令。

配置 Docker 守护程序，尤其是使用 systemd 来管理容器的 cgroup：

```shell
sudo mkdir /etc/docker
cat <<EOF | sudo tee /etc/docker/daemon.json
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "storage-driver": "overlay2"
}
EOF

sudo systemctl daemon-reload
sudo systemctl restart docker
```

### 验证与 gcr.io 容器镜像仓库的连通性

> 提示：在虚拟机：`k8s-1`、`k8s-2` 和 `k8s-3`中执行命令。

使用 `kubeadm config images pull` 验证与 gcr.io 容器镜像仓库的连通性，不过会失败。

我们通过拉取阿里云的Docker镜像仓库，然后在打 Tag 的方式，来曲线解决这个问题。

node 节点上部分 image 是使用不到的，我们暂时就都一股脑的都 pull 下来。

```shell
kubeadm config images list # 查看所需的镜像

# 替换 kube-apiserver
sudo docker pull registry.aliyuncs.com/google_containers/kube-apiserver:v1.22.0
sudo docker tag registry.aliyuncs.com/google_containers/kube-apiserver:v1.22.0 k8s.gcr.io/kube-apiserver:v1.22.0
sudo docker rmi registry.aliyuncs.com/google_containers/kube-apiserver:v1.22.0

# 替换 kube-controller-manager
sudo docker pull registry.aliyuncs.com/google_containers/kube-controller-manager:v1.22.0
sudo docker tag registry.aliyuncs.com/google_containers/kube-controller-manager:v1.22.0 k8s.gcr.io/kube-controller-manager:v1.22.0
sudo docker rmi registry.aliyuncs.com/google_containers/kube-controller-manager:v1.22.0

# 替换 kube-scheduler
sudo docker pull registry.aliyuncs.com/google_containers/kube-scheduler:v1.22.0
sudo docker tag registry.aliyuncs.com/google_containers/kube-scheduler:v1.22.0 k8s.gcr.io/kube-scheduler:v1.22.0
sudo docker rmi registry.aliyuncs.com/google_containers/kube-scheduler:v1.22.0

# 替换 kube-proxy
sudo docker pull registry.aliyuncs.com/google_containers/kube-proxy:v1.22.0
sudo docker tag registry.aliyuncs.com/google_containers/kube-proxy:v1.22.0 k8s.gcr.io/kube-proxy:v1.22.0
sudo docker rmi registry.aliyuncs.com/google_containers/kube-proxy:v1.22.0

# 替换 pause
sudo docker pull registry.aliyuncs.com/google_containers/pause:3.5
sudo docker tag registry.aliyuncs.com/google_containers/pause:3.5 k8s.gcr.io/pause:3.5
sudo docker rmi registry.aliyuncs.com/google_containers/pause:3.5

# 替换 etcd
sudo docker pull registry.aliyuncs.com/google_containers/etcd:3.5.0-0
sudo docker tag registry.aliyuncs.com/google_containers/etcd:3.5.0-0 k8s.gcr.io/etcd:3.5.0-0
sudo docker rmi registry.aliyuncs.com/google_containers/etcd:3.5.0-0

# 替换 coredns
sudo docker pull coredns/coredns:1.8.4
sudo docker tag coredns/coredns:1.8.4 k8s.gcr.io/coredns/coredns:v1.8.4
sudo docker rmi coredns/coredns:1.8.4
```

### 初始化 master 节点

> 提示：在虚拟机：`k8s-1`中执行命令。

```shell
sudo kubeadm init --kubernetes-version=v1.22.0 --apiserver-advertise-address=192.168.205.10  --pod-network-cidr=10.244.0.0/16
# sudo kubeadm reset # 尽最大努力还原通过 'kubeadm init' 或者 'kubeadm join' 操作对主机所做的更改
```

根据返回的提示：设置：

```shell
# To start using your cluster, you need to run the following as a regular user:

mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

最后返回 join 的信息：

```shell
# Then you can join any number of worker nodes by running the following on each as root:

kubeadm join 192.168.205.10:6443 --token g012n6.65ete4bw7ys92tuv \
        --discovery-token-ca-cert-hash sha256:fdae044c194ed166f7b1b0746f5106008660ede517dd4cf436dfe68cc446c878
```

### node 加入集群

> 提示：在虚拟机：`k8s-2` 和 `k8s-3`中执行命令。

两个参数 token 和 hash 的值替换成你自己集群返回的值：

```shell
sudo kubeadm join 192.168.205.10:6443 --token g012n6.65ete4bw7ys92tuv \
        --discovery-token-ca-cert-hash sha256:fdae044c194ed166f7b1b0746f5106008660ede517dd4cf436dfe68cc446c878
# sudo kubeadm reset # 尽最大努力还原通过 'kubeadm init' 或者 'kubeadm join' 操作对主机所做的更改
```

### 安装 Pod 网络附加组件

> 提示：在虚拟机：`k8s-1`中执行命令。

这里选择 [calico](https://docs.projectcalico.org/getting-started/kubernetes/self-managed-onprem/onpremises)：

不选择 flannel 的原因是 flannel 暂不支持 NetworkPolicy 等新特性。

```shell
kubectl apply -f https://docs.projectcalico.org/manifests/calico.yaml
```

等待约 1 分钟，Node 的状态会从 NotReady->Ready

### 查看 node 状态

> 提示：在虚拟机：`k8s-1`中执行命令。

```shell
kubectl get nodes
```

返回：

```shell
NAME    STATUS   ROLES                  AGE   VERSION
k8s-1   Ready    control-plane,master   12m   v1.22.0
k8s-2   Ready    <none>                 12m   v1.22.0
k8s-3   Ready    <none>                 11m   v1.22.0
```

### 设置 KUBELET_EXTRA_ARGS

`unable to upgrade connection: pod does not exist`的[解决方案](https://github.com/kubernetes/kubernetes/issues/63702)：

> 提示：在虚拟机：`k8s-1`中执行命令。

```shell
cat <<EOF | sudo tee /etc/default/kubelet
KUBELET_EXTRA_ARGS="--node-ip=192.168.205.10"
EOF
sudo systemctl restart kubelet
```

> 提示：在虚拟机：`k8s-2`中执行命令。

```shell
cat <<EOF | sudo tee /etc/default/kubelet
KUBELET_EXTRA_ARGS="--node-ip=192.168.205.11"
EOF
sudo systemctl restart kubelet
```

> 提示：在虚拟机：`k8s-3`中执行命令。

```shell
cat <<EOF | sudo tee /etc/default/kubelet
KUBELET_EXTRA_ARGS="--node-ip=192.168.205.12"
EOF
sudo systemctl restart kubelet
```

> 提示：在虚拟机：`k8s-1`中执行命令。

```shell
kubectl get nodes -o wide
```

返回：

```shell
NAME    STATUS   ROLES                  AGE   VERSION   INTERNAL-IP      EXTERNAL-IP   OS-IMAGE             KERNEL-VERSION       CONTAINER-RUNTIME
k8s-1   Ready    control-plane,master   49m   v1.22.0   192.168.205.10   <none>        Ubuntu 18.04.5 LTS   4.15.0-144-generic   docker://20.10.8
k8s-2   Ready    <none>                 48m   v1.22.0   192.168.205.11   <none>        Ubuntu 18.04.5 LTS   4.15.0-144-generic   docker://20.10.8
k8s-3   Ready    <none>                 37m   v1.22.0   192.168.205.12   <none>        Ubuntu 18.04.5 LTS   4.15.0-144-generic   docker://20.10.8
```

### 清理

销毁虚拟机：

```shell
vagrant destroy
```

### 参考

* [Vagrant box ubuntu/bionic64](https://app.vagrantup.com/ubuntu/boxes/bionic64)
* [Install Docker Engine on Ubuntu](https://docs.docker.com/engine/install/ubuntu/)
* [Installing kubeadm](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/)
* [Configure the Docker daemon, in particular to use systemd for the management of the container’s cgroups](https://kubernetes.io/docs/setup/production-environment/container-runtimes/#docker)
* [Creating a cluster with kubeadm](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/)
* [Cluster Networking](https://kubernetes.io/docs/concepts/cluster-administration/networking/#how-to-implement-the-kubernetes-networking-model)
* [阿里云 kubernetes yum 仓库镜像](https://developer.aliyun.com/article/433817)
* [让 K8S 在 GFW 内愉快的航行](https://developer.aliyun.com/article/759310)