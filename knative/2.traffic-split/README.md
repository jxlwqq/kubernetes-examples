# 流量分流

### 更新 Hello 服务
```shell
kubectl apply -f hello-update.yaml
```

访问：
```shell
curl http://hello.default.127.0.0.1.sslip.io
```

返回：
```shell
Hello Knative!
```

### 查看修订历史

```shell
kn revisions list
```

返回：

```shell
NAME            SERVICE   TRAFFIC   TAGS   GENERATION   AGE     CONDITIONS   READY   REASON
hello-knative   hello     100%             2            2m25s   4 OK / 4     True
hello-world     hello                      1            38m     3 OK / 4     True
```

可以看出，hello-knative 版本分流了所有流量。

### 分流

```shell
kubectl apply -f hello-split.yaml
```

再次查看修订历史：

```shell
kn revisions list
```

```shell
NAME            SERVICE   TRAFFIC   TAGS   GENERATION   AGE     CONDITIONS   READY   REASON
hello-knative   hello     50%              2            7m13s   3 OK / 4     True
hello-world     hello     50%              1            43m     3 OK / 4     True
```

现在变更为了各 50%。

再次访问：
```shell
curl http://hello.default.127.0.0.1.sslip.io
```

返回：

```shell
Hello Knative!
```

或者：

```shell
Hello World!
```
