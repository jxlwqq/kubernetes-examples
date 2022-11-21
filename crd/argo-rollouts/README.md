# Argo Rollouts

Argo Rollouts（https://argoproj.github.io/argo-rollouts/） 是一个 Kubernetes Controller 和一组 CRD，为 Kubernetes 提供高级部署功能，例如蓝绿、金丝雀、金丝雀分析、实验和渐进式交付功能。

原生 Kubernetes Deployment 对象支持 RollingUpdate 策略，该策略在更新期间提供一组基本的安全保证（就绪探测）。 然而滚动更新策略面临许多限制：

无法控制流量到新版本
无法查询外部指标来验证更新
可以停止部署进程，但无法自动中止和回滚更新
由于这些原因，在大规模的大批量生产环境中，滚动更新通常被认为更新过程的风险太大，因为它无法控制爆炸半径，可能会过于激进地 Rollout，并且不会在失败时提供自动回滚。

#### 安装

```bash
kubectl create namespace argo-rollouts
kubectl apply -n argo-rollouts -f https://github.com/argoproj/argo-rollouts/releases/latest/download/install.yaml
 
brew install argoproj/tap/kubectl-argo-rollouts
kubectl argo rollouts version
```
