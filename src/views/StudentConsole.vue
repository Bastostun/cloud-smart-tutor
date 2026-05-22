<template>
  <div class="min-h-screen bg-apple-canvas">
    <!-- Apple 子导航栏 -->
    <div class="apple-sub-nav fixed top-[44px] left-0 right-0 z-40">
      <div class="max-w-[1440px] mx-auto px-4 h-full flex items-center justify-between">
        <div class="flex items-center space-x-4">
          <h2 class="apple-tagline">虚拟机控制台</h2>
          <span class="apple-caption text-apple-ink-muted-48">·</span>
          <span class="apple-caption text-apple-ink-muted-48">{{ studentInfo.name }} · {{ studentInfo.classroom }}</span>
        </div>
        <div class="flex items-center space-x-3">
          <div class="apple-badge apple-badge-success">
            <div class="w-1.5 h-1.5 rounded-full bg-current"></div>
            <span>虚拟机运行中</span>
          </div>
          <button @click="logout" class="apple-btn-utility">退出</button>
        </div>
      </div>
    </div>

    <!-- 主内容区 -->
    <main class="pt-[96px] pb-[44px] min-h-screen">
      <div class="max-w-[1440px] mx-auto px-4">
        <!-- 加载状态 -->
        <div v-if="isLoading" class="flex items-center justify-center h-64">
          <div class="text-center">
            <svg class="animate-spin w-8 h-8 mx-auto mb-4 text-apple-blue" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <p class="apple-caption text-apple-ink-muted-48">加载虚拟机信息中...</p>
          </div>
        </div>
        
        <!-- 标签导航 -->
        <div v-else class="flex items-center space-x-apple-sm mb-apple-lg border-b border-apple-divider">
          <button v-for="tab in tabs" :key="tab.id"
            @click="activeTab = tab.id"
            class="apple-tab"
            :class="{ 'apple-tab-active': activeTab === tab.id }">
            {{ tab.label }}
          </button>
        </div>

        <!-- 概览模块 -->
        <div v-if="activeTab === 'overview'" class="space-y-apple-lg">
          <!-- 用户信息卡片 -->
          <div class="apple-card">
            <h2 class="apple-tagline mb-apple-lg">用户信息</h2>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-apple-lg">
              <div>
                <div class="apple-caption text-apple-ink-muted-48">用户名</div>
                <div class="apple-body-strong mt-apple-xxs">{{ studentInfo.name }}</div>
              </div>
              <div>
                <div class="apple-caption text-apple-ink-muted-48">班级</div>
                <div class="apple-body-strong mt-apple-xxs">{{ studentInfo.classroom }}</div>
              </div>
              <div>
                <div class="apple-caption text-apple-ink-muted-48">学号</div>
                <div class="apple-body-strong mt-apple-xxs">{{ studentInfo.id }}</div>
              </div>
              <div>
                <div class="apple-caption text-apple-ink-muted-48">状态</div>
                <div class="apple-badge apple-badge-success mt-apple-xxs">
                  <div class="w-1.5 h-1.5 rounded-full bg-current"></div>
                  <span>在线</span>
                </div>
              </div>
            </div>
          </div>

          <!-- 硬件配置 -->
          <div class="apple-card">
            <h2 class="apple-tagline mb-apple-lg">硬件配置</h2>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-apple-lg">
              <div>
                <div class="apple-caption text-apple-ink-muted-48">CPU 核心数</div>
                <div class="apple-headline-md mt-apple-xxs">{{ vmConfig.cpu_cores }} 核</div>
              </div>
              <div>
                <div class="apple-caption text-apple-ink-muted-48">内存</div>
                <div class="apple-headline-md mt-apple-xxs">{{ vmConfig.memory_mb }} MB</div>
              </div>
              <div>
                <div class="apple-caption text-apple-ink-muted-48">系统盘</div>
                <div class="apple-headline-md mt-apple-xxs">{{ vmConfig.disk_gb }} GB SSD</div>
              </div>
              <div>
                <div class="apple-caption text-apple-ink-muted-48">数据盘</div>
                <div class="apple-headline-md mt-apple-xxs">100 GB NVMe</div>
              </div>
            </div>
          </div>

          <!-- 租期信息 -->
          <div class="apple-card">
            <h2 class="apple-tagline mb-apple-lg">租期信息</h2>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-apple-lg">
              <div>
                <div class="apple-caption text-apple-ink-muted-48">创建时间</div>
                <div class="apple-body-strong mt-apple-xxs">{{ leaseInfo.created_at }}</div>
              </div>
              <div>
                <div class="apple-caption text-apple-ink-muted-48">到期时间</div>
                <div class="apple-body-strong mt-apple-xxs">{{ leaseInfo.expires_at }}</div>
              </div>
              <div>
                <div class="apple-caption text-apple-ink-muted-48">剩余时间</div>
                <div class="apple-body-strong mt-apple-xxs text-apple-blue">{{ leaseInfo.remaining }}</div>
              </div>
              <div>
                <div class="apple-caption text-apple-ink-muted-48">远程访问</div>
                <div class="apple-body-strong mt-apple-xxs font-mono">{{ vmInfo.ip_address }}:{{ vmInfo.port }}</div>
              </div>
            </div>
          </div>

          <!-- 操作系统与重装 -->
          <div class="apple-card">
            <h2 class="apple-tagline mb-apple-lg">操作系统</h2>
            <div class="flex items-center justify-between">
              <div>
                <div class="apple-body-strong">{{ vmConfig.os_image }}</div>
                <div class="apple-caption text-apple-ink-muted-48 mt-apple-xxs">Linux 课程定制版</div>
              </div>
              <button @click="showReinstallDialog = true" class="apple-btn-secondary">重装系统</button>
            </div>
          </div>

          <!-- 存储与网络信息 -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-apple-lg">
            <div class="apple-card">
              <h2 class="apple-tagline mb-apple-lg">存储信息</h2>
              <div class="space-y-apple-md">
                <div>
                  <div class="flex justify-between apple-caption mb-apple-xxs">
                    <span>系统盘</span>
                    <span>12.5 GB / 30 GB</span>
                  </div>
                  <div class="apple-progress">
                    <div class="apple-progress-fill" style="width: 42%"></div>
                  </div>
                </div>
                <div>
                  <div class="flex justify-between apple-caption mb-apple-xxs">
                    <span>数据盘</span>
                    <span>23.1 GB / 100 GB</span>
                  </div>
                  <div class="apple-progress">
                    <div class="apple-progress-fill success" style="width: 23%"></div>
                  </div>
                </div>
              </div>
            </div>

            <div class="apple-card">
              <h2 class="apple-tagline mb-apple-lg">网络信息</h2>
              <div class="space-y-apple-sm">
                <div class="flex justify-between">
                  <span class="apple-caption text-apple-ink-muted-48">内网 IP</span>
                  <span class="apple-caption font-mono">{{ vmInfo.ip_address }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="apple-caption text-apple-ink-muted-48">带宽</span>
                  <span class="apple-caption">5 Mbps（双向）</span>
                </div>
                <div class="flex justify-between">
                  <span class="apple-caption text-apple-ink-muted-48">网络模式</span>
                  <span class="apple-caption">NAT（主机转发）</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 网络模块 -->
        <div v-if="activeTab === 'network'" class="space-y-apple-lg">
          <div class="apple-card">
            <h2 class="apple-tagline mb-apple-lg">VPN 配置</h2>
            <div class="flex items-center justify-between mb-apple-lg">
              <div>
                <div class="apple-body-strong">OpenVPN 隧道</div>
                <div class="apple-caption text-apple-ink-muted-48 mt-apple-xxs">已连接 · 延迟 12ms</div>
              </div>
              <div class="apple-badge apple-badge-success">
                <div class="w-1.5 h-1.5 rounded-full bg-current"></div>
                <span>已连接</span>
              </div>
            </div>
            <div class="grid grid-cols-2 gap-apple-md">
              <div>
                <div class="apple-caption text-apple-ink-muted-48">服务器地址</div>
                <div class="apple-body-strong font-mono mt-apple-xxs">vpn.school.local:1194</div>
              </div>
              <div>
                <div class="apple-caption text-apple-ink-muted-48">虚拟 IP</div>
                <div class="apple-body-strong font-mono mt-apple-xxs">10.8.0.15</div>
              </div>
            </div>
          </div>

          <div class="apple-card">
            <h2 class="apple-tagline mb-apple-lg">DHCP 状态</h2>
            <div class="flex items-center justify-between">
              <div>
                <div class="apple-body-strong">自动分配 IP</div>
                <div class="apple-caption text-apple-ink-muted-48 mt-apple-xxs">租期剩余 23 小时 45 分钟</div>
              </div>
              <button class="apple-btn-secondary">续期</button>
            </div>
          </div>

          <div class="apple-card">
            <h2 class="apple-tagline mb-apple-lg">网络流量统计（近 24 小时）</h2>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-apple-lg mb-apple-lg">
              <div class="apple-card-parchment text-center p-apple-md">
                <div class="apple-headline-md text-apple-blue">1.2 GB</div>
                <div class="apple-caption">下载流量</div>
              </div>
              <div class="apple-card-parchment text-center p-apple-md">
                <div class="apple-headline-md text-green-600">450 MB</div>
                <div class="apple-caption">上传流量</div>
              </div>
              <div class="apple-card-parchment text-center p-apple-md">
                <div class="apple-headline-md text-apple-blue">2.3 Mbps</div>
                <div class="apple-caption">平均速率</div>
              </div>
              <div class="apple-card-parchment text-center p-apple-md">
                <div class="apple-headline-md text-green-600">99.8%</div>
                <div class="apple-caption">在线率</div>
              </div>
            </div>

            <div class="h-48 bg-apple-parchment rounded-apple-lg flex items-end justify-between px-apple-md pb-apple-md">
              <div v-for="(bar, i) in networkChartData" :key="i"
                class="flex-1 bg-apple-blue/30 hover:bg-apple-blue/50 transition-all duration-200 rounded-t-apple-xs mx-apple-xxs"
                :style="{ height: bar + '%' }"
                :title="`${23-i}小时前: ${bar}MB`">
              </div>
            </div>
          </div>
        </div>

        <!-- 防火墙模块 -->
        <div v-if="activeTab === 'firewall'" class="space-y-apple-lg">
          <div class="flex items-center justify-between">
            <h2 class="apple-tagline">防火墙规则</h2>
            <button @click="showFirewallDialog = true" class="apple-btn-primary">添加规则</button>
          </div>

          <div class="apple-card">
            <div class="space-y-apple-md">
              <div v-for="rule in firewallRules" :key="rule.id" class="apple-card-parchment p-apple-md">
                <div class="flex items-center justify-between">
                  <div class="flex items-center space-x-apple-md">
                    <div class="w-10 h-10 rounded-full flex items-center justify-center"
                      :class="rule.direction === 'in' ? 'bg-green-100' : 'bg-blue-100'">
                      <svg class="w-5 h-5" :class="rule.direction === 'in' ? 'text-green-600' : 'text-apple-blue'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path v-if="rule.direction === 'in'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 14l-7 7m0 0l-7-7m7 7V3"/>
                        <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 10l7-7m0 0l7 7m-7-7v18"/>
                      </svg>
                    </div>
                    <div>
                      <div class="apple-body-strong">{{ rule.name }}</div>
                      <div class="apple-caption text-apple-ink-muted-48 mt-apple-xxs">
                        {{ rule.direction === 'in' ? '入站' : '出站' }} · {{ rule.protocol.toUpperCase() }} · 端口 {{ rule.port }}
                      </div>
                    </div>
                  </div>
                  <div class="flex items-center space-x-apple-sm">
                    <span class="apple-badge" :class="rule.action === 'accept' ? 'apple-badge-success' : 'apple-badge-error'">
                      {{ rule.action === 'accept' ? '允许' : '拒绝' }}
                    </span>
                    <button @click="deleteFirewallRule(rule.id)" class="apple-btn-utility text-red-500">删除</button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 安全组模块 -->
        <div v-if="activeTab === 'security'" class="space-y-apple-lg">
          <div class="flex items-center justify-between">
            <h2 class="apple-tagline">安全组管理</h2>
            <button @click="showSecurityGroupDialog = true" class="apple-btn-primary">新建安全组</button>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-apple-lg">
            <div v-for="sg in securityGroups" :key="sg.id" class="apple-card">
              <div class="flex items-center justify-between mb-apple-md">
                <h3 class="apple-tagline">{{ sg.name }}</h3>
                <button @click="deleteSecurityGroup(sg.id)" class="apple-btn-utility text-red-500">删除</button>
              </div>
              <p class="apple-caption text-apple-ink-muted-48 mb-apple-lg">{{ sg.description }}</p>
              <div class="space-y-apple-xxs">
                <div v-for="rule in sg.rules" :key="rule.id" class="flex items-center justify-between apple-caption">
                  <span>{{ rule.protocol.toUpperCase() }} :{{ rule.port }}</span>
                  <span class="apple-badge" :class="rule.action === 'allow' ? 'apple-badge-success' : 'apple-badge-error'">
                    {{ rule.action === 'allow' ? '开放' : '关闭' }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 日志模块 -->
        <div v-if="activeTab === 'logs'" class="space-y-apple-lg">
          <div class="flex items-center space-x-apple-sm mb-apple-lg">
            <button v-for="logType in logTypes" :key="logType"
              @click="activeLogType = logType"
              class="apple-btn-secondary text-apple-caption"
              :class="{ 'apple-btn-primary': activeLogType === logType }">
              {{ logType }}
            </button>
          </div>

          <div class="apple-card">
            <div class="space-y-apple-xxs max-h-[600px] overflow-y-auto">
              <div v-for="log in filteredLogs" :key="log.id" class="apple-log-entry flex items-center justify-between" :class="log.level">
                <div class="flex-1">
                  <div class="flex items-center space-x-apple-sm">
                    <span class="apple-caption-strong">{{ log.timestamp }}</span>
                    <span class="apple-caption text-apple-ink-muted-48">·</span>
                    <span class="apple-caption">{{ log.source }}</span>
                  </div>
                  <div class="apple-body mt-apple-xxs">{{ log.message }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 财务模块 -->
        <div v-if="activeTab === 'billing'" class="space-y-apple-lg">
          <div class="grid grid-cols-1 md:grid-cols-3 gap-apple-lg">
            <div class="apple-card text-center">
              <div class="apple-headline-md text-apple-blue">¥{{ billingInfo.balance.toFixed(2) }}</div>
              <div class="apple-caption text-apple-ink-muted-48 mt-apple-xxs">账户余额</div>
            </div>
            <div class="apple-card text-center">
              <div class="apple-headline-md">¥{{ billingInfo.total_cost.toFixed(2) }}</div>
              <div class="apple-caption text-apple-ink-muted-48 mt-apple-xxs">累计消费</div>
            </div>
            <div class="apple-card text-center">
              <div class="apple-headline-md">¥{{ billingInfo.monthly_cost.toFixed(2) }}</div>
              <div class="apple-caption text-apple-ink-muted-48 mt-apple-xxs">本月费用</div>
            </div>
          </div>

          <div class="apple-card">
            <h2 class="apple-tagline mb-apple-lg">消费记录</h2>
            <div class="space-y-apple-md">
              <div v-for="record in billingRecords" :key="record.id" class="flex items-center justify-between apple-card-parchment p-apple-md">
                <div>
                  <div class="apple-body-strong">{{ record.description }}</div>
                  <div class="apple-caption text-apple-ink-muted-48 mt-apple-xxs">{{ record.date }}</div>
                </div>
                <div class="text-right">
                  <div class="apple-body-strong" :class="record.type === 'cost' ? 'text-red-500' : 'text-green-600'">
                    {{ record.type === 'cost' ? '-' : '+' }}¥{{ record.amount.toFixed(2) }}
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="apple-card">
            <h2 class="apple-tagline mb-apple-lg">资源套餐信息</h2>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-apple-lg">
              <div class="apple-card-parchment p-apple-lg">
                <h3 class="apple-tagline mb-apple-md">当前套餐</h3>
                <div class="space-y-apple-sm">
                  <div class="flex justify-between apple-caption">
                    <span>CPU</span>
                    <span>{{ vmConfig.cpu_cores }} 核</span>
                  </div>
                  <div class="flex justify-between apple-caption">
                    <span>内存</span>
                    <span>{{ vmConfig.memory_mb }} MB</span>
                  </div>
                  <div class="flex justify-between apple-caption">
                    <span>存储</span>
                    <span>{{ vmConfig.disk_gb }} GB + 100 GB</span>
                  </div>
                  <div class="flex justify-between apple-caption">
                    <span>带宽</span>
                    <span>5 Mbps</span>
                  </div>
                  <div class="flex justify-between apple-caption-strong border-t border-apple-divider pt-apple-sm mt-apple-sm">
                    <span>月费</span>
                    <span class="text-apple-blue">¥{{ billingInfo.monthly_cost.toFixed(2) }}</span>
                  </div>
                </div>
              </div>
              <div class="apple-card-parchment p-apple-lg">
                <h3 class="apple-tagline mb-apple-md">升级选项</h3>
                <div class="space-y-apple-md">
                  <button class="apple-btn-secondary w-full">升级至 4 核 / ¥45/月</button>
                  <button class="apple-btn-secondary w-full">升级至 4GB 内存 / ¥35/月</button>
                  <button class="apple-btn-secondary w-full">升级至 10 Mbps / ¥20/月</button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- 终端浮动按钮 -->
    <button @click="showTerminal = !showTerminal"
      class="fixed bottom-apple-xl right-apple-xl w-14 h-14 rounded-full bg-apple-blue text-white shadow-apple-subtle flex items-center justify-center hover:scale-105 transition-transform">
      <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
      </svg>
    </button>

    <!-- 终端抽屉 -->
    <transition name="apple-fade">
      <div v-if="showTerminal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
        <div class="apple-card w-full max-w-4xl mx-4 max-h-[80vh]">
          <div class="flex items-center justify-between mb-apple-md">
            <div class="flex items-center space-x-apple-sm">
              <div class="w-2 h-2 rounded-full bg-green-500"></div>
              <h2 class="apple-tagline">实训终端</h2>
            </div>
            <button @click="showTerminal = false" class="apple-btn-utility">关闭</button>
          </div>
          <div class="apple-terminal h-[500px] overflow-y-auto">
            <div v-for="(line, index) in terminalLines" :key="index" class="apple-terminal-line"
              :class="{
                'apple-terminal-error': line.type === 'error',
                'apple-terminal-success': line.type === 'success',
                'apple-terminal-prompt': line.type === 'command'
              }">
              <span v-if="line.type === 'command'">student@vm:~$ </span>
              {{ line.content }}
            </div>
            <div class="flex items-center mt-apple-xs">
              <span class="apple-terminal-prompt">student@vm:~$ </span>
              <input v-model="terminalInput" @keyup.enter="executeCommand" class="apple-terminal-input flex-1 bg-transparent outline-none" />
            </div>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '@/services/api'

const route = useRoute()
const router = useRouter()

const activeTab = ref('overview')
const showTerminal = ref(false)
const terminalInput = ref('')
const showReinstallDialog = ref(false)
const showFirewallDialog = ref(false)
const showSecurityGroupDialog = ref(false)
const activeLogType = ref('操作日志')
const isLoading = ref(false)

const studentInfo = ref({
  id: '',
  name: '',
  classroom: '',
  status: 'offline'
})

const vmConfig = ref({
  cpu_cores: 0,
  memory_mb: 0,
  disk_gb: 0,
  os_image: ''
})

const vmInfo = ref({
  ip_address: '',
  port: 0,
  status: 'stopped',
  uptime_seconds: 0
})

const leaseInfo = ref({
  created_at: '',
  expires_at: '',
  remaining: ''
})

const tabs = [
  { id: 'overview', label: '概览' },
  { id: 'network', label: '网络' },
  { id: 'firewall', label: '防火墙' },
  { id: 'security', label: '安全组' },
  { id: 'logs', label: '日志' },
  { id: 'billing', label: '财务' },
]

const networkChartData = ref([35, 42, 28, 55, 40, 60, 45, 38, 52, 48, 65, 42, 30, 45, 58, 50, 40, 35, 48, 55, 42, 38, 45, 50])

const firewallRules = ref([
  { id: 1, name: '允许 SSH 连接', direction: 'in', protocol: 'tcp', port: '22', action: 'accept' },
  { id: 2, name: '允许 HTTP 访问', direction: 'in', protocol: 'tcp', port: '80', action: 'accept' },
  { id: 3, name: '允许 HTTPS 访问', direction: 'in', protocol: 'tcp', port: '443', action: 'accept' },
  { id: 4, name: '阻止所有入站 ICMP', direction: 'in', protocol: 'icmp', port: '*', action: 'drop' },
])

const securityGroups = ref([
  {
    id: 1,
    name: 'Web 服务器',
    description: '允许 HTTP/HTTPS 流量',
    rules: [
      { id: 1, protocol: 'tcp', port: '80', action: 'allow' },
      { id: 2, protocol: 'tcp', port: '443', action: 'allow' },
    ]
  },
  {
    id: 2,
    name: '数据库',
    description: '仅允许内网访问 MySQL',
    rules: [
      { id: 3, protocol: 'tcp', port: '3306', action: 'allow' },
      { id: 4, protocol: 'tcp', port: '5432', action: 'allow' },
    ]
  }
])

const logTypes = ['操作日志', '登录日志', '命令记录', '资源日志']

const logs = ref({
  '操作日志': [
    { id: 1, timestamp: '2026-05-22 10:23:15', source: '系统', message: '虚拟机启动成功', level: 'success' },
    { id: 2, timestamp: '2026-05-22 09:45:32', source: 'Nginx', message: '配置文件已更新', level: 'info' },
    { id: 3, timestamp: '2026-05-22 09:12:08', source: '系统', message: '安全组规则已修改', level: 'warning' },
  ],
  '登录日志': [
    { id: 4, timestamp: '2026-05-22 08:30:00', source: 'SSH', message: '用户 root 从 192.168.1.100 登录', level: 'info' },
    { id: 5, timestamp: '2026-05-21 17:45:12', source: 'Web', message: '用户 student 登录控制台', level: 'info' },
  ],
  '命令记录': [
    { id: 6, timestamp: '2026-05-22 10:15:23', source: '终端', message: 'sudo systemctl restart nginx', level: 'info' },
    { id: 7, timestamp: '2026-05-22 10:10:45', source: '终端', message: 'vim /etc/nginx/nginx.conf', level: 'info' },
  ],
  '资源日志': [
    { id: 8, timestamp: '2026-05-22 10:00:00', source: '监控', message: 'CPU 使用率 78%', level: 'warning' },
    { id: 9, timestamp: '2026-05-22 09:00:00', source: '监控', message: '内存使用率 65%', level: 'info' },
  ]
})

const billingInfo = ref({
  balance: 125.50,
  total_cost: 234.80,
  monthly_cost: 29.90
})

const billingRecords = ref([
  { id: 1, date: '2026-05-22', description: '虚拟机租赁费', amount: 29.90, type: 'cost' },
  { id: 2, date: '2026-05-15', description: '账户充值', amount: 200.00, type: 'credit' },
  { id: 3, date: '2026-05-15', description: '虚拟机创建费', amount: 4.90, type: 'cost' },
])

const terminalLines = ref([
  { type: 'success', content: 'Welcome to Debian 12.0 (Linux 课程定制版)' },
  { type: 'info', content: 'Last login: Wed May 22 08:30:00 2026 from 192.168.1.100' },
  { type: 'command', content: 'ls -la' },
  { type: 'output', content: 'total 32\ndrwxr-xr-x 4 student student 4096 May 22 09:00 .\ndrwxr-xr-x 3 root    root    4096 May 15 09:00 ..\n-rw-r--r-- 1 student student  220 May 15 09:00 .bash_logout' },
])

const filteredLogs = computed(() => {
  return logs.value[activeLogType.value] || []
})

const loadStudentVMInfo = async () => {
  try {
    const vmid = localStorage.getItem('student_vmid')
    if (!vmid) return

    const res = await api.vms.getDetail(vmid)
    vmInfo.value = {
      ip_address: res.ip_address,
      port: res.port || 8006,
      status: res.status,
      uptime_seconds: res.uptime_seconds || 0
    }

    vmConfig.value = {
      cpu_cores: res.cpu?.cores || 0,
      memory_mb: res.memory_mb || 0,
      disk_gb: res.disk_gb || 0,
      os_image: res.os_image || res.os?.image || ''
    }

    studentInfo.value = {
      id: res.student_id || route.params.studentId,
      name: res.student_name || '',
      classroom: res.classroom || '',
      status: res.status === 'running' ? 'online' : 'offline'
    }

    if (res.lease_info) {
      leaseInfo.value = {
        created_at: res.lease_info.created_at,
        expires_at: res.lease_info.expires_at,
        remaining: res.lease_info.remaining
      }
    }
  } catch (err) {
    console.error('Failed to load VM info:', err)
  }
}

const loadFirewallRules = async () => {
  try {
    const vmid = localStorage.getItem('student_vmid')
    if (!vmid) return

    const res = await api.firewall.getRules()
    firewallRules.value = (res.rules || []).filter((r) => r.vmid === vmid || !r.vmid)
  } catch (err) {
    console.error('Failed to load firewall rules:', err)
  }
}

const executeCommand = async () => {
  if (!terminalInput.value.trim()) return

  const cmd = terminalInput.value.trim()
  terminalLines.value.push({ type: 'command', content: cmd })

  try {
    const vmid = localStorage.getItem('student_vmid')
    if (vmid) {
      const res = await api.vms.executeCommand?.(vmid, cmd)
      if (res.output) {
        terminalLines.value.push({ type: 'output', content: res.output })
      }
    }
  } catch (err) {
    terminalLines.value.push({ type: 'error', content: `执行失败: ${err.message}` })
  }

  terminalInput.value = ''
}

const deleteFirewallRule = async (ruleId) => {
  try {
    // TODO: 调用后端 API 删除防火墙规则
    firewallRules.value = firewallRules.value.filter(r => r.id !== ruleId)
  } catch (err) {
    console.error('Failed to delete firewall rule:', err)
  }
}

const deleteSecurityGroup = async (sgId) => {
  securityGroups.value = securityGroups.value.filter(sg => sg.id !== sgId)
}

const logout = async () => {
  try {
    await api.auth.logout()
  } catch (err) {}
  localStorage.removeItem('student_token')
  localStorage.removeItem('student_vmid')
  router.push({ name: 'StudentLogin' })
}

let refreshInterval = null

onMounted(async () => {
  isLoading.value = true
  try {
    await Promise.all([
      loadStudentVMInfo(),
      loadFirewallRules()
    ])
  } finally {
    isLoading.value = false
  }

  refreshInterval = setInterval(() => {
    loadStudentVMInfo()
  }, 10000)
})

onBeforeUnmount(() => {
  if (refreshInterval) clearInterval(refreshInterval)
})
</script>
