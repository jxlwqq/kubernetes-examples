# Kubernetes Ingress with Nginx 示例

#### 实验前提

* 需要你有 macOS 开发环境，本文以此为例，其他类型的开发环境请自行搭建。
* 需要你对 YAML 这一专门用来写配置文件的语言有所了解。
* 需要你对 Docker 有一些基本的了解。
* 需要你对 Kubernetes 中的 Node、Pod、ReplicaSet、Deployment、Service、Ingress、ConfigMap 等一些核心基础概念有一定的了解。

#### 安装 Docker for Mac

下载地址：https://hub.docker.com/editions/community/docker-ce-desktop-mac

启动并开启 Kubernetes 功能，功能开启过程中，Docker将会自动拉取相关镜像，所以全程需要科学上网。

#### 镜像准备

在等待 Kubernetes 的状态变成 running 的同时，我们可以手动拉取以下两个镜像，避免后续操作等待时间过长。

```bash
docker pull quay.io/kubernetes-ingress-controller/nginx-ingress-controller:0.30.0
docker pull hashicorp/http-echo
```

#### 本地端口准备

请确保本地localhost的80端口没有被占用，已在使用的请在实验期间暂时关闭占用80端口的服务。

#### 切换集群

如果你本地有多个 Kubernetes 的集群配置，请先切换至名为 docker-desktop 的集群：

````bash
kubectl config use-context docker-desktop
````

#### 创建 Ingress 控制器

这里我们选择 ingress-nginx

```bash
kubectl apply -f ../ingress-nginx/ingress-nginx-deployment-and-other-resources-mandatory.yaml
kubectl apply -f ../ingress-nginx/ingress-nginx-service.yaml
```

注：
* ingress-nginx-deployment-and-other-resources-mandatory.yaml 文件来源自：https://github.com/kubernetes/ingress-nginx/blob/master/deploy/static/mandatory.yaml
* ingress-nginx-service.yaml 文件来源自：https://github.com/kubernetes/ingress-nginx/blob/master/deploy/static/provider/cloud-generic.yaml

详细操作说明见：https://github.com/kubernetes/ingress-nginx/blob/master/docs/deploy/index.md

#### 创建 Apple 资源（kind: Deployment 与 kind: Service）

```bash
kubectl apply -f apple-deployment-and-service.yaml
```

#### 创建 Banana 资源（kind: Deployment 与 kind: Service）

```bash
kubectl apply -f apple-deployment-and-service.yaml
```

#### 创建 Ingress 资源（kind: Ingress）

创建 Ingress 路由规则：

```bash
kubectl apply -f ingress.yaml
```

#### 查看资源状态
```bash
kubectl get all --all-namespaces
```

#### 访问

```bash
curl http://localhost/apple  # return apple
curl http://localhost/banana # return banana
curl http://localhost/pear   # return 404
```

结束，撒花。

#### 尾巴

如果你想使用 minikube 作为本地集群的话，可以使用以下命令安装：

minikube 的安装需要一些外部的依赖组件：详见阿里云的这篇文章：[Minikube - Kubernetes本地实验环境](https://yq.aliyun.com/articles/221687)

```bash
brew update
brew install minikube
minikube delete && rm -fr ~/.minikube/
minikube start --vm-driver=virtualbox \
               --image-mirror-country cn \
               --image-repository=registry.cn-hangzhou.aliyuncs.com/google_containers \
               --iso-url=https://kubernetes.oss-cn-hangzhou.aliyuncs.com/minikube/iso/minikube-v1.7.3.iso \
               --registry-mirror=https://xxxxxx.mirror.aliyuncs.com
```