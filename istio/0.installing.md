# 快速入门

#### 前提

* 启动 Docker For Mac，并开启 Kubernetes 功能；
* 如果 istio 相关组件在集群中无法正常运行，需要在 Docker 首选项的 Advanced 面板下增加 Docker 的 CPU 和内存限制；
* 提前拉取容器镜像

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

#### 最小步骤

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

### 查看仪表板

打入流量：

```shell
while true; do curl http://127.0.0.1/productpage; done
```

访问 Kiali 仪表板:

```shell
istioctl dashboard kiali
```

#### 参考：

* [Istio 文档：平台安装 Docker Desktop](https://istio.io/latest/zh/docs/setup/platform-setup/docker/)
* [Istio 文档：入门](https://istio.io/latest/zh/docs/setup/getting-started/)
