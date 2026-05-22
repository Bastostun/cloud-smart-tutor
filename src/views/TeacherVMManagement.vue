<template>
  <div class="max-w-[1440px] mx-auto px-4 py-apple-section">
    <div class="mb-apple-xl flex items-center justify-between">
      <div>
        <h1 class="apple-headline-lg">虚拟机管理</h1>
        <p class="apple-lead text-apple-ink-muted-48 mt-apple-sm">批量分配 · 实时监控 · 统一回收</p>
      </div>
      <button @click="showAllocateDialog = true" class="apple-btn-primary flex items-center space-x-apple-xs">
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
        </svg>
        <span>分配虚拟机</span>
      </button>
    </div>

    <div v-if="isLoading" class="flex items-center justify-center h-64">
      <div class="text-center">
        <svg class="animate-spin w-8 h-8 mx-auto mb-4 text-apple-blue" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <p class="apple-caption text-apple-ink-muted-48">加载虚拟机列表中...</p>
      </div>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-4 gap-apple-lg mb-apple-xl">
      <div class="apple-card text-center">
        <div class="apple-headline-md text-apple-blue">{{ vmStats.total }}</div>
        <div class="apple-caption text-apple-ink-muted-48 mt-apple-xxs">虚拟机总数</div>
      </div>
      <div class="apple-card text-center">
        <div class="apple-headline-md text-green-600">{{ vmStats.running }}</div>
        <div class="apple-caption text-apple-ink-muted-48 mt-apple-xxs">运行中</div>
      </div>
      <div class="apple-card text-center">
        <div class="apple-headline-md text-orange-500">{{ vmStats.stopped }}</div>
        <div class="apple-caption text-apple-ink-muted-48 mt-apple-xxs">已停止</div>
      </div>
      <div class="apple-card text-center">
        <div class="apple-headline-md text-red-500">{{ vmStats.error }}</div>
        <div class="apple-caption text-apple-ink-muted-48 mt-apple-xxs">异常</div>
      </div>
    </div>

    <div class="apple-card">
      <div class="flex items-center justify-between mb-apple-lg">
        <h2 class="apple-tagline">已分配虚拟机</h2>
        <div class="flex items-center space-x-apple-sm">
          <input v-model="searchQuery" class="apple-input text-apple-caption" placeholder="搜索虚拟机..." />
          <select v-model="statusFilter" class="apple-input text-apple-caption">
            <option value="">全部状态</option>
            <option value="running">运行中</option>
            <option value="stopped">已停止</option>
            <option value="error">异常</option>
          </select>
        </div>
      </div>

      <div class="overflow-x-auto">
        <table class="apple-table w-full">
          <thead>
            <tr>
              <th class="apple-caption-strong text-left py-apple-md">虚拟机名称</th>
              <th class="apple-caption-strong text-left">学生</th>
              <th class="apple-caption-strong text-left">状态</th>
              <th class="apple-caption-strong text-left">IP 地址</th>
              <th class="apple-caption-strong text-left">配置</th>
              <th class="apple-caption-strong text-left">创建时间</th>
              <th class="apple-caption-strong text-right">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="vm in filteredVMs" :key="vm.vm_id" class="border-t border-apple-divider">
              <td class="py-apple-md">
                <div class="apple-caption-strong">{{ vm.name }}</div>
                <div class="apple-caption text-apple-ink-muted-48">{{ vm.vm_id }}</div>
              </td>
              <td>
                <div class="apple-caption">{{ vm.student_id }}</div>
                <div class="apple-caption text-apple-ink-muted-48">{{ vm.classroom_id }}</div>
              </td>
              <td>
                <span class="apple-badge" :class="{
                  'apple-badge-success': vm.status === 'running',
                  'apple-badge-warning': vm.status === 'stopped',
                  'apple-badge-error': vm.status === 'error'
                }">
                  <div class="w-1.5 h-1.5 rounded-full bg-current"></div>
                  <span>{{ statusLabel(vm.status) }}</span>
                </span>
              </td>
              <td><span class="apple-caption font-mono">{{ vm.ip_address }}</span></td>
              <td>
                <div class="apple-caption">{{ vm.config.cpu_cores }}核 / {{ vm.config.memory_mb }}MB</div>
                <div class="apple-caption text-apple-ink-muted-48">{{ vm.config.disk_gb }}GB</div>
              </td>
              <td><span class="apple-caption">{{ formatDate(vm.created_at) }}</span></td>
              <td class="text-right">
                <div class="flex items-center justify-end space-x-apple-xxs">
                  <button v-if="vm.status === 'running'" @click="openVNC(vm)" class="apple-btn-utility">远程</button>
                  <button v-if="vm.status === 'stopped'" @click="startVM(vm.vm_id)" class="apple-btn-utility">启动</button>
                  <button v-if="vm.status === 'running'" @click="stopVM(vm.vm_id)" class="apple-btn-utility">停止</button>
                  <button @click="deleteVM(vm.vm_id)" class="apple-btn-utility text-red-500">删除</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <transition name="apple-fade">
      <div v-if="showAllocateDialog" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
        <div class="apple-card w-full max-w-2xl mx-4 max-h-[90vh] overflow-y-auto">
          <div class="flex items-center justify-between mb-apple-lg">
            <h2 class="apple-tagline">批量分配虚拟机</h2>
            <button @click="showAllocateDialog = false" class="apple-btn-utility">关闭</button>
          </div>

          <div class="space-y-apple-lg">
            <div>
              <label class="apple-caption-strong block mb-apple-xxs">班级选择</label>
              <select v-model="allocateForm.classroom_id" class="apple-input w-full">
                <option value="">请选择班级</option>
                <option value="classroom-01">班级 01</option>
                <option value="classroom-02">班级 02</option>
              </select>
            </div>

            <div>
              <label class="apple-caption-strong block mb-apple-xxs">学生选择（多选）</label>
              <div class="grid grid-cols-3 gap-apple-xxs max-h-40 overflow-y-auto">
                <label v-for="student in students" :key="student.id" class="apple-checkbox-label flex items-center space-x-apple-xxs">
                  <input type="checkbox" :value="student.id" v-model="allocateForm.student_ids" class="rounded border-apple-hairline" />
                  <span class="apple-caption">{{ student.name }}</span>
                </label>
              </div>
            </div>

            <div>
              <label class="apple-caption-strong block mb-apple-xxs">操作系统镜像</label>
              <select v-model="allocateForm.os_image" class="apple-input w-full">
                <option value="debian-12.0_x64">Debian 12.0 x64（Linux 课程定制版）</option>
                <option value="ubuntu-22.04_x64">Ubuntu 22.04 LTS</option>
                <option value="centos-9_x64">CentOS Stream 9</option>
              </select>
            </div>

            <div class="grid grid-cols-3 gap-apple-md">
              <div>
                <label class="apple-caption-strong block mb-apple-xxs">CPU 核心数</label>
                <input v-model.number="allocateForm.cpu_cores" type="number" min="1" max="16" class="apple-input w-full" />
              </div>
              <div>
                <label class="apple-caption-strong block mb-apple-xxs">内存 (MB)</label>
                <input v-model.number="allocateForm.memory_mb" type="number" min="512" max="32768" class="apple-input w-full" />
              </div>
              <div>
                <label class="apple-caption-strong block mb-apple-xxs">存储 (GB)</label>
                <input v-model.number="allocateForm.disk_gb" type="number" min="10" max="500" class="apple-input w-full" />
              </div>
            </div>

            <div class="border-t border-apple-divider pt-apple-lg">
              <div class="flex items-center justify-between mb-apple-md">
                <label class="apple-tagline">AI 伴学功能</label>
                <label class="relative inline-flex items-center cursor-pointer">
                  <input type="checkbox" v-model="allocateForm.ai_enabled" class="sr-only peer" />
                  <div class="w-11 h-6 bg-apple-divider peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-apple-blue"></div>
                </label>
              </div>

              <div v-if="allocateForm.ai_enabled" class="space-y-apple-md">
                <div>
                  <label class="apple-caption-strong block mb-apple-xxs">AI 模型选择</label>
                  <select v-model="allocateForm.ai_model" class="apple-input w-full">
                    <option value="ollama-llama3">Ollama Llama 3（本地部署）</option>
                    <option value="ollama-qwen">Ollama Qwen（中文优化）</option>
                    <option value="custom-api">自定义 API</option>
                  </select>
                </div>
                <div>
                  <label class="apple-caption-strong block mb-apple-xxs">API 服务地址</label>
                  <input v-model="allocateForm.ai_api_endpoint" class="apple-input w-full" placeholder="http://localhost:11434" />
                </div>
              </div>
            </div>
          </div>

          <div class="flex items-center justify-end space-x-apple-sm mt-apple-xl">
            <button @click="showAllocateDialog = false" class="apple-btn-secondary">取消</button>
            <button @click="allocateVMs" class="apple-btn-primary">确认分配</button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { api } from '@/services/api'

const isLoading = ref(false)
const showAllocateDialog = ref(false)
const searchQuery = ref('')
const statusFilter = ref('')

const vmStats = ref({
  total: 0,
  running: 0,
  stopped: 0,
  error: 0
})

const vms = ref([])
const students = ref([])

const allocateForm = ref({
  student_ids: [],
  classroom_id: '',
  os_image: 'debian-12.0_x64',
  cpu_cores: 2,
  memory_mb: 2048,
  disk_gb: 30,
  ai_enabled: true,
  ai_model: 'ollama-llama3',
  ai_api_endpoint: 'http://localhost:11434'
})

const filteredVMs = computed(() => {
  let result = vms.value

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(vm =>
      vm.name.toLowerCase().includes(query) ||
      vm.student_id.toLowerCase().includes(query) ||
      vm.vm_id.toLowerCase().includes(query)
    )
  }

  if (statusFilter.value) {
    result = result.filter(vm => vm.status === statusFilter.value)
  }

  return result
})

const statusLabel = (status) => {
  const labels = { running: '运行中', stopped: '已停止', error: '异常', paused: '已暂停' }
  return labels[status] || status
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

const loadStudents = async () => {
  try {
    const res = await api.permissions.getUsers()
    students.value = (res.users || []).filter(u => u.role === 'student')
  } catch (err) {
    console.error('Failed to load students:', err)
  }
}

const loadVMs = async () => {
  try {
    const res = await api.vms.getList()
    vms.value = res.vms || []

    vmStats.value.total = vms.value.length
    vmStats.value.running = vms.value.filter(v => v.status === 'running').length
    vmStats.value.stopped = vms.value.filter(v => v.status === 'stopped').length
    vmStats.value.error = vms.value.filter(v => v.status === 'error').length
  } catch (err) {
    console.error('Failed to load VMs:', err)
  }
}

const allocateVMs = async () => {
  try {
    for (const studentId of allocateForm.value.student_ids) {
      await api.vms.create({
        student_id: studentId,
        classroom_id: allocateForm.value.classroom_id,
        os_image: allocateForm.value.os_image,
        cpu_cores: allocateForm.value.cpu_cores,
        memory_mb: allocateForm.value.memory_mb,
        disk_gb: allocateForm.value.disk_gb,
        ai_enabled: allocateForm.value.ai_enabled,
        ai_model: allocateForm.value.ai_model,
        ai_api_endpoint: allocateForm.value.ai_api_endpoint
      })
    }
    showAllocateDialog.value = false
    loadVMs()
  } catch (err) {
    console.error('Failed to allocate VMs:', err)
  }
}

const startVM = async (vmId) => {
  try {
    await api.vms.start(vmId)
    loadVMs()
  } catch (err) {
    console.error('Failed to start VM:', err)
  }
}

const stopVM = async (vmId) => {
  try {
    await api.vms.stop(vmId)
    loadVMs()
  } catch (err) {
    console.error('Failed to stop VM:', err)
  }
}

const deleteVM = async (vmId) => {
  if (!confirm('确定要删除此虚拟机吗？')) return
  try {
    await api.vms.delete(vmId)
    loadVMs()
  } catch (err) {
    console.error('Failed to delete VM:', err)
  }
}

const openVNC = async (vm) => {
  try {
    const res = await api.vms.getConsole(vm.vm_id)
    window.open(res.url, '_blank', 'width=1280,height=720')
  } catch (err) {
    console.error('Failed to open VNC:', err)
  }
}

let refreshInterval = null

onMounted(async () => {
  isLoading.value = true
  try {
    await Promise.all([
      loadVMs(),
      loadStudents()
    ])
  } finally {
    isLoading.value = false
  }

  refreshInterval = setInterval(() => {
    loadVMs()
  }, 10000)
})

onBeforeUnmount(() => {
  if (refreshInterval) clearInterval(refreshInterval)
})
</script>
