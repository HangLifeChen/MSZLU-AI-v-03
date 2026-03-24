# Mac本地K8s部署指南

## 当前部署状态

### 已成功部署的组件：
✅ PVC存储（使用hostpath）
✅ RBAC权限配置
✅ 数据库服务（PostgreSQL, Redis, Elasticsearch, Etcd, MinIO, Milvus）
✅ 应用配置（ConfigMap）
✅ 应用服务（Service）
✅ HPA自动扩缩容

### 需要解决的问题：
❌ 应用Pod无法启动 - 镜像拉取失败
- 原因：无法访问私有Harbor仓库 `192.168.8.101:30002/mszlu-ai/app:v1.0.0`
- 缺少`harbor-secret`镜像拉取密钥

❌ Ingress部署失败
- 原因：Nginx Ingress Controller webhook连接失败

## 解决方案

### 方案1：使用私有Harbor仓库（推荐用于生产环境）

如果你有访问Harbor仓库的权限，需要：

1. 创建imagePullSecret：
```bash
kubectl create secret docker-registry harbor-secret \
  --docker-server=192.168.8.101:30002 \
  --docker-username=<你的用户名> \
  --docker-password=<你的密码> \
  --docker-email=<你的邮箱> \
  -n mszlu-ai
```

2. 重新部署应用：
```bash
kubectl delete deployment mszlu-ai-app -n mszlu-ai
kubectl apply -f k8s/app-deployments.yml
```

### 方案2：构建本地镜像（推荐用于本地开发）

1. 构建应用镜像：
```bash
# 在项目根目录执行
cd app
docker build -t mszlu-ai-app:local .
```

2. 修改deployment配置，将镜像地址改为本地镜像：
```bash
kubectl set image deployment/mszlu-ai-app app=mszlu-ai-app:local -n mszlu-ai
```

### 方案3：使用NodePort访问应用（临时方案）

如果Ingress不可用，可以使用NodePort或port-forward访问应用：

```bash
# 方法1：使用port-forward
kubectl port-forward svc/mszlu-ai-app-service 8888:80 -n mszlu-ai

# 然后访问 http://localhost:8888

# 方法2：修改Service为NodePort类型
kubectl patch svc mszlu-ai-app-service -n mszlu-ai -p '{"spec":{"type":"NodePort"}}'
```

## 验证部署状态

### 查看所有Pod状态：
```bash
kubectl get pods -n mszlu-ai
```

### 查看服务状态：
```bash
kubectl get svc -n mszlu-ai
```

### 查看应用日志：
```bash
kubectl logs -f deployment/mszlu-ai-app -n mszlu-ai
```

### 查看特定Pod的详细信息：
```bash
kubectl describe pod <pod-name> -n mszlu-ai
```

## 数据库连接信息

部署成功后，数据库可以通过以下方式访问：

### PostgreSQL:
- Service: `postgres-service.mszlu-ai.svc.cluster.local:5432`
- User: `mszluai`
- Password: `mszlu123456`
- Database: `faber-ai`

### Redis:
- Service: `redis-service.mszlu-ai.svc.cluster.local:6379`

### Elasticsearch:
- Service: `elasticsearch-service.mszlu-ai.svc.cluster.local:9200`
- Username: `elastic` (如果启用安全功能)
- Password: `mszlu123456!@#$`

### Milvus:
- Service: `milvus-service.mszlu-ai.svc.cluster.local:19530`

### MinIO:
- Service: `minio-service.mszlu-ai.svc.cluster.local:9000`
- Access Key: `minioadmin`
- Secret Key: `minioadmin`

## 清理部署

如果需要清理所有资源：
```bash
kubectl delete namespace mszlu-ai
```

## 下一步

请选择以下方案之一继续：

1. **如果你有Harbor访问权限**：提供Harbor的用户名和密码，我会帮你创建imagePullSecret
2. **如果你想构建本地镜像**：我需要查看是否有Dockerfile，或者帮你创建一个
3. **如果你想先测试数据库**：数据库已经部署成功，可以先验证数据库连接

请告诉我你想选择哪个方案，我会继续帮你完成部署。
