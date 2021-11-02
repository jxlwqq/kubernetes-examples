# 自动扩缩

访问：

```shell
curl http://http-echo.default.127.0.0.1.sslip.io
```

另开一个窗口，观察 Pod：
```shell
kubectl get pod -l serving.knative.dev/service=http-echo --watch
```

watch 返回：

> 扩缩至零，过程可能需要1-2分钟时间。

```shell
NAME                                      READY   STATUS        RESTARTS   AGE
http-echo-v1-deployment-948b9746b-275j2   2/2     Running       0          110s
http-echo-v1-deployment-948b9746b-275j2   2/2     Terminating   0          118s
http-echo-v1-deployment-948b9746b-275j2   1/2     Terminating   0          2m1s
http-echo-v1-deployment-948b9746b-275j2   0/2     Terminating   0          2m29s
http-echo-v1-deployment-948b9746b-275j2   0/2     Terminating   0          2m30s
http-echo-v1-deployment-948b9746b-275j2   0/2     Terminating   0          2m30s
```

再次访问：

```shell
curl http://http-echo.default.127.0.0.1.sslip.io
```

watch 返回：

> 迅速启动。

```shell
NAME                                      READY   STATUS              RESTARTS   AGE
http-echo-v1-deployment-948b9746b-275j2   2/2     Running             0          110s
http-echo-v1-deployment-948b9746b-275j2   2/2     Terminating         0          118s
http-echo-v1-deployment-948b9746b-275j2   1/2     Terminating         0          2m1s
http-echo-v1-deployment-948b9746b-275j2   0/2     Terminating         0          2m29s
http-echo-v1-deployment-948b9746b-275j2   0/2     Terminating         0          2m30s
http-echo-v1-deployment-948b9746b-275j2   0/2     Terminating         0          2m30s
http-echo-v1-deployment-948b9746b-5vwgt   0/2     Pending             0          0s
http-echo-v1-deployment-948b9746b-5vwgt   0/2     Pending             0          0s
http-echo-v1-deployment-948b9746b-5vwgt   0/2     ContainerCreating   0          0s
http-echo-v1-deployment-948b9746b-5vwgt   1/2     Running             0          2s
http-echo-v1-deployment-948b9746b-5vwgt   2/2     Running             0          2s
```