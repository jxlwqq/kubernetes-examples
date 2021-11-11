# 在线精品店

[Online Boutique](https://github.com/GoogleCloudPlatform/microservices-demo) 是一个云原生微服务演示应用程序。Online Boutique 是一个由 10 个微服务组成的应用。该应用是一个基于 Web 的电子商务应用，用户可以浏览商品，将其添加到购物车，并购买商品。

### 预先拉取镜像

```shell
docker pull gcr.io/google-samples/microservices-demo/emailservice:v0.3.0
docker pull gcr.io/google-samples/microservices-demo/checkoutservice:v0.3.0
docker pull gcr.io/google-samples/microservices-demo/recommendationservice:v0.3.0
docker pull gcr.io/google-samples/microservices-demo/frontend:v0.3.0
docker pull gcr.io/google-samples/microservices-demo/paymentservice:v0.3.0
docker pull gcr.io/google-samples/microservices-demo/productcatalogservice:v0.3.0
docker pull gcr.io/google-samples/microservices-demo/cartservice:v0.3.0
docker pull gcr.io/google-samples/microservices-demo/loadgenerator:v0.3.0
docker pull gcr.io/google-samples/microservices-demo/currencyservice:v0.3.0
docker pull gcr.io/google-samples/microservices-demo/shippingservice:v0.3.0
docker pull gcr.io/google-samples/microservices-demo/adservice:v0.3.0
```

### 部署服务

```shell
git clone git@github.com:GoogleCloudPlatform/microservices-demo.git
cd microservices-demo
kubectl apply -f release/kubernetes-manifests.yaml
kubectl apply -f release/istio-manifests.yaml
```