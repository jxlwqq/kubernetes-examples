# 高级流量路由

本示例为一个简单的 Customer List 微服务应用，由 [Tetrate](https://academy.tetrate.io/) 创建。访问 Web 页面，将展示顾客信息。


我们将部署 Web 前端、Customers v1、Customers v2，以及相应的 Gateway、 VirtualServices 和 DestinationRule。其中 Customers v1 仅返回顾客姓名，而 Customers v2 返回顾客的姓名和所在城市。


### 部署网关

```shell
kubectl apply -f istio-gateway.yaml
```

信息如下：

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: sample-gateway
spec:
  selector:
    istio: ingressgateway
  servers:
    - port:
        number: 80
        name: http
        protocol: HTTP
      hosts:
        - 'web.example.com'
        - 'svc.example.com'
```

通过 `curl -H "Host: web.example.com" 127.0.0.1`，我们可以访问 Web 前端服务。

通过 `curl -H "Host: svc.example.com" 127.0.0.1`，我们可以访问 Customers 后端服务。


### 部署 Customers 后端服务

Customers 后端服务同时对服务网格内部和外部暴露服务。

部署 Deployment 和 Service：

```shell
kubectl apply -f customers-deployment-v1.yaml
kubectl apply -f customers-deployment-v2.yaml
kubectl apply -f customers-service.yaml
```

部署 VirtualService 和 DestinationRule：

```shell
kubectl apply -f customers-virtual-service.yaml
kubectl apply -f customers-destination-rule.yaml
```

信息如下：

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: customers
spec:
  hosts:
    - 'svc.example.com' # 集群外访问
    - 'customers.default.svc.cluster.local' # 集群内访问
  gateways:
    - sample-gateway # 集群外访问
    - mesh # 集群内访问
  http:
    - match: # 匹配到集群内网关的访问，**并且**HEAD信息的用户是 debug
        - gateways:
            - mesh
          headers:
            user:
              exact: debug
      route:
        - destination:
            host: customers.default.svc.cluster.local
            port:
              number: 80
            subset: v2
    - match: # 匹配到集群内网关的访问
        - gateways:
            - mesh
      route:
        - destination:
            host: customers.default.svc.cluster.local
            port:
              number: 80
            subset: v1
    - match: # 匹配是集群外的访问（即边缘网关）
        - gateways:
            - sample-gateway
      route:
        - destination:
            host: customers.default.svc.cluster.local
            port:
              number: 80
            subset: v2
---
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: customers
spec:
  host: customers.default.svc.cluster.local
  subsets:
    - name: v1
      labels:
        version: v1
    - name: v2
      labels:
        version: v2
```

### 部署 Web 前端服务

Web 前端服务只对集群外部暴露服务。

部署 Deployment 和 Service：

```shell
kubectl apply -f webfrontend-deployment.yaml
kubectl apply -f webfrontend-service.yaml
```


部署 VirtualService：

```shell
kubectl apply -f webfrontend-virtual-service.yaml
```

信息如下：

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: web-frontend
spec:
  hosts:
    - 'web.example.com' # 集群外访问
  gateways:
    - sample-gateway # 集群外访问
  http:
    - route:
        - destination:
            host: web-frontend.default.svc.cluster.local
            port:
              number: 80
```

### 访问

```shell
curl -H "Host: web.example.com"  -H "User: debug" 127.0.0.1 # 访问 web 页面，展示城市和用户信息
curl -H "Host: web.example.com"  -H "User: abc" 127.0.0.1   # 访问 web 页面，仅展示用户信息
curl -H "Host: svc.example.com" 127.0.0.1 # 直接访问后端服务，返回 JSON 数据，包含城市和用户信息
```










