# 云端智师 PVE 管理系统 - 后端 API 文档

> **版本**: v1.0.0
> **基础 URL**: `http://localhost:8080/api/v1`
> **WebSocket URL**: `ws://localhost:8080`
> **数据格式**: JSON
> **字符编码**: UTF-8

---

## 状态码定义

| 状态码 | 说明 | 处理方式 |
|--------|------|---------|
| 200 | 请求成功 | 正常处理返回数据 |
| 201 | 创建成功 | 正常处理返回数据 |
| 400 | 请求参数错误 | 检查请求参数格式 |
| 401 | 未认证 | 跳转到登录页 |
| 403 | 权限不足 | 提示权限不足 |
| 404 | 资源不存在 | 提示资源不存在 |
| 500 | 服务器内部错误 | 提示稍后重试 |

---

## 认证与授权

### 登录

```
POST /api/v1/auth/login
```

**请求体**:
```json
{
  "username": "admin",
  "password": "your_password",
  "realm": "pve",
  "otp": "123456",
  "language": "zh-CN"
}
```

**参数说明**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |
| realm | string | 是 | 认证域（pve/pam/ldap/ad） |
| otp | string | 否 | 双因子验证码 |
| language | string | 否 | 界面语言（默认 zh-CN） |

**返回**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "admin",
    "username": "admin",
    "role": "admin",
    "real_name": "系统管理员",
    "email": "admin@school.local"
  },
  "expires_at": "2026-05-23T10:00:00Z"
}
```

### 登出

```
POST /api/v1/auth/logout
```

**请求头**: `Authorization: Bearer <token>`

**返回**:
```json
{ "status": "logged_out" }
```

---

## 数据中心 API

### 获取集群概要

```
GET /api/v1/datacenter/summary
```

**返回**:
```json
{
  "cluster_name": "cloud-smart-tutor",
  "version": "8.2.4",
  "nodes": 3,
  "vms": 52,
  "containers": 12,
  "storage_pools": 4,
  "cpu_usage_percent": 45.2,
  "memory_usage_percent": 62.8,
  "disk_usage_percent": 38.5,
  "ksm_saving_percent": 58.3,
  "health_status": "healthy"
}
```

### 获取集群节点列表

```
GET /api/v1/datacenter/nodes
```

**返回**:
```json
{
  "nodes": [
    {
      "id": "node-01",
      "name": "pve-node-01",
      "status": "online",
      "cpu_usage": 45.2,
      "memory_total_gb": 64,
      "memory_used_gb": 42.3,
      "disk_total_gb": 2000,
      "disk_used_gb": 780,
      "uptime_hours": 720,
      "version": "8.2.4",
      "ip_address": "192.168.1.1"
    }
  ]
}
```

### 搜索资源

```
GET /api/v1/datacenter/search?q=keyword
```

**参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| q | string | 是 | 搜索关键词 |

**返回**:
```json
{
  "results": [
    {
      "id": "vm-100",
      "name": "student-01",
      "type": "vm",
      "node": "pve-node-01",
      "status": "running"
    }
  ]
}
```

---

## 虚拟机 API

### 获取虚拟机列表

```
GET /api/v1/vms?node=node-01&type=vm
```

**参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| node | string | 否 | 节点 ID |
| type | string | 否 | 类型（vm/ct） |
| status | string | 否 | 状态过滤 |
| page | int | 否 | 页码（默认 1） |
| per_page | int | 否 | 每页数量（默认 20） |

**返回**:
```json
{
  "vms": [
    {
      "vmid": 100,
      "name": "student-01",
      "type": "vm",
      "status": "running",
      "node": "pve-node-01",
      "cpu_usage": 45.2,
      "memory_mb": 2048,
      "disk_gb": 30,
      "ip_address": "192.168.1.101",
      "uptime_seconds": 86400,
      "tags": ["classroom-01", "linux-course"]
    }
  ],
  "total": 52,
  "page": 1,
  "per_page": 20
}
```

### 创建虚拟机

```
POST /api/v1/vms
```

**请求体**:
```json
{
  "node": "pve-node-01",
  "vmid": 100,
  "name": "student-01",
  "type": "vm",
  "os_image": "debian-12.0_x64",
  "cpu_cores": 2,
  "memory_mb": 2048,
  "disk_gb": 30,
  "network": {
    "bridge": "vmbr0",
    "model": "virtio"
  },
  "ai_enabled": true,
  "ai_model": "ollama-llama3",
  "tags": ["classroom-01"]
}
```

**返回**:
```json
{
  "task_id": "UPID:pve-node-01:00012345:00012346:663E1234:vzcreate:100:admin@pve:",
  "vmid": 100,
  "status": "creating"
}
```

### 获取虚拟机详情

```
GET /api/v1/vms/{vmid}
```

**返回**:
```json
{
  "vmid": 100,
  "name": "student-01",
  "type": "vm",
  "status": "running",
  "node": "pve-node-01",
  "description": "学生虚拟机",
  "cpu": {
    "cores": 2,
    "sockets": 1,
    "usage_percent": 45.2
  },
  "memory": {
    "total_mb": 2048,
    "used_mb": 1280,
    "usage_percent": 62.5,
    "balloon_mb": 2048
  },
  "disks": [
    {
      "name": "scsi0",
      "size_gb": 30,
      "used_gb": 12.5,
      "type": "qcow2",
      "storage": "local-lvm"
    }
  ],
  "networks": [
    {
      "name": "net0",
      "model": "virtio",
      "bridge": "vmbr0",
      "mac_address": "32:33:34:35:36:37",
      "ip_address": "192.168.1.101"
    }
  ],
  "os": {
    "image": "debian-12.0_x64",
    "type": "linux"
  },
  "created_at": "2026-05-15T09:00:00Z",
  "uptime_seconds": 86400,
  "tags": ["classroom-01", "linux-course"]
}
```

### 启动虚拟机

```
POST /api/v1/vms/{vmid}/start
```

**返回**:
```json
{
  "task_id": "UPID:pve-node-01:00012347:00012348:663E1235:qmstart:100:admin@pve:",
  "status": "starting"
}
```

### 停止虚拟机

```
POST /api/v1/vms/{vmid}/stop
```

**返回**:
```json
{
  "task_id": "UPID:pve-node-01:00012349:0001234A:663E1236:qmstop:100:admin@pve:",
  "status": "stopping"
}
```

### 重启虚拟机

```
POST /api/v1/vms/{vmid}/reboot
```

### 暂停虚拟机

```
POST /api/v1/vms/{vmid}/suspend
```

### 恢复虚拟机

```
POST /api/v1/vms/{vmid}/resume
```

### 删除虚拟机

```
DELETE /api/v1/vms/{vmid}
```

### 克隆虚拟机

```
POST /api/v1/vms/{vmid}/clone
```

**请求体**:
```json
{
  "newid": 200,
  "name": "student-01-clone",
  "full": false,
  "node": "pve-node-01"
}
```

### 迁移虚拟机

```
POST /api/v1/vms/{vmid}/migrate
```

**请求体**:
```json
{
  "target_node": "pve-node-02",
  "online": true
}
```

### 获取 noVNC 控制台 URL

```
GET /api/v1/vms/{vmid}/console/novnc
```

**返回**:
```json
{
  "url": "wss://192.168.1.1:8006/api2/json/nodes/pve-node-01/qemu/100/vncwebsocket",
  "ticket": "PVEVNC:...",
  "port": 5900,
  "password": "xyz123"
}
```

### 获取 SPICE 控制台 URL

```
GET /api/v1/vms/{vmid}/console/spice
```

### 获取虚拟机选项配置

```
GET /api/v1/vms/{vmid}/options
```

### 更新虚拟机选项

```
PUT /api/v1/vms/{vmid}/options
```

### 获取虚拟机快照列表

```
GET /api/v1/vms/{vmid}/snapshots
```

### 创建快照

```
POST /api/v1/vms/{vmid}/snapshots
```

**请求体**:
```json
{
  "snapname": "pre-update",
  "description": "更新前快照"
}
```

### 回滚快照

```
POST /api/v1/vms/{vmid}/snapshots/{snapname}/rollback
```

### 删除快照

```
DELETE /api/v1/vms/{vmid}/snapshots/{snapname}
```

---

## 节点 API

### 获取节点概要

```
GET /api/v1/nodes/{nodeid}/summary
```

**返回**:
```json
{
  "node": "pve-node-01",
  "status": "online",
  "cpu": {
    "cores": 16,
    "sockets": 1,
    "usage_percent": 45.2,
    "model": "Intel Xeon E5-2680 v4"
  },
  "memory": {
    "total_gb": 64,
    "used_gb": 42.3,
    "free_gb": 21.7,
    "usage_percent": 62.8
  },
  "disk": {
    "total_gb": 2000,
    "used_gb": 780,
    "free_gb": 1220,
    "usage_percent": 38.5
  },
  "uptime_hours": 720,
  "version": "8.2.4",
  "ip_address": "192.168.1.1"
}
```

### 获取节点网络配置

```
GET /api/v1/nodes/{nodeid}/network
```

### 获取节点磁盘信息

```
GET /api/v1/nodes/{nodeid}/disks
```

### 获取节点系统日志

```
GET /api/v1/nodes/{nodeid}/syslog
```

### 获取节点任务记录

```
GET /api/v1/nodes/{nodeid}/tasks
```

---

## 存储 API

### 获取存储列表

```
GET /api/v1/storage
```

### 获取存储详情

```
GET /api/v1/storage/{storageid}
```

### 创建存储

```
POST /api/v1/storage
```

### 获取存储内容

```
GET /api/v1/storage/{storageid}/content
```

---

## 任务 API

### 获取当前运行任务

```
GET /api/v1/tasks/running
```

### 获取任务日志

```
GET /api/v1/tasks/{taskid}/log
```

### 停止任务

```
POST /api/v1/tasks/{taskid}/stop
```

---

## 权限 API

### 获取用户列表

```
GET /api/v1/users
```

### 创建用户

```
POST /api/v1/users
```

### 获取用户组列表

```
GET /api/v1/groups
```

### 获取角色列表

```
GET /api/v1/roles
```

---

## 防火墙 API

### 获取集群防火墙规则

```
GET /api/v1/firewall/rules
```

### 创建防火墙规则

```
POST /api/v1/firewall/rules
```

---

## 实时推送（WebSocket）

### 连接 WebSocket

```
GET /ws/{type}/{identifier}
```

**类型**:
| 类型 | 标识符 | 说明 |
|------|--------|------|
| idc | - | IDC 主控大屏 |
| teacher | classroom_id | 教师控制台 |
| student | student_id | 学生控制台 |
| enterprise | classroom_id | 企业沙箱 |

### 消息类型

| 类型 | 说明 | 触发条件 |
|------|------|---------|
| node_metrics | 节点指标 | 每 5 秒推送 |
| vm_status_change | VM 状态变更 | 启动/停止/迁移时 |
| task_progress | 任务进度 | 任务执行中 |
| error_cluster | 错误聚类 | 探针上报错误时 |
| teaching_alert | 教学预警 | 错误率超阈值时 |
| ai_companion | AI 伴学响应 | 大模型返回时 |

---

## 前端 Mock 数据接入指南

### 1. 配置 Mock API 服务

在前端项目中创建 `src/services/mock-api.js`，所有 API 调用通过统一的 `api` 模块进行：

```javascript
// src/services/api.js
const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:8080/api/v1'
const USE_MOCK = import.meta.env.VITE_USE_MOCK === 'true'

import { mockApi } from './mock-api'

const api = {
  get: async (path, params = {}) => {
    if (USE_MOCK) return mockApi.get(path, params)
    const res = await fetch(`${API_BASE}${path}?${new URLSearchParams(params)}`)
    return res.json()
  },
  post: async (path, data = {}) => {
    if (USE_MOCK) return mockApi.post(path, data)
    const res = await fetch(`${API_BASE}${path}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data)
    })
    return res.json()
  },
  // ... 其他方法
}

export default api
```

### 2. 环境变量控制

```env
# .env.development
VITE_API_BASE=http://localhost:8080/api/v1
VITE_USE_MOCK=true  # 切换为 false 即可接入真实后端
```

### 3. 联调机制

1. **前端使用 Mock 数据开发**：所有功能先基于 Mock 数据开发完成
2. **后端 API 就绪后**：将 `VITE_USE_MOCK=false` 即可切换至真实后端
3. **数据格式一致**：Mock 数据格式与真实 API 返回格式完全一致
4. **定期同步**：每周进行一次前后端联调，解决接口差异问题
