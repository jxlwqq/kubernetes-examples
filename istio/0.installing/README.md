# 快速入门

### 前提

* 启动 Docker For Mac，并开启 Kubernetes 功能；
* 如果 Istio 相关组件在集群中无法正常运行，需要在 Docker 首选项的 Advanced 面板下增加 Docker 的 CPU 和内存限制；
* （可选）提前拉取容器镜像，提升 Pod 启动速度：

```shell
# 核心组件
docker pull docker.io/istio/pilot
docker pull docker.io/istio/proxyv2
# 相关插件
docker pull docker.io/grafana/grafana
docker pull docker.io/jaegertracing/all-in-one
docker pull docker.io/prom/prometheus
docker pull docker.io/jimmidyson/configmap-reload
docker pull quay.io/kiali/kiali
# bookinfo 示例
docker pull docker.io/istio/examples-bookinfo-details-v1
docker pull docker.io/istio/examples-bookinfo-ratings-v1
docker pull docker.io/istio/examples-bookinfo-reviews-v1
docker pull docker.io/istio/examples-bookinfo-reviews-v2
docker pull docker.io/istio/examples-bookinfo-reviews-v3
docker pull docker.io/istio/examples-bookinfo-productpage-v1
```

### 最小步骤

一共分为 9 个步骤，如下所示：

```shell
# 第1步：安装 istioctl
brew install istioctl
# 第2步：切换至目标集群
kubectl config use-context docker-desktop
# 第3步：在目标集群中安装 istio
istioctl install --set profile=demo -y
# 第4步：自动注入 Envoy 边车代理
kubectl label namespace default istio-injection=enabled
# 第5步：获取示例代码
git clone git@github.com:istio/istio.git && cd istio
# 第6步：部署 Bookinfo 示例应用
kubectl apply -f samples/bookinfo/platform/kube/bookinfo.yaml
# 第7步：把应用关联到 Istio 网关
kubectl apply -f samples/bookinfo/networking/bookinfo-gateway.yaml
# 第8步：安装 Kiali、Prometheus、Grafana、Jaeger 等插件
kubectl apply -f samples/addons
# 第9步：访问应用
curl http://127.0.0.1/productpage
```

> ⚠️  Apple chip 用户注意，需要 patch 下 istio-egressgateway 和 istio-ingressgateway 这两个 Deployment 资源，它们的默认节点亲和性只有 `amd64`，需要手动新增 `arm64` 这个值，否则 Pod 将一直处于 Pending 状态。[详见](https://github.com/istio/istio/issues/21094)

解决方案：
```shell
# 使用社区（非官方）构建的 arm64 镜像（https://github.com/querycap/istio）：
istioctl install --set hub=docker.io/querycapistio --set profile=demo -y

# 更新 istio-egressgateway 和 istio-ingressgateway 资源
kubectl patch deployments.apps \
  istio-ingressgateway \
  --namespace istio-system \
  --type='json' \
  -p='[
  {"op": "replace", "path": "/spec/template/spec/affinity/nodeAffinity/preferredDuringSchedulingIgnoredDuringExecution/0/preference/matchExpressions/0/values", "value": [amd64,arm64]},
  {"op": "replace", "path": "/spec/template/spec/affinity/nodeAffinity/requiredDuringSchedulingIgnoredDuringExecution/nodeSelectorTerms/0/matchExpressions/0/values", "value": [amd64,arm64,ppc64le,s390x]}
  ]'
  
kubectl patch deployments.apps \
  istio-egressgateway \
  --namespace istio-system \
  --type='json' \
  -p='[
  {"op": "replace", "path": "/spec/template/spec/affinity/nodeAffinity/preferredDuringSchedulingIgnoredDuringExecution/0/preference/matchExpressions/0/values", "value": [amd64,arm64]},
  {"op": "replace", "path": "/spec/template/spec/affinity/nodeAffinity/requiredDuringSchedulingIgnoredDuringExecution/nodeSelectorTerms/0/matchExpressions/0/values", "value": [amd64,arm64,ppc64le,s390x]}
  ]'
```

### 命令行自动补全

Zsh 用户：

```
mkdir -p ~/completions && istioctl collateral --zsh -o ~/completions
echo "source ~/completions/_istioctl" >> ~/.zshrc
```

`~/completions` 目录用于存放补全提示文件，可自定义。

### 查看仪表板

打入流量：

```shell
while true; do curl http://127.0.0.1/productpage; done
```

访问 Kiali 仪表板:

```shell
istioctl dashboard kiali
```

### 清理 bookinfo

```shell
kubectl delete -f samples/bookinfo/platform/kube/bookinfo.yaml
kubectl delete -f samples/bookinfo/networking/bookinfo-gateway.yaml
```

Istio 的组件和插件暂时不要清理。后续实验需要接着用。

#### 参考：

* [Istio 文档：平台安装 Docker Desktop](https://istio.io/latest/zh/docs/setup/platform-setup/docker/)
* [Istio 文档：入门](https://istio.io/latest/zh/docs/setup/getting-started/)
* [什么是微服务？](https://www.redhat.com/zh/topics/microservices/what-are-microservices)
* [什么是服务网格？](https://www.redhat.com/zh/topics/microservices/what-is-a-service-mesh)
* [什么是 Istio？](https://www.redhat.com/zh/topics/microservices/what-is-istio)
* [什么是 Jaeger？原理，用途，组件及相关术语盘点](https://www.redhat.com/zh/topics/microservices/what-is-jaeger)
