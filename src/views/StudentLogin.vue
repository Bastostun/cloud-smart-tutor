<template>
  <div class="min-h-screen bg-gradient-to-br from-[#000000] via-[#0a0a0a] to-[#1d1d1f] flex items-center justify-center p-4">
    <div class="w-full max-w-md">
      <!-- Logo 和标题 -->
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-[#0066cc]/10 mb-4">
          <svg class="w-8 h-8 text-[#0066cc]" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/>
          </svg>
        </div>
        <h1 class="text-2xl font-semibold text-white mb-1">智驭伴学舱</h1>
        <p class="text-sm text-gray-400">云端智师 · 学生实训终端</p>
      </div>

      <!-- 登录表单 -->
      <div v-if="!isCheckingEnvironment" class="bg-[#1d1d1f]/80 backdrop-blur-xl border border-[#333] rounded-2xl p-6 shadow-2xl">
        <!-- 认证方式切换 -->
        <div class="flex items-center space-x-1 mb-6 bg-[#272729] rounded-lg p-1">
          <button v-for="mode in authModes" :key="mode.id"
            @click="authMode = mode.id"
            class="flex-1 py-1.5 text-xs font-medium rounded-md transition-colors"
            :class="authMode === mode.id ? 'bg-[#0066cc] text-white' : 'text-gray-400 hover:text-white'">
            {{ mode.label }}
          </button>
        </div>

        <!-- 账号密码登录 -->
        <form v-if="authMode === 'local'" @submit.prevent="handleLogin" class="space-y-4">
          <div>
            <label class="block text-xs font-medium text-gray-300 mb-1.5">学号</label>
            <input v-model="loginForm.student_id" type="text" 
              class="w-full bg-[#272729] border border-[#333] rounded-lg px-3 py-2.5 text-sm text-white outline-none focus:border-[#0066cc] transition-colors placeholder:text-gray-500"
              placeholder="请输入学号" required />
          </div>

          <div>
            <label class="block text-xs font-medium text-gray-300 mb-1.5">密码</label>
            <input v-model="loginForm.password" type="password" 
              class="w-full bg-[#272729] border border-[#333] rounded-lg px-3 py-2.5 text-sm text-white outline-none focus:border-[#0066cc] transition-colors placeholder:text-gray-500"
              placeholder="请输入密码" required />
          </div>

          <div class="flex items-center justify-between">
            <label class="flex items-center space-x-2 cursor-pointer">
              <input type="checkbox" v-model="loginForm.remember" 
                class="w-4 h-4 rounded bg-[#272729] border-[#333] text-[#0066cc] focus:ring-[#0066cc]" />
              <span class="text-xs text-gray-400">记住登录状态</span>
            </label>
            <a href="#" class="text-xs text-[#0066cc] hover:underline">忘记密码？</a>
          </div>

          <button type="submit" 
            class="w-full py-2.5 text-sm font-medium text-white bg-[#0066cc] hover:bg-[#0071e3] rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            :disabled="isLoggingIn">
            <span v-if="isLoggingIn" class="flex items-center justify-center space-x-2">
              <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <span>登录中...</span>
            </span>
            <span v-else>登录</span>
          </button>
        </form>

        <!-- SSO 登录 -->
        <div v-else class="space-y-4">
          <div class="text-center text-sm text-gray-400 mb-4">
            将重定向到学校统一身份认证平台
          </div>
          <button @click="handleSSOLogin" 
            class="w-full py-2.5 text-sm font-medium text-white bg-[#0066cc] hover:bg-[#0071e3] rounded-lg transition-colors flex items-center justify-center space-x-2">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/>
            </svg>
            <span>跳转至学校 SSO 登录</span>
          </button>
          <div class="text-center text-xs text-gray-500">
            支持 OAuth 2.0 / CAS / LDAP 协议
          </div>
        </div>
      </div>

      <!-- 环境健康检测进度 -->
      <div v-else class="bg-[#1d1d1f]/80 backdrop-blur-xl border border-[#333] rounded-2xl p-6 shadow-2xl">
        <div class="text-center mb-6">
          <div class="inline-flex items-center justify-center w-12 h-12 rounded-full bg-[#0066cc]/10 mb-3">
            <svg class="animate-spin w-6 h-6 text-[#0066cc]" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          </div>
          <h2 class="text-lg font-medium text-white mb-1">环境预检中</h2>
          <p class="text-xs text-gray-400">正在为您准备实训环境...</p>
        </div>

        <!-- 进度条 -->
        <div class="w-full bg-[#272729] rounded-full h-1.5 mb-6">
          <div class="bg-[#0066cc] h-1.5 rounded-full transition-all duration-500"
            :style="{ width: environmentProgress + '%' }"></div>
        </div>

        <!-- 检测项 -->
        <div class="space-y-3">
          <div v-for="check in environmentChecks" :key="check.id" 
            class="flex items-center justify-between p-3 bg-[#272729] rounded-lg">
            <div class="flex items-center space-x-3">
              <div class="w-6 h-6 rounded-full flex items-center justify-center"
                :class="check.status === 'success' ? 'bg-green-500/20' : 
                        check.status === 'error' ? 'bg-red-500/20' : 
                        check.status === 'checking' ? 'bg-yellow-500/20' : 'bg-gray-500/20'">
                <svg v-if="check.status === 'checking'" class="animate-spin w-3 h-3 text-yellow-400" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <svg v-else-if="check.status === 'success'" class="w-3 h-3 text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"/>
                </svg>
                <svg v-else-if="check.status === 'error'" class="w-3 h-3 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M6 18L18 6M6 6l12 12"/>
                </svg>
                <div v-else class="w-2 h-2 rounded-full bg-gray-500"></div>
              </div>
              <div>
                <div class="text-xs font-medium text-white">{{ check.label }}</div>
                <div class="text-[10px] text-gray-500 mt-0.5">{{ check.detail }}</div>
              </div>
            </div>
            <div v-if="check.status === 'success'" class="text-xs font-mono text-green-400">{{ check.value }}</div>
            <div v-else-if="check.status === 'error'" class="text-[10px] text-red-400">异常</div>
          </div>
        </div>

        <!-- 错误提示 -->
        <div v-if="environmentError" class="mt-4 p-3 bg-red-500/10 border border-red-500/20 rounded-lg">
          <div class="flex items-start space-x-2">
            <svg class="w-4 h-4 text-red-400 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
            </svg>
            <div>
              <div class="text-xs font-medium text-red-400">环境检测异常</div>
              <div class="text-[10px] text-gray-400 mt-1">{{ environmentError }}</div>
              <button @click="retryEnvironment" class="mt-2 text-[10px] text-[#0066cc] hover:underline">重试</button>
            </div>
          </div>
        </div>
      </div>

      <!-- 底部协议信息 -->
      <div class="mt-6 text-center">
        <p class="text-xs text-gray-500">
          登录即表示同意
          <a href="#" class="text-[#0066cc] hover:underline">服务条款</a>
          和
          <a href="#" class="text-[#0066cc] hover:underline">隐私政策</a>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const authModes = [
  { id: 'local', label: '学号登录' },
  { id: 'sso', label: '学校 SSO' }
]

const authMode = ref('local')
const isLoggingIn = ref(false)
const isCheckingEnvironment = ref(false)
const environmentProgress = ref(0)
const environmentError = ref('')

const loginForm = reactive({
  student_id: '',
  password: '',
  remember: false
})

const environmentChecks = reactive([
  { id: 'network', label: '边缘节点网络延迟', detail: 'Ping 测试至最近 PVE 节点', status: 'pending', value: '' },
  { id: 'vm', label: '虚拟机虚拟网卡连通性', detail: '检测 vnet0 网卡状态', status: 'pending', value: '' },
  { id: 'ai', label: '本地 AI 智能体推理接口', detail: '健康度检测与响应测试', status: 'pending', value: '' }
])

const handleLogin = async () => {
  isLoggingIn.value = true
  try {
    // 模拟登录 API 调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    localStorage.setItem('student_id', loginForm.student_id)
    localStorage.setItem('student_token', 'mock-token-' + Date.now())
    localStorage.setItem('student_vmid', 'vm-100')
    
    // 登录成功后开始环境检测
    await startEnvironmentCheck()
  } catch (err) {
    console.error('Login failed:', err)
  } finally {
    isLoggingIn.value = false
  }
}

const handleSSOLogin = async () => {
  // 重定向到 SSO 登录页面
  // window.location.href = 'https://sso.school.edu.cn/oauth/authorize?...'
  alert('将重定向到学校 SSO 认证平台\n支持 OAuth 2.0 / CAS / LDAP')
}

const startEnvironmentCheck = async () => {
  isCheckingEnvironment.value = true
  environmentError.value = ''
  
  // 重置检测项
  environmentChecks.forEach(check => check.status = 'pending')
  
  // 1. 检测网络延迟
  environmentChecks[0].status = 'checking'
  environmentProgress.value = 15
  await new Promise(resolve => setTimeout(resolve, 800))
  const latency = Math.floor(Math.random() * 30) + 5
  environmentChecks[0].status = 'success'
  environmentChecks[0].value = `${latency}ms`
  environmentChecks[0].detail = `PVE 节点响应正常`
  environmentProgress.value = 35

  // 2. 检测虚拟机网卡
  environmentChecks[1].status = 'checking'
  environmentProgress.value = 50
  await new Promise(resolve => setTimeout(resolve, 600))
  const isVnetOk = Math.random() > 0.1
  if (isVnetOk) {
    environmentChecks[1].status = 'success'
    environmentChecks[1].value = '已连接'
    environmentChecks[1].detail = 'vnet0 网卡工作正常'
  } else {
    environmentChecks[1].status = 'error'
    environmentChecks[1].detail = '虚拟网卡连通性异常，建议切换备用节点'
  }
  environmentProgress.value = 70

  // 3. 检测 AI 接口
  environmentChecks[2].status = 'checking'
  environmentProgress.value = 80
  await new Promise(resolve => setTimeout(resolve, 500))
  environmentChecks[2].status = 'success'
  environmentChecks[2].value = '健康'
  environmentChecks[2].detail = 'AI 智能体推理接口正常'
  environmentProgress.value = 100

  // 检测完成后判断是否全部成功
  const allSuccess = environmentChecks.every(c => c.status === 'success')
  if (!allSuccess) {
    environmentError.value = '部分环境检测未通过。您可以尝试重试，或联系管理员切换备用节点。'
    return
  }

  // 环境检测通过后，进入实训控制台
  await new Promise(resolve => setTimeout(resolve, 500))
  router.push({ name: 'StudentConsole', params: { studentId: loginForm.student_id } })
}

const retryEnvironment = () => {
  environmentProgress.value = 0
  startEnvironmentCheck()
}
</script>
