<template>
  <div class="min-h-screen bg-apple-canvas text-apple-ink">
    <!-- Apple 全球导航栏 -->
    <nav class="apple-global-nav fixed top-0 left-0 right-0 z-50">
      <div class="max-w-[1440px] mx-auto px-4 h-full flex items-center justify-between">
        <div class="flex items-center space-x-6">
          <!-- Logo -->
          <router-link to="/idc-dashboard" class="flex items-center space-x-2 text-apple-nav font-display">
            <svg class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/>
            </svg>
            <span class="text-white/90 hover:text-white transition-colors">云端智师</span>
          </router-link>
          
          <!-- 导航链接 -->
          <div class="hidden md:flex items-center space-x-5">
            <router-link v-for="item in navItems" :key="item.path" :to="item.path"
              class="text-white/80 hover:text-white transition-colors text-apple-nav"
              :class="{ 'text-white': route.path === item.path }">
              {{ item.label }}
            </router-link>
          </div>
        </div>
        
        <div class="flex items-center space-x-4">
          <!-- 教室状态 -->
          <div class="flex items-center space-x-2 text-white/80">
            <div class="w-1.5 h-1.5 rounded-full bg-green-400"></div>
            <span class="text-apple-nav">局域网</span>
          </div>
          <!-- 时钟 -->
          <span class="text-white/60 text-apple-nav hidden lg:block">{{ currentTime }}</span>
        </div>
      </div>
    </nav>

    <!-- Apple 子导航栏（毛玻璃效果） -->
    <div class="apple-sub-nav fixed top-[44px] left-0 right-0 z-40">
      <div class="max-w-[1440px] mx-auto px-4 h-full flex items-center justify-between">
        <div class="flex items-center space-x-4">
          <h2 class="text-apple-tagline text-apple-ink">{{ currentPageTitle }}</h2>
          <span class="text-apple-caption">·</span>
          <span class="text-apple-caption">{{ currentPageSubtitle }}</span>
        </div>
        <div class="flex items-center space-x-3">
          <!-- 教室 LAN 状态指示 -->
          <div class="flex items-center space-x-2 apple-badge apple-badge-success">
            <div class="w-1.5 h-1.5 rounded-full bg-current"></div>
            <span>教室集群在线</span>
          </div>
          <button class="apple-btn-utility">帮助</button>
        </div>
      </div>
    </div>

    <!-- 主内容区 -->
    <main class="pt-[96px] min-h-screen pb-[44px]">
      <router-view v-slot="{ Component }">
        <transition name="apple-fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>

    <!-- Toast 通知 -->
    <Toast />

    <!-- 底部状态栏 -->
    <footer class="fixed bottom-0 left-0 right-0 bg-apple-parchment border-t border-apple-hairline px-6 py-2">
      <div class="max-w-[1440px] mx-auto flex items-center justify-between text-apple-fine text-apple-ink-muted-48">
        <div class="flex items-center space-x-6">
          <span>节点: <span class="text-apple-ink font-medium">52</span></span>
          <span>虚拟机: <span class="text-apple-ink font-medium">1,247</span></span>
          <span>KSM 合并率: <span class="text-apple-ink font-medium">62.3%</span></span>
          <span>教室 LAN: <span class="text-apple-ink font-medium">活跃</span></span>
        </div>
        <div class="flex items-center space-x-4">
          <span>v2.0.0 Apple Design</span>
          <span>局域网集群配置</span>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import Toast from '@/components/Toast.vue'

const route = useRoute()
const currentTime = ref('')

const navItems = [
  { path: '/idc-dashboard', label: 'IDC 总控' },
  { path: '/student/workspace', label: '学生孪生舱' },
  { path: '/teacher/audit', label: '教改审计' },
  { path: '/enterprise/sandbox', label: '企业沙箱' }
]

const pageTitles = {
  '/idc-dashboard': { title: 'IDC 全局拓扑', subtitle: '教室局域网集群监控' },
  '/student/workspace': { title: '学生伴学孪生舱', subtitle: 'AI 伴随式智能排错' },
  '/teacher/audit': { title: '教改审计面板', subtitle: '高频错误聚类分析' },
  '/enterprise/sandbox': { title: '企业安全沙箱', subtitle: '零信任安全隔离' }
}

const currentPageTitle = computed(() => pageTitles[route.path]?.title || '云端智师')
const currentPageSubtitle = computed(() => pageTitles[route.path]?.subtitle || '')

const updateClock = () => {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  })
}

let timer = null
onMounted(() => {
  updateClock()
  timer = setInterval(updateClock, 60000)
})

onUnmounted(() => {
  clearInterval(timer)
})
</script>

<style>
.apple-fade-enter-active,
.apple-fade-leave-active {
  transition: opacity 0.2s ease;
}

.apple-fade-enter-from,
.apple-fade-leave-to {
  opacity: 0;
}
</style>
