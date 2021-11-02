# 流量分流

### 更新 Hello 服务
```shell
kubectl apply -f http-echo-update.yaml
```

访问：
```shell
curl http://http-echo.default.127.0.0.1.sslip.io
```

返回：
```shell
v2
```

### 查看修订历史

```shell
kn revisions list
```

返回：

```shell
NAME           SERVICE     TRAFFIC   TAGS   GENERATION   AGE    CONDITIONS   READY   REASON
http-echo-v2   http-echo   100%             2            25s    4 OK / 4     True    
http-echo-v1   http-echo                    1            9m5s   3 OK / 4     True   
```

可以看出，http-echo-v2 版本分流了所有流量。

### 分流

```shell
kubectl apply -f http-echo-split.yaml
```

再次查看修订历史：

```shell
kn revisions list
```

```shell
NAME           SERVICE     TRAFFIC   TAGS   GENERATION   AGE     CONDITIONS   READY   REASON
http-echo-v2   http-echo   50%              2            61s     4 OK / 4     True    
http-echo-v1   http-echo   50%              1            9m41s   3 OK / 4     True  
```

现在变更为了各 50%。

再次访问：
```shell
curl http://http-echo.default.127.0.0.1.sslip.io
```

返回：

```shell
v1
```

或者：

```shell
v2
```
