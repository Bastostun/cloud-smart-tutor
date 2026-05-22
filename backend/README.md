# 云端智师 - 微型 IDC 高并发 Go 后端

## 架构概述

本后端是"软件定义算力、AI 伴随审计"的微型 IDC 控制中枢，不是传统的 CRUD 信息系统。Go 后端同时：
- **向下**调度 Proxmox VE 虚拟化集群硬件资源
- **向内**对虚拟机探针上传的报错进行 Map 计数与动态聚类统计
- **向外**对接大模型 API 进行 RAG 编排
- **通过 WebSocket 全双工管道**实时反哺给前端 Vue 3 多端控制台

## 项目结构

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # 应用入口，启动所有子系统
├── internal/
│   ├── config/
│   │   └── config.go            # 环境变量配置加载
│   ├── domain/
│   │   ├── models.go            # 核心数据结构（EdgeNode、CloneRequest 等）
│   │   └── events.go            # 事件定义（ProbeEvent、AuditLog 等）
│   ├── websocket/
│   │   ├── hub.go               # 高并发 WebSocket 连接池（房间广播）
│   │   └── client.go            # 客户端连接管理（读写分离 Pump）
│   ├── pve/
│   │   └── client.go            # PVE API 客户端 + 并发调度器
│   ├── cluster/
│   │   └── engine.go            # 内存错误聚类引擎（动态阈值告警）
│   ├── sandbox/
│   │   └── audit.go             # 零信任沙箱审计流（秒级拦截）
│   ├── probe/
│   │   └── ingestor.go          # 探针数据接收 + AI 伴随式响应管道
│   ├── monitor/
│   │   └── registry.go          # 节点注册表 + 心跳监测 + 指标广播
│   └── handlers/
│       └── handlers.go          # Gin 路由 + WebSocket 端点
├── go.mod
└── go.sum
```

## 核心子系统

### 1. WebSocket Hub 高并发连接池
- **房间（Room）广播机制**：按教室/角色/学生隔离消息流
- **读写分离 Pump**：每个客户端独立的 ReadPump/WritePump
- **Ping/Pong 心跳**：30 秒保活，60 秒超时断开
- **非阻塞发送**：Send channel buffer full 时自动丢弃，不阻塞 Hub

### 2. PVE 并发调度器
- **Goroutine 工作池**：可配置 worker 数量（默认 10）
- **Channel 任务分发**：CloneTask 通过 Channel 异步分发
- **批量克隆**：50 台沙箱 10 秒内批量拉起，不阻塞主线程
- **结果收集**：通过 Channel 收集克隆结果，实时更新批次状态

### 3. 内存错误聚类引擎
- **FNV-32a 哈希聚类**：错误文本归一化后哈希，快速聚类
- **sync.RWMutex 线程安全**：读写锁保护内存 Map
- **动态阈值触发**：默认 20% 学生触发教学预警
- **5 分钟冷却**：同一错误 5 分钟内不重复告警
- **AI 教学建议生成**：内置 Nginx、端口、权限等常见错误分析

### 4. 零信任沙箱审计流
- **6 条默认安全规则**：剪贴板、网络外发、键盘审计、DLP、截屏、USB
- **4 级干预动作**：none / warn / block / suspend_vm
- **风险评分**：动态计算 VM 风险分数（0-100）
- **毫秒级推送**：安全事件通过 WebSocket 实时推送

### 5. AI 伴随式响应管道
- **自动错误分析**：6 类常见错误内置解析逻辑
- **解决方案生成**：分步骤的排查指南
- **置信度评分**：0.60-0.95 置信度
- **主动触发**：探针上报后自动触发，无需学生请求

## API 端点

### REST API
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/probe/event` | 接收探针事件 |
| POST | `/api/v1/heartbeat` | 节点心跳 |
| POST | `/api/v1/nodes/register` | 注册节点 |
| GET | `/api/v1/nodes` | 获取所有节点 |
| GET | `/api/v1/nodes/:id` | 获取单个节点 |
| GET | `/api/v1/classroom/:id/stats` | 教室统计 |
| GET | `/api/v1/classroom/:id/errors` | 教室错误排行 |
| POST | `/api/v1/clone` | 批量克隆 VM |
| GET | `/api/v1/clone/:batchId/status` | 克隆批次状态 |
| GET | `/api/v1/sandbox/policies` | 沙箱策略 |
| GET | `/api/v1/health` | 健康检查 |

### WebSocket 端点
| Path | Description |
|------|-------------|
| `/ws/idc` | IDC 主控大屏 |
| `/ws/teacher/:classroomId` | 教师控制台 |
| `/ws/student/:studentId` | 学生伴学舱 |
| `/ws/enterprise/:classroomId` | 企业沙箱 |
| `/ws/probe` | 探针数据上报 |

## WebSocket 消息类型
| Type | Description |
|------|-------------|
| `node_metrics` | 节点指标数据 |
| `probe_event` | 探针事件 |
| `error_cluster` | 错误聚类更新 |
| `teaching_alert` | 教学预警 |
| `ai_companion` | AI 伴学响应 |
| `audit_log` | 审计日志 |
| `clone_status` | 克隆状态 |
| `heartbeat` | 心跳 |

## 环境变量

```bash
# 服务器配置
SERVER_ADDRESS=0.0.0.0
SERVER_PORT=8080
GIN_MODE=release

# PVE 集群配置
PVE_ENDPOINT=https://pve-cluster.local:8006
PVE_USERNAME=root@pam
PVE_PASSWORD=your_password
PVE_REALM=pam
PVE_TLS_SKIP_VERIFY=true
PVE_TIMEOUT_SECONDS=30

# 并发配置
MAX_CONCURRENT_CLONES=50
PVE_WORKERS=10

# 监测配置
NODE_POLL_INTERVAL=5           # 分钟
ERROR_CLEANUP_MAX_AGE=30       # 分钟
ERROR_ALERT_THRESHOLD=20.0     # 百分比

# 教室配置
DEFAULT_CLASSROOM_ID=classroom-01
DEFAULT_TOTAL_STUDENTS=52
```

## 构建与运行

```bash
# 下载依赖
go mod tidy

# 运行（开发模式）
GIN_MODE=debug go run cmd/server/main.go

# 构建
go build -o cloud-smart-tutor ./cmd/server/

# 运行编译后的二进制文件
./cloud-smart-tutor
```

## 前端对接

Vue 3 前端通过 WebSocket 连接后端：

```javascript
// IDC 大屏
const ws = new WebSocket('ws://localhost:8080/ws/idc')

// 教师控制台
const ws = new WebSocket('ws://localhost:8080/ws/teacher/classroom-01')

// 学生伴学舱
const ws = new WebSocket('ws://localhost:8080/ws/student/student-001')

// 企业沙箱
const ws = new WebSocket('ws://localhost:8080/ws/enterprise/classroom-01')
```

## 线程安全设计

- `websocket.Hub`：`sync.RWMutex` 保护 Clients 和 Rooms Map
- `cluster.ErrorClusteringEngine`：`sync.RWMutex` 保护 clusters Map，独立 `sync.Mutex` 保护 alertCooldown
- `pve.PVEScheduler`：`sync.RWMutex` 保护 Results Map
- `sandbox.AuditStream`：`sync.RWMutex` 保护 VMStates Map
- `monitor.NodeRegistry`：`sync.RWMutex` 保护 Nodes 和 Classrooms Map

所有并发访问均通过互斥锁保护，避免竞态条件。
