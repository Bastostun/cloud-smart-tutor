import { mockApi } from './mock-api'

const API_BASE = import.meta.env.VITE_API_BASE || '/api/v1'
const USE_MOCK = import.meta.env.VITE_USE_MOCK === 'true'

const getAuthHeader = () => {
  const token = localStorage.getItem('pve_token')
  return token ? { Authorization: `Bearer ${token}` } : {}
}

const request = async (method, path, data = {}) => {
  if (USE_MOCK) {
    switch (method) {
      case 'GET': return mockApi.get(path, data)
      case 'POST': return mockApi.post(path, data)
      case 'PUT': return mockApi.put(path, data)
      case 'DELETE': return mockApi.delete(path)
      default: throw new Error(`Unsupported method: ${method}`)
    }
  }

  let url = path.startsWith('http') ? path : `${API_BASE}${path}`
  const options = {
    method,
    headers: {
      'Content-Type': 'application/json',
      ...getAuthHeader()
    }
  }

  if (method === 'GET' && Object.keys(data).length > 0) {
    const params = new URLSearchParams(data)
    url = `${url}?${params}`
  } else if (method !== 'GET' && data) {
    options.body = JSON.stringify(data)
  }

  const response = await fetch(url, options)

  if (!response.ok) {
    if (response.status === 401) {
      localStorage.removeItem('pve_token')
      localStorage.removeItem('pve_user')
      window.location.href = '/login'
    }
    const errorData = await response.json().catch(() => ({}))
    throw new Error(errorData.error || `HTTP ${response.status}: ${response.statusText}`)
  }

  return response.json()
}

export const api = {
  // 认证
  auth: {
    login: (data) => request('POST', '/auth/login', data),
    logout: () => request('POST', '/auth/logout'),
    verify: () => request('GET', '/auth/verify')
  },

  // 数据中心
  datacenter: {
    getSummary: () => request('GET', '/datacenter/summary'),
    getNodes: () => request('GET', '/datacenter/nodes'),
    search: (q) => request('GET', '/datacenter/search', { q })
  },

  // 虚拟机
  vms: {
    getList: (params) => request('GET', '/vms', params),
    getDetail: (vmid) => request('GET', `/vms/${vmid}`),
    create: (data) => request('POST', '/vms', data),
    start: (vmid) => request('POST', `/vms/${vmid}/start`),
    stop: (vmid) => request('POST', `/vms/${vmid}/stop`),
    reboot: (vmid) => request('POST', `/vms/${vmid}/reboot`),
    suspend: (vmid) => request('POST', `/vms/${vmid}/suspend`),
    resume: (vmid) => request('POST', `/vms/${vmid}/resume`),
    delete: (vmid) => request('DELETE', `/vms/${vmid}`),
    clone: (vmid, data) => request('POST', `/vms/${vmid}/clone`, data),
    migrate: (vmid, data) => request('POST', `/vms/${vmid}/migrate`, data),
    getConsole: (vmid) => request('GET', `/vms/${vmid}/console/novnc`),
    getOptions: (vmid) => request('GET', `/vms/${vmid}/options`),
    updateOptions: (vmid, data) => request('PUT', `/vms/${vmid}/options`, data),
    getSnapshots: (vmid) => request('GET', `/vms/${vmid}/snapshots`),
    createSnapshot: (vmid, data) => request('POST', `/vms/${vmid}/snapshots`, data),
    rollbackSnapshot: (vmid, snapname) => request('POST', `/vms/${vmid}/snapshots/${snapname}/rollback`),
    deleteSnapshot: (vmid, snapname) => request('DELETE', `/vms/${vmid}/snapshots/${snapname}`),
    getOSImages: () => request('GET', '/vm/os-images'),
    getAIModels: () => request('GET', '/vm/ai-models'),
    getStatus: (vmid) => request('GET', `/vms/${vmid}/status`),
    reinstall: (vmid, data) => request('POST', `/vms/${vmid}/reinstall`, data)
  },

  // 节点
  nodes: {
    getSummary: (nodeId) => request('GET', `/nodes/${nodeId}/summary`),
    getNetwork: (nodeId) => request('GET', `/nodes/${nodeId}/network`),
    getDisks: (nodeId) => request('GET', `/nodes/${nodeId}/disks`),
    getSyslog: (nodeId) => request('GET', `/nodes/${nodeId}/syslog`),
    getTasks: (nodeId) => request('GET', `/nodes/${nodeId}/tasks`)
  },

  // 存储
  storage: {
    getList: () => request('GET', '/storage'),
    getDetail: (storageId) => request('GET', `/storage/${storageId}`),
    getContent: (storageId) => request('GET', `/storage/${storageId}/content`)
  },

  // 任务
  tasks: {
    getRunning: () => request('GET', '/tasks/running'),
    getLog: (taskId) => request('GET', `/tasks/${taskId}/log`),
    stop: (taskId) => request('POST', `/tasks/${taskId}/stop`)
  },

  // 权限
  permissions: {
    getUsers: () => request('GET', '/users'),
    getGroups: () => request('GET', '/groups'),
    getRoles: () => request('GET', '/roles')
  },

  // 防火墙
  firewall: {
    getRules: () => request('GET', '/firewall/rules'),
    createRule: (data) => request('POST', '/firewall/rules', data)
  },

  // 设置
  settings: {
    get: () => request('GET', '/settings'),
    update: (data) => request('PUT', '/settings', data)
  },

  // 健康检查
  health: () => request('GET', '/health')
}

export default api
