<template>
  <div class="max-w-[1440px] mx-auto px-4 py-apple-section h-[calc(100vh-140px)] flex flex-col">
    <!-- 标题 -->
    <div class="mb-apple-xl flex items-center justify-between">
      <div>
        <h1 class="apple-headline-lg">企业安全沙箱</h1>
        <p class="apple-lead text-apple-ink-muted-48 mt-apple-sm">零信任安全隔离 · 本地 AI 离线审计</p>
      </div>
      <div class="flex items-center space-x-apple-md">
        <div class="apple-badge apple-badge-success">
          <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4z"/>
          </svg>
          <span>零信任沙箱已激活</span>
        </div>
        <div class="apple-badge apple-badge-error">
          <span>拦截 {{ blockCount }} 次</span>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-apple-xl flex-1 min-h-0">
      <!-- 左侧：沙箱工作区 -->
      <div class="lg:col-span-2 relative">
        <div class="apple-card h-full flex flex-col relative overflow-hidden">
          <!-- 动态水印层 -->
          <div class="absolute inset-0 pointer-events-none z-10 overflow-hidden">
            <div v-for="(watermark, index) in watermarks" :key="index"
              class="apple-watermark"
              :style="{ left: watermark.x + '%', top: watermark.y + '%' }">
              CONFIDENTIAL · {{ employeeId }} · {{ currentDate }}
            </div>
          </div>

          <!-- 沙箱头部 -->
          <div class="flex items-center justify-between mb-apple-md relative z-20">
            <div class="flex items-center space-x-apple-sm">
              <div class="w-2 h-2 rounded-full bg-red-500"></div>
              <h2 class="apple-tagline">安全隔离代码沙箱</h2>
              <span class="apple-caption text-apple-ink-muted-48">禁止外发</span>
            </div>
            <div class="flex items-center space-x-apple-md">
              <div class="apple-badge apple-badge-success">
                <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 5.636a9 9 0 010 12.728M5.636 18.364a9 9 0 010-12.728"/>
                </svg>
                <span>断网模式</span>
              </div>
              <div class="apple-badge apple-badge-warning">
                <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                </svg>
                <span>防泄漏保护</span>
              </div>
            </div>
          </div>

          <!-- 代码编辑区 -->
          <div class="apple-terminal flex-1 relative z-20 overflow-y-auto">
            <div class="flex">
              <!-- 行号 -->
              <div class="text-apple-ink-muted-48 select-none pr-apple-md text-right">
                <div v-for="n in codeLines.length" :key="n" class="leading-6">{{ n }}</div>
              </div>
              <!-- 代码内容 -->
              <div class="flex-1">
                <div v-for="(line, index) in codeLines" :key="index" class="leading-6"
                  :class="line.type === 'keyword' ? 'text-apple-terminal-prompt' :
                          line.type === 'string' ? 'apple-terminal-success' :
                          line.type === 'comment' ? 'text-apple-ink-muted-48 italic' : 'text-apple-terminal-output'">
                  {{ line.content }}
                </div>
                <div class="text-apple-blue apple-pulse mt-apple-xs">▋</div>
              </div>
            </div>
          </div>

          <!-- 拦截提示 -->
          <transition name="apple-fade">
            <div v-if="showBlockAlert" class="absolute bottom-apple-md left-1/2 -translate-x-1/2 z-30 apple-alert apple-alert-error px-apple-xl py-apple-md">
              <div class="flex items-center space-x-apple-sm">
                <div class="w-2 h-2 rounded-full bg-red-500 apple-pulse"></div>
                <span class="apple-caption-strong text-red-600">{{ blockMessage }}</span>
              </div>
            </div>
          </transition>
        </div>
      </div>

      <!-- 右侧：AI 合规审计日志 -->
      <div class="flex flex-col">
        <div class="apple-card h-full flex flex-col">
          <div class="flex items-center justify-between mb-apple-md">
            <h2 class="apple-tagline">AI 离线合规审计</h2>
            <div class="apple-badge apple-badge-info">
              <div class="w-1.5 h-1.5 rounded-full bg-current apple-pulse"></div>
              <span>实时监控</span>
            </div>
          </div>

          <!-- 审计统计 -->
          <div class="grid grid-cols-2 gap-apple-sm mb-apple-md">
            <div class="apple-card-parchment text-center p-apple-md">
              <div class="apple-headline-md text-green-600">{{ scanCount }}</div>
              <div class="apple-caption">扫描次数</div>
            </div>
            <div class="apple-card-parchment text-center p-apple-md">
              <div class="apple-headline-md text-red-600">{{ blockCount }}</div>
              <div class="apple-caption">拦截次数</div>
            </div>
            <div class="apple-card-parchment text-center p-apple-md">
              <div class="apple-headline-md text-apple-blue">活跃</div>
              <div class="apple-caption">键盘审计</div>
            </div>
            <div class="apple-card-parchment text-center p-apple-md">
              <div class="apple-headline-md text-green-600">已阻断</div>
              <div class="apple-caption">网络外发</div>
            </div>
          </div>

          <!-- 审计日志流 -->
          <div class="flex-1 overflow-y-auto space-y-apple-xxs mb-apple-md">
            <div v-for="log in auditLogs" :key="log.id" class="apple-log-entry" :class="log.type === 'block' ? 'error' : log.type === 'scan' ? 'success' : log.type === 'warning' ? 'warning' : 'info'">
              <div class="flex items-start space-x-apple-xxs">
                <span class="apple-caption-strong flex-1">{{ log.message }}</span>
                <span class="apple-caption text-apple-ink-muted-48 whitespace-nowrap">{{ log.time }}</span>
              </div>
            </div>
          </div>

          <!-- 安全策略 -->
          <div class="border-t border-apple-divider pt-apple-md">
            <h3 class="apple-caption-strong mb-apple-sm">安全策略状态</h3>
            <div class="space-y-apple-xxs">
              <div v-for="policy in policies" :key="policy.name" class="flex items-center justify-between apple-caption">
                <span>{{ policy.name }}</span>
                <span class="apple-badge" :class="policy.enabled ? 'apple-badge-success' : 'apple-badge-error'">
                  {{ policy.enabled ? '启用' : '禁用' }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const employeeId = 'EMP-2024-0847'
const currentDate = new Date().toISOString().split('T')[0]
const blockCount = ref(23)
const scanCount = ref(1547)
const showBlockAlert = ref(false)
const blockMessage = ref('')

const watermarks = ref([])
for (let i = 0; i < 20; i++) {
  watermarks.value.push({ x: Math.random() * 100, y: Math.random() * 100 })
}

const codeLines = ref([
  { content: '#!/usr/bin/env python3', type: 'comment' },
  { content: '# 企业级数据分析模块 - 机密', type: 'comment' },
  { content: '', type: 'normal' },
  { content: 'import numpy as np', type: 'keyword' },
  { content: 'from sklearn.model_selection import train_test_split', type: 'keyword' },
  { content: '', type: 'normal' },
  { content: 'class DataAnalyzer:', type: 'keyword' },
  { content: '    """核心数据分析引擎"""', type: 'comment' },
  { content: '', type: 'normal' },
  { content: '    def __init__(self, config):', type: 'keyword' },
  { content: '        self.config = config', type: 'normal' },
  { content: '        self.model = None', type: 'normal' },
  { content: '        self.api_key = "sk-xxxx-xxxx-xxxx"', type: 'string' },
  { content: '', type: 'normal' },
  { content: '    def train(self, data):', type: 'keyword' },
  { content: '        X_train, X_test = train_test_split(data)', type: 'normal' },
  { content: '        self.model.fit(X_train)', type: 'normal' },
  { content: '        return self.model.score(X_test)', type: 'normal' },
])

const policies = ref([
  { name: '剪贴板拦截', enabled: true },
  { name: '网络外发阻断', enabled: true },
  { name: '键盘行为审计', enabled: true },
  { name: '动态水印', enabled: true },
  { name: '截屏防护', enabled: true },
  { name: 'USB 存储禁用', enabled: true },
])

const auditLogs = ref([
  { id: 1, time: '14:23:45', type: 'scan', message: '代码行扫描：无敏感信息泄漏' },
  { id: 2, time: '14:23:44', type: 'block', message: '拦截 Ctrl+C 复制操作' },
  { id: 3, time: '14:23:42', type: 'scan', message: '合规检查通过：API 密钥已脱敏' },
  { id: 4, time: '14:23:40', type: 'block', message: '拦截剪贴板写入请求' },
  { id: 5, time: '14:23:38', type: 'warning', message: '检测到异常键盘输入频率' },
  { id: 6, time: '14:23:35', type: 'scan', message: '网络请求检测：无外发连接' },
  { id: 7, time: '14:23:32', type: 'block', message: '拦截 HTTP POST 外发尝试' },
])

const blockTemplates = ['拦截 Ctrl+C 复制操作', '拦截剪贴板写入请求', '拦截网络外发尝试', '检测到截屏行为 - 已模糊化']
const scanTemplates = ['代码扫描完成：合规', '键盘行为审计：正常', '网络检测：无异常连接', '敏感词检测：通过']

let logId = 7
const addAuditLog = () => {
  const now = new Date().toTimeString().slice(0, 8)
  const isBlock = Math.random() > 0.7
  if (isBlock) {
    const msg = blockTemplates[Math.floor(Math.random() * blockTemplates.length)]
    blockCount.value++
    blockMessage.value = msg
    showBlockAlert.value = true
    setTimeout(() => showBlockAlert.value = false, 3000)
    auditLogs.value.unshift({ id: ++logId, time: now, type: 'block', message: msg })
  } else {
    const msg = scanTemplates[Math.floor(Math.random() * scanTemplates.length)]
    scanCount.value++
    auditLogs.value.unshift({ id: ++logId, time: now, type: 'scan', message: msg })
  }
  if (auditLogs.value.length > 50) auditLogs.value.pop()
}

onMounted(() => setInterval(addAuditLog, 1500))
</script>
