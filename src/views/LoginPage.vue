<template>
  <div class="min-h-screen bg-gradient-to-br from-apple-ink via-apple-tile1 to-apple-tile2 flex items-center justify-center px-4">
    <div class="w-full max-w-md">
      <!-- Logo -->
      <div class="text-center mb-apple-xl">
        <svg class="w-20 h-20 mx-auto mb-apple-md text-white" viewBox="0 0 24 24" fill="currentColor">
          <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/>
        </svg>
        <h1 class="apple-display-md text-white mb-apple-sm">云端智师</h1>
        <p class="apple-lead text-white/60">Proxmox VE 管理控制台</p>
        <p class="apple-fine-print text-white/40 mt-apple-xxs">v8.2.4</p>
      </div>

      <!-- 登录卡片 -->
      <div class="apple-card bg-apple-tile1 border-apple-tile2 p-apple-xl">
        <form @submit.prevent="handleLogin" class="space-y-apple-lg">
          <!-- 语言选择 -->
          <div>
            <label class="apple-caption-strong text-white block mb-apple-xxs">语言</label>
            <select v-model="loginForm.language" class="apple-input w-full bg-apple-tile2 border-apple-hairline text-white">
              <option value="zh-CN">简体中文</option>
              <option value="en">English</option>
              <option value="de">Deutsch</option>
              <option value="fr">Français</option>
              <option value="es">Español</option>
              <option value="ja">日本語</option>
              <option value="ko">한국어</option>
              <option value="ru">Русский</option>
            </select>
          </div>

          <!-- 用户名 -->
          <div>
            <label class="apple-caption-strong text-white block mb-apple-xxs">用户名</label>
            <input v-model="loginForm.username" type="text" class="apple-input w-full bg-apple-tile2 border-apple-hairline text-white" placeholder="admin" required />
          </div>

          <!-- 密码 -->
          <div>
            <label class="apple-caption-strong text-white block mb-apple-xxs">密码</label>
            <input v-model="loginForm.password" type="password" class="apple-input w-full bg-apple-tile2 border-apple-hairline text-white" placeholder="••••••••" required />
          </div>

          <!-- 认证域 -->
          <div>
            <label class="apple-caption-strong text-white block mb-apple-xxs">认证方式</label>
            <select v-model="loginForm.realm" class="apple-input w-full bg-apple-tile2 border-apple-hairline text-white">
              <option value="pve">Proxmox VE 认证</option>
              <option value="pam">Linux PAM 认证</option>
              <option value="ldap">LDAP 目录服务</option>
              <option value="ad">微软 Active Directory</option>
            </select>
          </div>

          <!-- 双因子认证 -->
          <div v-if="showOTP">
            <label class="apple-caption-strong text-white block mb-apple-xxs">双因子验证码 (OTP)</label>
            <input v-model="loginForm.otp" type="text" class="apple-input w-full bg-apple-tile2 border-apple-hairline text-white" placeholder="123456" maxlength="6" />
            <p class="apple-fine-print text-white/40 mt-apple-xxs">请输入 OATH Token 或 YubiKey 生成的 6 位验证码</p>
          </div>

          <!-- 保存用户名 -->
          <div class="flex items-center space-x-apple-xxs">
            <input type="checkbox" v-model="loginForm.remember" class="rounded border-apple-hairline" />
            <label class="apple-caption text-white/80">保存用户名</label>
          </div>

          <!-- 错误提示 -->
          <div v-if="errorMsg" class="apple-alert apple-alert-error">
            <div class="flex items-center space-x-apple-sm">
              <div class="w-2 h-2 rounded-full bg-red-500 apple-pulse"></div>
              <span class="apple-caption-strong text-red-400">{{ errorMsg }}</span>
            </div>
          </div>

          <!-- 登录按钮 -->
          <button type="submit" class="apple-btn-primary w-full bg-apple-blue-focus hover:bg-apple-blue" :disabled="isLoggingIn">
            <span v-if="isLoggingIn">
              <svg class="animate-spin -ml-1 mr-apple-sm h-4 w-4 text-white inline" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              登录中...
            </span>
            <span v-else>登录</span>
          </button>
        </form>
      </div>

      <!-- 底部信息 -->
      <div class="mt-apple-xl text-center">
        <p class="apple-fine-print text-white/40">
          云端智师 · 边缘超融合微型 IDC · 基于 Proxmox VE
        </p>
        <p class="apple-fine-print text-white/30 mt-apple-xxs">
          © 2026 Cloud Smart Tutor
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '@/services/api'

const router = useRouter()
const isLoggingIn = ref(false)
const errorMsg = ref('')
const showOTP = ref(false)

const loginForm = ref({
  username: '',
  password: '',
  realm: 'pve',
  otp: '',
  language: 'zh-CN',
  remember: false
})

const handleLogin = async () => {
  isLoggingIn.value = true
  errorMsg.value = ''

  try {
    const res = await api.auth.login(loginForm.value)

    if (res.error) {
      if (res.error.includes('OTP') || res.error.includes('双因子')) {
        showOTP.value = true
        errorMsg.value = '请输入双因子验证码'
      } else {
        errorMsg.value = res.error
      }
      return
    }

    localStorage.setItem('pve_token', res.token)
    localStorage.setItem('pve_user', JSON.stringify(res.user))

    if (loginForm.value.remember) {
      localStorage.setItem('pve_saved_user', loginForm.value.username)
    }

    router.push({ name: 'TeacherConsole' })
  } catch (err) {
    errorMsg.value = '登录失败：' + err.message
  } finally {
    isLoggingIn.value = false
  }
}

onMounted(() => {
  const savedUser = localStorage.getItem('pve_saved_user')
  if (savedUser) {
    loginForm.value.username = savedUser
    loginForm.value.remember = true
  }
})
</script>
