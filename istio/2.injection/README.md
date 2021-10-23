# 注入

### 自动注入

```shell
kubectl label namespace default istio-injection=enabled
```

使用 Istio 提供的`准入控制器变更 Webhook`，可以将 Sidecar 自动添加到可用的 Kubernetes Pod 中。当你在一个命名空间中设置了 istio-injection=enabled 标签，且 Injection webhook 被启用后，任何新的 Pod 都有将在创建时自动添加 Sidecar。

查看**准入控制器变更 Webhook**:

```shell
kubectl get mutatingwebhookconfigurations
```

返回：
```shell
NAME                     WEBHOOKS   AGE
istio-sidecar-injector   4          23h
```

什么是**准入控制 Webhook**:

准入 Webhook 是一种用于接收准入请求并对其进行处理的 HTTP 回调机制。 可以定义两种类型的准入 webhook，即 **验证性质的准入 Webhook** 和 **变更性质的准入 Webhook**。 变更性质的准入 Webhook 会先被调用。它们可以更改发送到 kube-apiserver 的对象以执行自定义的设置默认值操作。

在完成了所有对象修改并且 kube-apiserver 也验证了所传入的对象之后， 验证性质的 Webhook 会被调用，并通过拒绝请求的方式来强制实施自定义的策略。

更多信息请参阅：

* [使用准入控制器](https://kubernetes.io/zh/docs/reference/access-authn-authz/admission-controllers/)
* [动态准入控制](https://kubernetes.io/zh/docs/reference/access-authn-authz/extensible-admission-controllers/)

### 手动注入

```shell
# 由于 default 命令空间已经被设置为自动注入，所以我们需要新建一个命名空间，试验手动注入
kubectl create namespace tutorial
# 新建一个 nginx Deployment
kubectl create deployment nginx --image=nginx --dry-run=client --output=yaml > nginx.yaml
kubectl apply -f nginx.yaml -n tutorial
# 手动注入
istioctl kube-inject -f nginx.yaml > nginx-injection.yaml
kubectl apply -f nginx-injection.yaml -n tutorial
# 比较差异
diff nginx.yaml nginx-injection.yaml
```