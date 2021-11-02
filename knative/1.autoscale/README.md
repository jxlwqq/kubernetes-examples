# 自动扩缩

访问：

```shell
curl http://hello.default.127.0.0.1.sslip.io
```

另开一个窗口，观察 Pod：
```shell
kubectl get pod -l serving.knative.dev/service=hello --watch
```

watch 返回：

> 扩缩至零，过程可能需要1-2分钟时间。

```shell
NAME                                      READY   STATUS        RESTARTS   AGE
hello-world-deployment-55f4fb96b4-qzjsl   2/2     Running       0          21s
hello-world-deployment-55f4fb96b4-qzjsl   2/2     Terminating   0          115s
hello-world-deployment-55f4fb96b4-qzjsl   1/2     Terminating   0          117s
hello-world-deployment-55f4fb96b4-qzjsl   0/2     Terminating   0          2m26s
hello-world-deployment-55f4fb96b4-qzjsl   0/2     Terminating   0          2m27s
hello-world-deployment-55f4fb96b4-qzjsl   0/2     Terminating   0          2m27s
```

再次访问：

```shell
curl http://hello.default.127.0.0.1.sslip.io
```

watch 返回：

> 迅速启动。

```shell
NAME                                      READY   STATUS              RESTARTS   AGE
hello-world-deployment-55f4fb96b4-qzjsl   2/2     Running             0          21s
hello-world-deployment-55f4fb96b4-qzjsl   2/2     Terminating         0          115s
hello-world-deployment-55f4fb96b4-qzjsl   1/2     Terminating         0          117s
hello-world-deployment-55f4fb96b4-qzjsl   0/2     Terminating         0          2m26s
hello-world-deployment-55f4fb96b4-qzjsl   0/2     Terminating         0          2m27s
hello-world-deployment-55f4fb96b4-qzjsl   0/2     Terminating         0          2m27s
hello-world-deployment-55f4fb96b4-2mppf   0/2     Pending             0          0s
hello-world-deployment-55f4fb96b4-2mppf   0/2     Pending             0          0s
hello-world-deployment-55f4fb96b4-2mppf   0/2     ContainerCreating   0          0s
hello-world-deployment-55f4fb96b4-2mppf   1/2     Running             0          2s
hello-world-deployment-55f4fb96b4-2mppf   2/2     Running             0          2s
```