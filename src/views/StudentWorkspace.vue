<template>
  <div class="max-w-[1440px] mx-auto px-4 py-apple-section h-[calc(100vh-140px)] flex flex-col">
    <!-- 标题 -->
    <div class="mb-apple-xl flex items-center justify-between">
      <div>
        <h1 class="apple-headline-lg">学生伴学孪生舱</h1>
        <p class="apple-lead text-apple-ink-muted-48 mt-apple-sm">AI 伴随式智能排错 · 实时错误捕获</p>
      </div>
      <div class="flex items-center space-x-apple-md">
        <div class="apple-badge apple-badge-success">
          <div class="w-1.5 h-1.5 rounded-full bg-current"></div>
          <span>AI Agent 在线</span>
        </div>
        <div class="apple-badge apple-badge-warning">
          <span>捕获错误 3</span>
        </div>
      </div>
    </div>

    <!-- 分屏布局 -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-apple-xl flex-1 min-h-0">
      <!-- 左侧：实训终端 -->
      <div class="flex flex-col">
        <div class="apple-card flex-1 flex flex-col min-h-0">
          <div class="flex items-center justify-between mb-apple-md">
            <div class="flex items-center space-x-apple-sm">
              <div class="w-2 h-2 rounded-full bg-green-500"></div>
              <h2 class="apple-tagline">实训终端</h2>
              <span class="apple-caption text-apple-ink-muted-48">Ubuntu 22.04 LTS</span>
            </div>
            <div class="apple-badge apple-badge-success">
              <span>运行中</span>
            </div>
          </div>

          <div class="apple-terminal flex-1 overflow-y-auto">
            <div v-for="(line, index) in terminalLines" :key="index" class="apple-terminal-line"
              :class="line.type === 'error' ? 'apple-terminal-error' :
                      line.type === 'success' ? 'apple-terminal-success' :
                      line.type === 'command' ? 'apple-terminal-prompt' : 'apple-terminal-output'">
              <span v-if="line.type === 'command'">student@vm-042:~$ </span>
              {{ line.content }}
            </div>
            <div class="apple-terminal-line flex items-center mt-apple-xs">
              <span class="apple-terminal-prompt">student@vm-042:~$ </span>
              <span class="ml-1">{{ currentInput }}</span>
              <span class="inline-block w-2 h-4 bg-apple-blue ml-0.5 apple-pulse"></span>
            </div>
          </div>

          <!-- 错误捕获提示 -->
          <div v-if="lastError" class="mt-apple-md apple-alert apple-alert-error flex items-center space-x-apple-sm">
            <div class="w-2 h-2 rounded-full bg-red-500 apple-pulse"></div>
            <span class="apple-caption-strong text-red-600">探针已捕获 stderr</span>
            <span class="apple-caption text-apple-ink-muted-48">→ 已发送至 AI Agent</span>
          </div>
        </div>
      </div>

      <!-- 右侧：AI 伴学区 -->
      <div class="flex flex-col">
        <div class="apple-card flex-1 flex flex-col min-h-0">
          <div class="flex items-center justify-between mb-apple-md">
            <div class="flex items-center space-x-apple-sm">
              <h2 class="apple-tagline">AI 伴学助手</h2>
            </div>
            <div class="apple-badge apple-badge-info">
              <div class="w-1.5 h-1.5 rounded-full bg-current apple-pulse"></div>
              <span>实时监控中</span>
            </div>
          </div>

          <!-- AI 对话区 -->
          <div class="flex-1 overflow-y-auto space-y-apple-md mb-apple-md">
            <div v-for="msg in chatMessages" :key="msg.id"
              class="flex" :class="msg.sender === 'ai' ? 'justify-start' : 'justify-end'">
              <div class="max-w-[85%]"
                :class="msg.sender === 'ai' ? 'apple-chat-ai' : 'apple-chat-user'">
                <div v-if="msg.sender === 'ai'" class="flex items-center space-x-apple-xxs mb-apple-xxs">
                  <span class="apple-caption-strong text-apple-ink">AI Agent</span>
                  <span v-if="msg.proactive" class="apple-badge apple-badge-warning text-[10px]">主动捕获</span>
                </div>
                <p class="text-apple-body whitespace-pre-wrap">{{ msg.content }}</p>
                <div v-if="msg.solution" class="mt-apple-md apple-card-parchment">
                  <div class="flex items-center space-x-apple-xxs mb-apple-xxs">
                    <div class="w-1.5 h-1.5 rounded-full bg-green-500"></div>
                    <span class="apple-caption-strong text-green-700">解决方案</span>
                  </div>
                  <div class="apple-terminal text-sm">
                    <pre class="whitespace-pre-wrap text-apple-terminal-output">{{ msg.solution }}</pre>
                  </div>
                </div>
                <div class="apple-caption text-apple-ink-muted-48 mt-apple-xxs text-right">{{ msg.time }}</div>
              </div>
            </div>
          </div>

          <!-- 输入区 -->
          <div class="flex items-center space-x-apple-sm">
            <input v-model="userInput" @keyup.enter="sendMessage"
              class="apple-input flex-1"
              placeholder="向 AI Agent 提问..." />
            <button @click="sendMessage" class="apple-btn-primary">发送</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const terminalLines = ref([
  { content: 'Welcome to Ubuntu 22.04.3 LTS', type: 'info' },
  { content: '', type: 'info' },
])

const chatMessages = ref([
  { id: 1, sender: 'ai', content: '你好！我是你的 AI 伴学助手。我会实时监听你的终端操作，遇到问题随时向我提问！', time: new Date().toTimeString().slice(0, 5) }
])

const userInput = ref('')
const currentInput = ref('')
const lastError = ref(null)

const terminalSequence = [
  { type: 'command', content: 'sudo nginx -t', delay: 800 },
  { type: 'info', content: '[sudo] password for student: ', delay: 400 },
  { type: 'info', content: 'nginx: the configuration file /etc/nginx/nginx.conf syntax is ok', delay: 600 },
  { type: 'success', content: 'nginx: configuration file /etc/nginx/nginx.conf test is successful', delay: 400 },
  { type: 'command', content: 'sudo systemctl restart nginx', delay: 1000 },
  { type: 'error', content: 'Job for nginx.service failed because the control process exited with error code.', delay: 600 },
  { type: 'error', content: 'See "systemctl status nginx.service" and "journalctl -xe" for details.', delay: 400 },
  { type: 'command', content: 'sudo journalctl -xe | tail -20', delay: 800 },
  { type: 'error', content: 'nginx: [emerg] unexpected end of file, expecting ";" or "}" in /etc/nginx/sites-enabled/default:42', delay: 600 },
  { type: 'command', content: 'sudo nano /etc/nginx/sites-enabled/default', delay: 600 },
  { type: 'info', content: '(编辑配置文件...)', delay: 1000 },
  { type: 'command', content: 'sudo systemctl restart nginx', delay: 800 },
  { type: 'success', content: 'nginx.service started successfully.', delay: 600 },
  { type: 'command', content: 'curl -I http://localhost', delay: 600 },
  { type: 'success', content: 'HTTP/1.1 200 OK', delay: 400 },
]

let lineIndex = 0
let currentCharIndex = 0

const simulateTyping = () => {
  if (lineIndex >= terminalSequence.length) return
  const line = terminalSequence[lineIndex]
  if (line.type === 'command') {
    currentInput.value = line.content.slice(0, currentCharIndex + 1)
    currentCharIndex++
    if (currentCharIndex >= line.content.length) {
      terminalLines.value.push({ content: line.content, type: 'command' })
      currentInput.value = ''
      currentCharIndex = 0
      lineIndex++
      if (line.content.includes('systemctl restart nginx') && lineIndex === 6) {
        lastError.value = { time: new Date().toTimeString().slice(0, 8) }
        triggerAIResponse()
      }
      if (line.content.includes('curl -I')) triggerSuccessResponse()
      setTimeout(simulateTyping, 500)
    } else {
      setTimeout(simulateTyping, 80 + Math.random() * 40)
    }
  } else {
    terminalLines.value.push({ content: line.content, type: line.type })
    lineIndex++
    if (line.type === 'error') lastError.value = { time: new Date().toTimeString().slice(0, 8) }
    setTimeout(simulateTyping, line.delay)
  }
}

const triggerAIResponse = () => {
  chatMessages.value.push({
    id: chatMessages.value.length + 1, sender: 'ai', proactive: true,
    content: '我检测到 Nginx 启动失败！通过分析错误日志，我发现问题出在配置文件第 42 行缺少分号。',
    solution: '在 /etc/nginx/sites-enabled/default 第 42 行末尾添加分号：\nserver {\n    listen 80;\n    server_name example.com;\n}',
    time: new Date().toTimeString().slice(0, 5)
  })
}

const triggerSuccessResponse = () => {
  chatMessages.value.push({
    id: chatMessages.value.length + 1, sender: 'ai',
    content: 'Nginx 已成功重启并返回 200 OK。你解决了那个分号问题，做得漂亮！',
    time: new Date().toTimeString().slice(0, 5)
  })
}

const sendMessage = () => {
  if (!userInput.value.trim()) return
  chatMessages.value.push({ id: chatMessages.value.length + 1, sender: 'user', content: userInput.value, time: new Date().toTimeString().slice(0, 5) })
  userInput.value = ''
  setTimeout(() => {
    chatMessages.value.push({
      id: chatMessages.value.length + 1, sender: 'ai',
      content: '收到你的问题，让我分析一下...\n\n根据当前上下文，我建议你先检查 Nginx 配置文件中的语法错误。使用 nginx -t 命令可以快速验证配置是否正确。',
      time: new Date().toTimeString().slice(0, 5)
    })
  }, 1000)
}

onMounted(() => setTimeout(simulateTyping, 1500))
</script>
