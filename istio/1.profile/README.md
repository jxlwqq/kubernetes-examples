# 配置文件

在上一节安装环节中，我们在命令行中增加了 `--set profile=demo` 参数：

```shell
istioctl install --set profile=demo -y
```

profile 是 istioctl 内置的安装配置文件，有以下几个选择：

1. default：根据 IstioOperator API 的默认设置启动组件。 建议用于生产部署和 Multicluster Mesh 中的 Primary Cluster。
2. demo：这一配置具有适度的资源需求，旨在展示 Istio 的功能。 它适合运行 Bookinfo 应用程序和相关任务。 这是通过快速开始指导安装的配置。 
3. minimal：与默认配置文件相同，但只安装了控制平面组件。 它允许您使用 Separate Profile 配置控制平面和数据平面组件(例如 Gateway)。 
4. remote：配置 Multicluster Mesh 的 Remote Cluster。 
5. empty：不部署任何东西。可以作为自定义配置的基本配置文件。 
6. preview：预览文件包含的功能都是实验性。这是为了探索 Istio 的新功能。不确保稳定性、安全性和性能（使用风险需自负）。

本地开发环境一般选择 `demo`，生产部署选择 `default`。


标注 &#x2714; 的组件安装在每个配置文件中：

|     | default | demo | minimal | remote | empty | preview |
| --- | --- | --- | --- | --- | --- | --- |
| 核心组件 | | | | | | | |
| `istio-egressgateway` | | &#x2714; | | | | | | |
| `istio-ingressgateway` | &#x2714; | &#x2714; | | | | &#x2714; |
| `istiod` | &#x2714; | &#x2714; | &#x2714; | | | &#x2714; |

