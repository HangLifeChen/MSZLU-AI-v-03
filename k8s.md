# k8s1.34 完整部署文档（Ubuntu 24.04 + 3节点高可用）

虚拟机部署，使用VMware安装Ubuntu 24.04 操作系统。

每个节点4核8G，可以根据本机配置修改，最低配置2核4G。

云服务器部署道理一样。

## 一、架构规划

| 角色        | IP地址        | 说明                               |
| ----------- | ------------- | ---------------------------------- |
| **deploy**  | 192.168.8.101 | Ansible执行机（复用为master/node） |
| **master1** | 192.168.8.101 | 控制平面 + etcd                    |
| **master2** | 192.168.8.102 | 控制平面 + etcd                    |
| **master3** | 192.168.8.103 | 控制平面 + etcd                    |
| **node1**   | 192.168.8.101 | 工作节点（复用）                   |
| **node2**   | 192.168.8.102 | 工作节点（复用）                   |
| **node3**   | 192.168.8.103 | 工作节点（复用）                   |

> **说明**：3节点环境下，每台机器既是Master又是Node。如果生产环境允许，建议单独准备一台 Deploy 机器，但3节点复用是资源受限的标准做法。

---

ubuntu网络配置文件：

~~~yaml
network:
  version: 2
  renderer: networkd
  ethernets:
    ens33:  # 修改为你的网卡名称
      dhcp4: no
      addresses:
        - 192.168.8.101/24
      routes:
        - to: default
          via: 192.168.8.2  # 网关地址，根据你的网络修改
      nameservers:
        addresses:
          - 223.5.5.5       # 阿里 DNS
          - 114.114.114.114 # 114 DNS
~~~

~~~shell
# 运行下面命令使网络生效 上面的局域网ip地址可以自行设定
# sudo netplan apply 
~~~

配置国内镜像（一键脚本）：

~~~shell
#!/bin/bash
# Ubuntu 24.04 (Noble) 清华大学镜像源一键配置脚本
# 适配新的 DEB822 格式（/etc/apt/sources.list.d/ubuntu.sources）
# 支持 x86_64, ARM64 架构

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 检查是否以 root 权限运行
if [ "$EUID" -ne 0 ]; then 
    echo -e "${RED}错误：请使用 sudo 或以 root 用户运行此脚本${NC}"
    exit 1
fi

# 检查系统版本
if ! command -v lsb_release &> /dev/null; then
    apt-get update && apt-get install -y lsb-release
fi

UBUNTU_VERSION=$(lsb_release -rs)
UBUNTU_CODENAME=$(lsb_release -cs)
ARCH=$(dpkg --print-architecture)

echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}Ubuntu 24.04 APT 镜像源配置工具${NC}"
echo -e "${BLUE}================================${NC}"
echo -e "系统版本: ${GREEN}Ubuntu ${UBUNTU_VERSION} (${UBUNTU_CODENAME})${NC}"
echo -e "系统架构: ${GREEN}${ARCH}${NC}"
echo -e "配置文件: ${GREEN}/etc/apt/sources.list.d/ubuntu.sources${NC}"
echo ""

# 版本检查
if [[ "$UBUNTU_CODENAME" != "noble" && "$UBUNTU_CODENAME" != "oracular" ]]; then
    echo -e "${YELLOW}警告：此脚本主要为 Ubuntu 24.04/24.10 (noble/oracular) 设计${NC}"
    echo -e "${YELLOW}当前检测到: ${UBUNTU_CODENAME}${NC}"
    read -p "是否继续? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# 备份原配置
SOURCES_FILE="/etc/apt/sources.list.d/ubuntu.sources"
BACKUP_FILE="${SOURCES_FILE}.bak.$(date +%Y%m%d_%H%M%S)"

echo -e "${YELLOW}>>> 备份原配置...${NC}"
if [ -f "$SOURCES_FILE" ]; then
    cp "$SOURCES_FILE" "$BACKUP_FILE"
    echo -e "${GREEN}✓ 已备份到: ${BACKUP_FILE}${NC}"
else
    echo -e "${YELLOW}! 未找到现有配置文件，将创建新文件${NC}"
    mkdir -p /etc/apt/sources.list.d
fi

# 写入清华镜像源配置（DEB822 格式）
echo -e "${YELLOW}>>> 写入清华大学镜像源（DEB822 格式）...${NC}"

cat > "$SOURCES_FILE" << EOF
# Ubuntu 24.04 LTS (Noble Numbat) - Tsinghua University Mirror
# 使用 DEB822 格式，Ubuntu 24.04+ 新格式
# 帮助文档: https://mirrors.tuna.tsinghua.edu.cn/help/ubuntu/

Types: deb
URIs: https://mirrors.tuna.tsinghua.edu.cn/ubuntu/
Suites: noble noble-updates noble-backports
Components: main restricted universe multiverse
Signed-By: /usr/share/keyrings/ubuntu-archive-keyring.gpg

# 安全更新（ Security ）
Types: deb
URIs: https://mirrors.tuna.tsinghua.edu.cn/ubuntu/
Suites: noble-security
Components: main restricted universe multiverse
Signed-By: /usr/share/keyrings/ubuntu-archive-keyring.gpg

# 源码镜像（默认注释，需要时取消注释）
# Types: deb-src
# URIs: https://mirrors.tuna.tsinghua.edu.cn/ubuntu/
# Suites: noble noble-updates noble-backports noble-security
# Components: main restricted universe multiverse
# Signed-By: /usr/share/keyrings/ubuntu-archive-keyring.gpg

# 预发布软件源（不建议启用）
# Types: deb
# URIs: https://mirrors.tuna.tsinghua.edu.cn/ubuntu/
# Suites: noble-proposed
# Components: main restricted universe multiverse
# Signed-By: /usr/share/keyrings/ubuntu-archive-keyring.gpg
EOF

# 清理旧的 sources.list（避免冲突）
if [ -f /etc/apt/sources.list ]; then
    echo -e "${YELLOW}>>> 清理旧的 sources.list...${NC}"
    mv /etc/apt/sources.list /etc/apt/sources.list.bak.disabled
    echo -e "${GREEN}✓ 已禁用旧的 sources.list${NC}"
fi

echo -e "${GREEN}✓ 镜像源配置已写入: ${SOURCES_FILE}${NC}"

# 更新 apt 缓存
echo -e "${YELLOW}>>> 更新软件包索引...${NC}"
if apt-get update; then
    echo -e "${GREEN}✓ APT 缓存更新成功${NC}"
else
    echo -e "${RED}✗ APT 更新失败，正在恢复备份...${NC}"
    if [ -f "$BACKUP_FILE" ]; then
        cp "$BACKUP_FILE" "$SOURCES_FILE"
        apt-get update
        echo -e "${YELLOW}已恢复原配置${NC}"
    fi
    exit 1
fi

# 验证配置
echo -e "${YELLOW}>>> 验证配置...${NC}"
if apt-cache policy | grep -q "tsinghua.edu.cn"; then
    echo -e "${GREEN}✓ 镜像源已成功切换至清华大学${NC}"
    echo -e "${BLUE}  当前使用的镜像: https://mirrors.tuna.tsinghua.edu.cn/ubuntu/${NC}"
else
    echo -e "${RED}✗ 镜像源配置可能未生效${NC}"
    exit 1
fi

# 显示统计信息
echo ""
echo -e "${BLUE}================================${NC}"
echo -e "${GREEN}配置完成！${NC}"
echo -e "${BLUE}================================${NC}"
echo -e "主镜像源: ${GREEN}清华大学 (TUNA)${NC}"
echo -e "配置文件: ${YELLOW}${SOURCES_FILE}${NC}"
echo -e "备份文件: ${YELLOW}${BACKUP_FILE}${NC}"
echo ""
echo -e "${YELLOW}常用命令:${NC}"
echo -e "  更新软件包: ${BLUE}sudo apt update && sudo apt upgrade${NC}"
echo -e "  安装软件:   ${BLUE}sudo apt install <package>${NC}"
echo -e "  查看配置:   ${BLUE}cat /etc/apt/sources.list.d/ubuntu.sources${NC}"
echo -e "  恢复备份:   ${BLUE}sudo cp ${BACKUP_FILE} ${SOURCES_FILE}${NC}"
echo ""

# 可选：清理旧内核等（询问用户）
read -p "是否立即更新系统软件包? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}>>> 正在更新系统...${NC}"
    apt-get upgrade -y
    echo -e "${GREEN}✓ 系统更新完成${NC}"
fi

exit 0
~~~

## 二、基础环境准备（所有节点执行）

### 1. 系统基础配置

```bash
# 1. 设置主机名（分别在对应机器执行）
sudo hostnamectl set-hostname master1  # 101
sudo hostnamectl set-hostname master2  # 102  
sudo hostnamectl set-hostname master3  # 103

# 2. 配置 hosts（所有节点执行）
cat <<EOF | sudo tee /etc/hosts
192.168.8.101 master1
192.168.8.102 master2
192.168.8.103 master3
EOF

# 3. 关闭 Swap（必须）
sudo swapoff -a
sudo sed -i '/swap/d' /etc/fstab

# 4. 配置时区和时间同步
sudo timedatectl set-timezone Asia/Shanghai
sudo apt update && sudo apt install -y chrony
sudo systemctl enable chrony --now

# 5. 安装必要工具
sudo apt install -y wget git vim ipset ipvsadm conntrack socat ebtables apt-transport-https ca-certificates curl gnupg lsb-release

# 6. 加载内核模块（Ubuntu 24.04 需要）
cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf
overlay
br_netfilter
ip_vs
ip_vs_rr
ip_vs_wrr
ip_vs_sh
nf_conntrack
EOF

sudo modprobe overlay
sudo modprobe br_netfilter
sudo modprobe ip_vs
sudo modprobe ip_vs_rr
sudo modprobe ip_vs_wrr
sudo modprobe ip_vs_sh
sudo modprobe nf_conntrack

# 7. 网络参数优化
cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-iptables  = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward                 = 1
net.ipv4.conf.all.forwarding        = 1
net.ipv4.conf.default.forwarding    = 1
net.ipv4.ip_local_reserved_ports    = 30000-32767
vm.swappiness                       = 0
vm.overcommit_memory                = 1
vm.panic_on_oom                     = 0
fs.inotify.max_user_watches         = 655360
fs.inotify.max_user_instances       = 8192
fs.file-max                         = 52706963
fs.nr_open                          = 52706963
net.ipv4.tcp_keepalive_time         = 600
net.ipv4.tcp_keepalive_intvl        = 30
net.ipv4.tcp_keepalive_probes       = 10
EOF

sudo sysctl --system
```

### 2. SSH免密配置

设置root账户+设置root ssh登录+免密登录

~~~shell
#!/bin/bash
echo ">>> 配置 root SSH 登录..."

# 1. 设置 root 密码（如果还没设置）
echo "请设置 root 密码："
sudo passwd root

# 2. 修改 SSH 配置
sudo tee -a /etc/ssh/sshd_config << 'EOF'

# 允许 root 登录（Kubeasz 部署需要）
PermitRootLogin yes
PasswordAuthentication yes
PubkeyAuthentication yes
EOF

# 3. 重启 SSH 服务（Ubuntu 24.04 服务名是 ssh）
sudo systemctl restart ssh

# 4. 验证状态
echo ">>> SSH 服务状态："
sudo systemctl status ssh --no-pager | head -5

echo ">>> 当前 root 登录配置："
grep -E "^PermitRootLogin|^PasswordAuthentication|^PubkeyAuthentication" /etc/ssh/sshd_config | tail -2

echo ">>> 完成！现在可以测试：ssh root@localhost"
~~~

以下 仅在 101 执行：

```bash
# 生成密钥
ssh-keygen -t rsa -b 4096 -N "" -f /root/.ssh/id_rsa

# 分发公钥到所有节点（包括自己）
for ip in 192.168.8.101 192.168.8.102 192.168.8.103; do
  ssh-copy-id -o StrictHostKeyChecking=no root@$ip
done
```

---

## 三、部署节点配置（192.168.8.101）

### 1. 下载 Kubeasz 3.6.8

```bash
# 创建安装目录
sudo mkdir -p /etc/kubeasz && cd /etc/kubeasz
# 下载 ezdown 工具（3.6.8版本）
export release=3.6.8
sudo wget https://github.com/easzlab/kubeasz/releases/download/${release}/ezdown
sudo chmod +x ./ezdown

# 查看版本信息
sudo ./ezdown --help
```

### 2. 下载离线包（包含 K8s v1.34.x 组件）

```bash
# 默认下载最新稳定版（K8s v1.34.x + containerd 2.x + etcd 3.6.x）
# 这里注意会遇到403问题，是清华站阻止了疑似的不合法请求，所以sudo vim ezdown，然后找到wget -c --user-agent=Mozilla --no-check-certificate 改为 wget -c即可
# 由于docker镜像需要翻墙才能下载，所以需要配置docker代理 注意代理软件打开局域网代理支持
# 创建 Docker 代理配置目录
sudo mkdir -p /etc/systemd/system/docker.service.d
# 创建代理配置文件
sudo tee /etc/systemd/system/docker.service.d/http-proxy.conf << 'EOF'
[Service]
Environment="HTTP_PROXY=http://192.168.8.1:7890"
Environment="HTTPS_PROXY=http://192.168.8.1:7890"
Environment="NO_PROXY=localhost,127.0.0.1,easzlab.io.local,.local"
EOF

# 重载配置并重启 Docker
sudo systemctl daemon-reload
sudo systemctl restart docker

# 验证代理生效
sudo systemctl show --property=Environment docker
# 执行命令
sudo ./ezdown -D

# 如需指定 K8s 版本（例如 v1.34.1）
# ./ezdown -D -k v1.34.1

# 国内环境如果下载慢，可使用镜像（自动识别）
# ./ezdown -D -m standard
```

**下载完成后检查**：

```bash
ls -la /etc/kubeasz/
# 应包含：bin/  down/  roles/  playbooks/  ezctl  ezdown 等
```

### 3. 初始化集群配置

```bash
# 创建新集群配置（命名为 k8s-01）
cd /etc/kubeasz/
sudo apt install -y ansible
sudo ./ezctl new k8s-01

# 配置文件位置
ls /etc/kubeasz/clusters/k8s-01/
# config.yml  hosts
```

---

## 四、修改集群配置

### 1. 配置主机清单（hosts）

编辑 `/etc/kubeasz/clusters/k8s-01/hosts`：

```ini
# 'deploy' node is not in k8s cluster, and only one
[deploy]
192.168.8.101 ansible_connection=ssh NTP_ENABLED=no

# 'etcd' cluster must have odd member(s): 1,3,5,7...
[etcd]
192.168.8.101 NODE_NAME=etcd1
192.168.8.102 NODE_NAME=etcd2
192.168.8.103 NODE_NAME=etcd3

[kube_master]
192.168.8.101
192.168.8.102
192.168.8.103

[kube_node]
192.168.8.101
192.168.8.102
192.168.8.103

# if set 'LB_MODE="haproxy"', the nodes above also serve as load balancer
[lb]
#192.168.8.101 LB_ROLE=backup
#192.168.8.102 LB_ROLE=master
#192.168.8.103 LB_ROLE=backup
[ex_lb]
# if set 'NEW_INSTALL=yes', then set group 'harbor' to install harbor server
[harbor]
#192.168.8.103 HARBOR_DOMAIN="harbor.easzlab.io.local" NEW_INSTALL=no
[chrony]

[all:vars]
# --------- Main Variables ---------------
# Secure port for apiservers
SECURE_PORT="6443"

# Cluster container-runtime supported: containerd, docker
CONTAINER_RUNTIME="containerd"

# Network plugins supported: calico, cilium, flannel, kube-router, kube-ovn
CLUSTER_NETWORK="calico"

# Service proxy mode: 'iptables' or 'ipvs'
PROXY_MODE="ipvs"

# K8s Service CIDR, not overlap with node(host) networking
SERVICE_CIDR="10.68.0.0/16"

# Cluster CIDR (Pod CIDR), not overlap with node(host) networking
CLUSTER_CIDR="172.20.0.0/16"

# NodePort Range
NODE_PORT_RANGE="30000-32767"

# Cluster DNS Domain
CLUSTER_DNS_DOMAIN="cluster.local"

# -------- Additional Variables (don't change the default value right now) ---
# Binaries Directory
bin_dir="/usr/local/bin"

# Deploy Directory (kubeasz workspace), don't change
base_dir="/etc/kubeasz"

# Directory for docker/certificates
ca_dir="/etc/kubernetes/ssl"

# Default nic name for lb, if lb is needed
# LB_IF="eth0"
```

### 2. 调整集群参数（config.yml）

编辑 `/etc/kubeasz/clusters/k8s-01/config.yml`，关键修改：

```yaml
# 新增
cluster_dir: "/etc/kubeasz/clusters/k8s-01"
k8s_nodename: ''
CALICO_DATASTORE: "kubernetes"
# 修改 Calico 模式，支持 BGP 或 IPIP（VMware环境推荐 IPIP）
CALICO_NETWORK_MODE: "IPIP"
```

---

## 五、执行部署

### 1. 一键部署（约10-15分钟）

```bash
cd /etc/kubeasz

# 语法检查（可选但推荐）
ansible-playbook -i clusters/k8s-01/hosts playbooks/90.setup.yml --syntax-check
# 切换到root用户
su - root
# 开始安装
./ezctl setup k8s-01 all
```

**部署流程说明**：

1. `00. 预安装检查` - 检查系统环境
2. `01. 安装容器运行时` - 安装 containerd
3. `02. 安装 etcd 集群` - 部署 3 节点 etcd
4. `03. 安装 Docker（可选）` - 跳过（使用 containerd）
5. `04. 安装 Master 节点` - 部署 kube-apiserver/kube-controller-manager/kube-scheduler
6. `05. 安装 Node 节点` - 部署 kubelet/kube-proxy
7. `06. 安装网络插件` - 部署 Calico
8. `07. 安装集群 DNS` - 部署 CoreDNS
9. `08. 安装 Metrics Server` - 部署监控组件
10. `09. 安装其他插件` - 根据配置安装 Ingress、Local Path Provisioner 等

### 2. 部署过程监控

如果部署中断，可单步执行：

```bash
# 单独执行某一步骤
./ezctl setup k8s-01 01  # 安装容器运行时
./ezctl setup k8s-01 04  # 安装 Master
./ezctl setup k8s-01 06  # 安装网络
```

---

## 六、验证集群

### 1. 检查节点状态

```bash
# 在任意 Master 节点执行
kubectl get nodes -o wide
```

**预期输出**：

```
NAME                STATUS   ROLES    AGE     VERSION   INTERNAL-IP     EXTERNAL-IP   OS-IMAGE             KERNEL-VERSION     CONTAINER-RUNTIME
k8s-192-168-8-101   Ready    master   2m33s   v1.34.1   192.168.8.101   <none>        Ubuntu 24.04.4 LTS   6.8.0-52-generic   containerd://2.1.4
k8s-192-168-8-102   Ready    master   2m34s   v1.34.1   192.168.8.102   <none>        Ubuntu 24.04.4 LTS   6.8.0-52-generic   containerd://2.1.4
k8s-192-168-8-103   Ready    master   2m33s   v1.34.1   192.168.8.103   <none>        Ubuntu 24.04.4 LTS   6.8.0-52-generic   containerd://2.1.4
```

### 2. 检查核心组件

```bash
# 检查 Pods 状态
kubectl get pods -n kube-system

# 检查 Etcd 健康状态
kubectl get componentstatuses

# 检查 Calico 状态
kubectl get pods -n kube-system -l k8s-app=calico-node
```

### 3. 测试部署应用

```bash
# 创建测试 Deployment
kubectl create deployment nginx-test --image=nginx:alpine --replicas=3

# 暴露服务
kubectl expose deployment nginx-test --port=80 --type=NodePort

# 查看服务
kubectl get svc nginx-test
# 访问：curl http://<任意节点IP>:<NodePort>
```

---

## 七、高可用优化（可选）

### 1. 备份 Etcd（重要）

```bash
# 手动备份
./ezctl backup k8s-01

# 查看备份
ls /etc/kubeasz/clusters/k8s-01/backup/
```

### 2. 开启自动补全

```bash
# 配置 kubectl 自动补全
echo 'source <(kubectl completion bash)' >> ~/.bashrc
source ~/.bashrc
```

---

## 八、运维命令参考

| 操作             | 命令                               |
| ---------------- | ---------------------------------- |
| **查看集群状态** | `kubectl cluster-info`             |
| **扩容 Node**    | `./ezctl add-node k8s-01 <新IP>`   |
| **扩容 Master**  | `./ezctl add-master k8s-01 <新IP>` |
| **删除 Node**    | `./ezctl del-node k8s-01 <IP>`     |
| **升级集群**     | `./ezctl upgrade k8s-01`           |
| **销毁集群**     | `./ezctl destroy k8s-01`           |

---

## 九、常见问题处理（Ubuntu 24.04 特供）

### 1. 内核模块加载失败

```bash
# 如果 modprobe ip_vs 失败，安装额外模块
sudo apt install -y linux-modules-extra-$(uname -r)
sudo reboot
```

### 2. Containerd 服务启动失败

```bash
# 检查配置
sudo containerd config default > /tmp/containerd.toml
sudo vim /etc/containerd/config.toml  # 确保 SystemdCgroup = true

# 重启
sudo systemctl restart containerd
```

### 3. Calico 节点无法启动

VMware 环境下可能需要调整：

```bash
# 在 config.yml 中确保
CALICO_IPV4POOL_IPIP: "Always"
CALICO_IPV4POOL_VXLAN: "Never"
```

## 4. 如果重启后报错

如果重启后报如下错误：

~~~shell
huo@k8s-192-168-8-101:~$ kubectl get nodes
E0221 23:49:07.660992    8532 memcache.go:265] "Unhandled Error" err="couldn't get current server API group list: Get \"http://localhost:8080/api?timeout=32s\": dial tcp [::1]:8080: connect: connection refused"
E0221 23:49:07.662235    8532 memcache.go:265] "Unhandled Error" err="couldn't get current server API group list: Get \"http://localhost:8080/api?timeout=32s\": dial tcp [::1]:8080: connect: connection refused"
E0221 23:49:07.663528    8532 memcache.go:265] "Unhandled Error" err="couldn't get current server API group list: Get \"http://localhost:8080/api?timeout=32s\": dial tcp [::1]:8080: connect: connection refused"
E0221 23:49:07.664961    8532 memcache.go:265] "Unhandled Error" err="couldn't get current server API group list: Get \"http://localhost:8080/api?timeout=32s\": dial tcp [::1]:8080: connect: connection refused"
E0221 23:49:07.666058    8532 memcache.go:265] "Unhandled Error" err="couldn't get current server API group list: Get \"http://localhost:8080/api?timeout=32s\": dial tcp [::1]:8080: connect: connection refused"
The connection to the server localhost:8080 was refused - did you specify the right host or port?

~~~

解决方案，设置环境变量：

~~~shell
echo 'export KUBECONFIG=/etc/kubeasz/clusters/k8s-01/kubectl.kubeconfig' >> ~/.bashrc
source ~/.bashrc
~~~

