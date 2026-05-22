# 云端智师 - 前后端联调机制与开发协作指南

> **版本**: v1.0.0  
> **更新日期**: 2026-05-22

---

## 一、前后端开发协作流程

### 1. 开发阶段划分

| 阶段 | 前端任务 | 后端任务 | 联调重点 |
|------|---------|---------|---------|
| **第一阶段（1-2 周）** | 基于 Mock 数据完成所有页面开发 | API 基础框架搭建 + 认证系统 | 登录认证、Token 传递 |
| **第二阶段（3-4 周）** | 接入真实后端 API，逐步替换 Mock | 实现核心 VM 管理 API | VM 生命周期管理、状态同步 |
| **第三阶段（5-6 周）** | WebSocket 实时推送集成 | WebSocket Hub + 探针接入 | 实时数据推送、错误聚类 |
| **第四阶段（7-8 周）** | 全功能联调测试 + 优化 | PVE 完整对接 + 性能优化 | 端到端测试、性能调优 |

### 2. 联调机制

**每周三下午进行前后端联调会议**，内容包括：
1. 前端演示本周开发的功能（基于 Mock 或真实 API）
2. 后端演示已完成的 API 接口
3. 对照 API 文档逐项验证接口返回格式
4. 记录接口差异，当天修复

### 3. 接口差异处理

| 差异类型 | 处理方式 | 责任人 |
|---------|---------|--------|
| 返回字段缺失 | 前端添加容错处理，后端补充字段 | 后端开发 |
| 字段类型不一致 | 统一为文档定义类型 | 后端开发 |
| 状态码不一致 | 按文档统一状态码 | 后端开发 |
| 新增字段 | 更新文档，前端适配 | 双方协商 |

---

## 二、前端 Mock 数据方案

### 1. 启用/关闭 Mock

```env
# .env.development
VITE_USE_MOCK=true   # 使用 Mock 数据
VITE_USE_MOCK=false  # 使用真实后端 API
```

### 2. Mock 数据位置

- `src/services/mock-api.js` - 所有 Mock 数据与逻辑
- `src/services/api.js` - 统一 API 调用层（自动切换 Mock/真实）

### 3. Mock 数据格式保证

**Mock 数据格式与真实 API 完全一致**，包括：
- 字段名称
- 数据类型
- 嵌套结构
- 状态码含义

切换至真实后端时，**无需修改任何前端代码**。

---

## 三、后端 API 对接清单

### 已完成 API（可立即对接）

| API | 状态 | 前端使用方式 |
|-----|------|-------------|
| `POST /api/v1/auth/login` | ✅ 可用 | `api.auth.login()` |
| `POST /api/v1/auth/logout` | ✅ 可用 | `api.auth.logout()` |
| `GET /api/v1/datacenter/summary` | ✅ 可用 | `api.datacenter.getSummary()` |
| `GET /api/v1/datacenter/nodes` | ✅ 可用 | `api.datacenter.getNodes()` |
| `GET /api/v1/datacenter/search?q=` | ✅ 可用 | `api.datacenter.search()` |
| `GET /api/v1/vms` | ✅ 可用 | `api.vms.getList()` |
| `GET /api/v1/vms/{vmid}` | ✅ 可用 | `api.vms.getDetail()` |
| `POST /api/v1/vms` | ✅ 可用 | `api.vms.create()` |
| `POST /api/v1/vms/{vmid}/start` | ✅ 可用 | `api.vms.start()` |
| `POST /api/v1/vms/{vmid}/stop` | ✅ 可用 | `api.vms.stop()` |
| `POST /api/v1/vms/{vmid}/reboot` | ✅ 可用 | `api.vms.reboot()` |
| `DELETE /api/v1/vms/{vmid}` | ✅ 可用 | `api.vms.delete()` |
| `GET /api/v1/vms/{vmid}/console/novnc` | ✅ 可用 | `api.vms.getConsole()` |
| `GET /api/v1/vm/os-images` | ✅ 可用 | `api.vms.getOSImages()` |
| `GET /api/v1/vm/ai-models` | ✅ 可用 | `api.vms.getAIModels()` |
| `GET /api/v1/tasks/running` | ✅ 可用 | `api.tasks.getRunning()` |
| `GET /api/v1/tasks/{taskId}/log` | ✅ 可用 | `api.tasks.getLog()` |
| `POST /api/v1/tasks/{taskId}/stop` | ✅ 可用 | `api.tasks.stop()` |
| `GET /api/v1/storage` | ✅ 可用 | `api.storage.getList()` |
| `GET /api/v1/users` | ✅ 可用 | `api.permissions.getUsers()` |

### 待开发 API（当前使用 Mock 数据）

| API | 预计完成 | 前端影响 |
|-----|---------|---------|
| `POST /api/v1/vms/{vmid}/clone` | 第 3 周 | 克隆功能 |
| `POST /api/v1/vms/{vmid}/migrate` | 第 3 周 | 迁移功能 |
| `GET /api/v1/nodes/{nodeId}/summary` | 第 2 周 | 节点详情 |
| `GET /api/v1/nodes/{nodeId}/disks` | 第 2 周 | 磁盘信息 |
| `POST /api/v1/vms/{vmid}/snapshots` | 第 4 周 | 快照管理 |
| `GET /api/v1/firewall/rules` | 第 3 周 | 防火墙管理 |

---

## 四、前端开发规范

### 1. 设计风格

- **Apple 设计语言**（参考 `awesome-design-md-main/design-md/apple/DESIGN.md`）
- **Proxmox VE 布局**（四区域：标题栏/资源树/内容面板/日志面板）

### 2. 组件命名

- 页面组件：`PascalCase.vue`（如 `TeacherConsole.vue`）
- 通用组件：`AppleCard.vue`、`AppleButton.vue`
- 服务层：`camelCase.js`（如 `api.js`、`mock-api.js`）

### 3. API 调用规范

```javascript
// ✅ 正确：使用统一 api 模块
import { api } from '@/services/api'

const loadVMs = async () => {
  const res = await api.vms.getList()
  vms.value = res.vms
}

// ❌ 错误：直接调用 fetch
const res = await fetch('/api/v1/vms')
```

### 4. 错误处理

```javascript
try {
  await api.vms.start(vmid)
} catch (err) {
  if (err.message.includes('401')) {
    // 未认证，跳转登录
    router.push('/login')
  } else {
    // 其他错误，提示用户
    alert('操作失败：' + err.message)
  }
}
```

---

## 五、WebSocket 实时推送对接

### 1. 连接方式

```javascript
const connectWebSocket = (type, identifier) => {
  const ws = new WebSocket(`${import.meta.env.VITE_WS_BASE}/ws/${type}/${identifier}`)

  ws.onmessage = (event) => {
    const data = JSON.parse(event.data)
    handleWebSocketMessage(data)
  }

  ws.onclose = () => {
    // 自动重连
    setTimeout(() => connectWebSocket(type, identifier), 3000)
  }

  return ws
}

// 教师端
connectWebSocket('teacher', 'classroom-01')

// 学生端
connectWebSocket('student', 'stu-001')

// IDC 大屏
connectWebSocket('idc', '')
```

### 2. 消息处理

```javascript
const handleWebSocketMessage = (data) => {
  switch (data.event) {
    case 'node_metrics':
      updateNodeMetrics(data.payload)
      break
    case 'vm_status_change':
      updateVMStatus(data.payload)
      break
    case 'task_progress':
      updateTaskProgress(data.payload)
      break
    case 'error_cluster':
      showErrorCluster(data.payload)
      break
    case 'teaching_alert':
      showTeachingAlert(data.payload)
      break
    case 'ai_companion':
      showAIResponse(data.payload)
      break
  }
}
```

---

## 六、教师端功能需求清单

### 已实现功能

| 功能 | 状态 | 说明 |
|------|------|------|
| 登录认证 | ✅ 完成 | 多 Realm 认证、语言选择、保存用户名 |
| 标题栏 | ✅ 完成 | Logo、搜索、用户设置、创建 VM/CT、退出 |
| 资源树 | ✅ 完成 | 数据中心/节点/VM/存储、视图切换 |
| 内容面板 | ✅ 完成 | 概要/控制台/硬件/选项/快照/防火墙/权限 |
| 日志面板 | ✅ 完成 | 任务日志/集群日志、双击查看详情、可缩放隐藏 |
| 虚拟机操作 | ✅ 完成 | 启动/关机/停止/重置/迁移/控制台 |
| 设置对话框 | ✅ 完成 | 存储管理、xterm.js 设置、布局重置 |
| 创建虚拟机 | ✅ 完成 | 节点选择、名称、镜像、CPU/内存/磁盘 |

### 待完善功能

| 功能 | 优先级 | 预计完成 |
|------|--------|---------|
| SPICE 安全终端 | P1 | 第 3 周 |
| 双因子认证 (OATH/YubiKey) | P1 | 第 2 周 |
| 批量操作 (Bulk Start/Stop/Migrate) | P2 | 第 3 周 |
| 备份管理 | P2 | 第 4 周 |
| 复制任务管理 | P2 | 第 4 周 |
| HA 高可用管理 | P3 | 第 5 周 |
| Ceph 集群管理 | P3 | 第 5 周 |
| 订阅管理 | P3 | 第 5 周 |

---

## 七、快速开始

### 前端开发

```bash
# 安装依赖
npm install

# 启动开发服务器（使用 Mock 数据）
npm run dev

# 访问教师控制台
# http://localhost:5175/teacher-console
```

### 后端启动

```bash
cd backend
go mod tidy
go run cmd/server/main.go

# API 文档
# http://localhost:8080/api-docs
```

### 切换至真实后端

```env
# .env.development
VITE_USE_MOCK=false
```

---

## 八、联系人

| 角色 | 负责人 | 联系方式 |
|------|--------|---------|
| 后端开发 | AI Assistant | - |
| 前端开发 | AI Assistant | - |
| 项目经理 | 待定 | - |
| UI/UX 设计 | Apple Design System | `awesome-design-md-main/design-md/apple/` |
