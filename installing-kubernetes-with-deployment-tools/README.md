# 使用 kubeadm 创建集群

本地开发环境以 macOS 为例。

#### 安装 VirtualBox 和 Vagrant

* 安装 [VirtualBox](https://www.virtualbox.org/)
* 安装 [Vagrant](https://www.vagrantup.com

#### 启动虚拟机

虚拟机配置了3台，分别是 `k8s-1(192.168.205.10)`、`k8s-2(192.168.205.11)` 和 `k8s-3(192.168.205.11)`。配置详见 Vagrantfile 这个文件。

虚拟机初始化的时候，已经帮助你安装了 Docker 环境，详见 `config.vm.provision "shell"` 中信息。Vagrant 是用 Ruby 写的，语法都是通用的，应该能看懂。

```shell
cd installing-kubernetes-with-deployment-tools # 进入这个目录
vagrant box add centos/7 # 提前下载操作系统镜像文件，方便后续快速启动
vagrant up # 启动虚拟机
```

#### 登录虚拟机

```shell
vagrant ssh k8s-1 # 这台作为 master
vagrant ssh k8s-2 # node
vagrant ssh k8s-3 # node
```

#### 允许 iptables 检查桥接流量

所在虚拟机：

```shell
vagrant ssh k8s-1
```

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

#### 安装 kubeadm、kubelet 和 kubectl

所在虚拟机：

```shell
vagrant ssh k8s-1
```

使用阿里云的镜像进行安装：

```shell
sudo cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=http://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=0
repo_gpgcheck=0
gpgkey=http://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg
       http://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
EOF

# 将 SELinux 设置为 permissive 模式（相当于将其禁用）
sudo setenforce 0
sudo sed -i 's/^SELINUX=enforcing$/SELINUX=permissive/' /etc/selinux/config

# 关闭 swap
sudo swapoff -a

sudo yum install -y kubelet kubeadm kubectl --disableexcludes=kubernetes

sudo systemctl enable --now kubelet
```

#### 配置 cgroup 驱动程序

所在虚拟机：

```shell
vagrant ssh k8s-1
```

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

sudo systemctl restart docker
```

#### 验证与 gcr.io 容器镜像仓库的连通性

所在虚拟机：

```shell
vagrant ssh k8s-1
```

使用 `kubeadm config images pull` 验证与 gcr.io 容器镜像仓库的连通性，不过会失败。

我们通过拉取阿里云的Docker镜像仓库，然后在打 Tag 的方式，来曲线解决这个问题。

```shell
kubeadm config images list # 查看所需的镜像

# 替换 kube-apiserver
sudo docker pull registry.aliyuncs.com/google_containers/kube-apiserver:v1.21.0
sudo docker tag registry.aliyuncs.com/google_containers/kube-apiserver:v1.21.0 k8s.gcr.io/kube-apiserver:v1.21.0
sudo docker rmi registry.aliyuncs.com/google_containers/kube-apiserver:v1.21.0

# 替换 kube-controller-manager
sudo docker pull registry.aliyuncs.com/google_containers/kube-controller-manager:v1.21.0
sudo docker tag registry.aliyuncs.com/google_containers/kube-controller-manager:v1.21.0 k8s.gcr.io/kube-controller-manager:v1.21.0
sudo docker rmi registry.aliyuncs.com/google_containers/kube-controller-manager:v1.21.0

# 替换 kube-scheduler
sudo docker pull registry.aliyuncs.com/google_containers/kube-scheduler:v1.21.0
sudo docker tag registry.aliyuncs.com/google_containers/kube-scheduler:v1.21.0 k8s.gcr.io/kube-scheduler:v1.21.0
sudo docker rmi registry.aliyuncs.com/google_containers/kube-scheduler:v1.21.0

# 替换 kube-proxy
sudo docker pull registry.aliyuncs.com/google_containers/kube-proxy:v1.21.0
sudo docker tag registry.aliyuncs.com/google_containers/kube-proxy:v1.21.0 k8s.gcr.io/kube-proxy:v1.21.0
sudo docker rmi registry.aliyuncs.com/google_containers/kube-proxy:v1.21.0

# 替换 pause
sudo docker pull registry.aliyuncs.com/google_containers/pause:3.4.1
sudo docker tag registry.aliyuncs.com/google_containers/pause:3.4.1 k8s.gcr.io/pause:3.4.1
sudo docker rmi registry.aliyuncs.com/google_containers/pause:3.4.1

# 替换 etcd
sudo docker pull registry.aliyuncs.com/google_containers/etcd:3.4.13-0
sudo docker tag registry.aliyuncs.com/google_containers/etcd:3.4.13-0 k8s.gcr.io/etcd:3.4.13-0
sudo docker rmi registry.aliyuncs.com/google_containers/etcd:3.4.13-0

# 替换 coredns
sudo docker pull coredns/coredns:1.8.0
sudo docker tag coredns/coredns:1.8.0 k8s.gcr.io/coredns/coredns:v1.8.0
sudo docker rmi coredns/coredns:1.8.0
```

#### 初始化 master 节点

所在虚拟机：

```shell
vagrant ssh k8s-1
```

```shell
sudo kubeadm init --kubernetes-version=v1.21.0 --apiserver-advertise-address=192.168.205.10  --pod-network-cidr=10.244.0.0/16
```

根据返回的提示：设置：

```shell
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

最后返回 join 的信息：

```shell
kubeadm join 192.168.205.10:6443 --token g012n6.65ete4bw7ys92tuv \
        --discovery-token-ca-cert-hash sha256:fdae044c194ed166f7b1b0746f5106008660ede517dd4cf436dfe68cc446c878
```

#### node 加入集群

所在虚拟机：

```shell
vagrant ssh k8s-2
vagrant ssh k8s-3
```

```shell
kubeadm join 192.168.205.10:6443 --token g012n6.65ete4bw7ys92tuv \
        --discovery-token-ca-cert-hash sha256:fdae044c194ed166f7b1b0746f5106008660ede517dd4cf436dfe68cc446c878
```

#### 安装 Pod 网络附加组件

所在虚拟机：

```shell
vagrant ssh k8s-1
```

这里选择 flannel：

```shell
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
```

#### 查看 node 状态

所在虚拟机：

```shell
vagrant ssh k8s-1
```

```shell
kubectl get nodes
```

返回：

```shell
NAME    STATUS   ROLES                  AGE   VERSION
k8s-1   Ready    control-plane,master   12m   v1.21.0
k8s-2   Ready    <none>                 12m   v1.21.0
k8s-3   Ready    <none>                 11m   v1.21.0
```