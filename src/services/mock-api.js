const delay = (ms = 500) => new Promise(resolve => setTimeout(resolve, ms))

const mockData = {
  user: {
    id: 'admin',
    username: 'admin',
    role: 'admin',
    real_name: '系统管理员',
    email: 'admin@school.local',
    language: 'zh-CN'
  },

  cluster: {
    cluster_name: 'cloud-smart-tutor',
    version: '8.2.4',
    nodes: 3,
    vms: 52,
    containers: 12,
    storage_pools: 4,
    cpu_usage_percent: 45.2,
    memory_usage_percent: 62.8,
    disk_usage_percent: 38.5,
    ksm_saving_percent: 58.3,
    health_status: 'healthy'
  },

  nodes: [
    {
      id: 'node-01',
      name: 'pve-node-01',
      status: 'online',
      cpu_usage: 45.2,
      memory_total_gb: 64,
      memory_used_gb: 42.3,
      disk_total_gb: 2000,
      disk_used_gb: 780,
      uptime_hours: 720,
      version: '8.2.4',
      ip_address: '192.168.1.1'
    },
    {
      id: 'node-02',
      name: 'pve-node-02',
      status: 'online',
      cpu_usage: 38.5,
      memory_total_gb: 64,
      memory_used_gb: 35.2,
      disk_total_gb: 2000,
      disk_used_gb: 650,
      uptime_hours: 720,
      version: '8.2.4',
      ip_address: '192.168.1.2'
    },
    {
      id: 'node-03',
      name: 'pve-node-03',
      status: 'online',
      cpu_usage: 52.1,
      memory_total_gb: 64,
      memory_used_gb: 48.7,
      disk_total_gb: 2000,
      disk_used_gb: 920,
      uptime_hours: 720,
      version: '8.2.4',
      ip_address: '192.168.1.3'
    }
  ],

  vms: [],

  tasks: [
    {
      id: 'task-001',
      type: 'qmstart',
      status: 'running',
      progress: 65,
      node: 'pve-node-01',
      vmid: 100,
      started_at: '2026-05-22T10:00:00Z',
      description: 'Starting VM student-01'
    },
    {
      id: 'task-002',
      type: 'vzcreate',
      status: 'completed',
      progress: 100,
      node: 'pve-node-02',
      vmid: null,
      started_at: '2026-05-22T09:45:00Z',
      completed_at: '2026-05-22T09:52:00Z',
      description: 'Creating container web-server-01'
    },
    {
      id: 'task-003',
      type: 'qmclone',
      status: 'running',
      progress: 42,
      node: 'pve-node-01',
      vmid: 101,
      started_at: '2026-05-22T10:05:00Z',
      description: 'Cloning template to student-02'
    },
    {
      id: 'task-004',
      type: 'qmmigrate',
      status: 'completed',
      progress: 100,
      node: 'pve-node-01',
      vmid: 95,
      started_at: '2026-05-22T09:30:00Z',
      completed_at: '2026-05-22T09:40:00Z',
      description: 'Migrating VM student-95 to pve-node-02'
    }
  ],

  osImages: [
    { id: 'debian-12.0_x64', name: 'Debian 12.0 x64', description: 'Linux 课程定制版', size_gb: 2.1 },
    { id: 'ubuntu-22.04_x64', name: 'Ubuntu 22.04 LTS', description: '标准服务器版', size_gb: 1.8 },
    { id: 'centos-9_x64', name: 'CentOS Stream 9', description: '企业级 Linux', size_gb: 2.3 },
    { id: 'alpine-3.18_x64', name: 'Alpine 3.18', description: '轻量级容器版', size_gb: 0.15 },
    { id: 'ubuntu-20.04_x64', name: 'Ubuntu 20.04 LTS', description: '旧版兼容', size_gb: 1.6 },
    { id: 'windows-11_x64', name: 'Windows 11 Pro', description: 'Windows 开发环境', size_gb: 8.5 }
  ],

  aiModels: [
    { id: 'ollama-llama3', name: 'Ollama Llama 3', description: '本地部署，离线可用', provider: 'ollama', context_length: 8192 },
    { id: 'ollama-qwen', name: 'Ollama Qwen', description: '中文优化，本地部署', provider: 'ollama', context_length: 4096 },
    { id: 'ollama-codellama', name: 'Ollama CodeLlama', description: '代码专用', provider: 'ollama', context_length: 4096 },
    { id: 'custom-api', name: '自定义 API', description: '对接第三方大模型', provider: 'custom', context_length: 0 }
  ],

  firewallRules: [
    { id: 1, name: '允许 SSH', direction: 'in', protocol: 'tcp', port: '22', action: 'accept', source: '0.0.0.0/0' },
    { id: 2, name: '允许 HTTP', direction: 'in', protocol: 'tcp', port: '80', action: 'accept', source: '0.0.0.0/0' },
    { id: 3, name: '允许 HTTPS', direction: 'in', protocol: 'tcp', port: '443', action: 'accept', source: '0.0.0.0/0' },
    { id: 4, name: '阻止 ICMP', direction: 'in', protocol: 'icmp', port: '*', action: 'drop', source: '0.0.0.0/0' },
    { id: 5, name: '允许 DNS', direction: 'out', protocol: 'udp', port: '53', action: 'accept', source: '' }
  ],

  securityGroups: [
    { id: 1, name: 'Web 服务器', description: '允许 HTTP/HTTPS 流量', rules: [
      { id: 1, protocol: 'tcp', port: '80', action: 'allow' },
      { id: 2, protocol: 'tcp', port: '443', action: 'allow' }
    ]},
    { id: 2, name: '数据库', description: '仅允许内网访问', rules: [
      { id: 3, protocol: 'tcp', port: '3306', action: 'allow' },
      { id: 4, protocol: 'tcp', port: '5432', action: 'allow' }
    ]},
    { id: 3, name: '教学环境', description: '教学用开放端口', rules: [
      { id: 5, protocol: 'tcp', port: '8080', action: 'allow' },
      { id: 6, protocol: 'tcp', port: '3000', action: 'allow' },
      { id: 7, protocol: 'tcp', port: '22', action: 'allow' }
    ]}
  ],

  students: [
    { id: 'stu-001', name: '张三', classroom: '班级 01' },
    { id: 'stu-002', name: '李四', classroom: '班级 01' },
    { id: 'stu-003', name: '王五', classroom: '班级 01' },
    { id: 'stu-004', name: '赵六', classroom: '班级 01' },
    { id: 'stu-005', name: '钱七', classroom: '班级 01' },
    { id: 'stu-006', name: '孙八', classroom: '班级 02' },
    { id: 'stu-007', name: '周九', classroom: '班级 02' },
    { id: 'stu-008', name: '吴十', classroom: '班级 02' }
  ],

  settings: {
    xterm: {
      fontFamily: 'Menlo, Monaco, "Courier New", monospace',
      fontSize: 13,
      letterSpacing: 0,
      lineHeight: 1.2
    },
    language: 'zh-CN',
    layout: {
      treeWidth: 250,
      logHeight: 200,
      logHidden: false
    }
  }
}

// 生成 VM 数据
const generateVMs = () => {
  const vms = []
  const statuses = ['running', 'stopped', 'running', 'running', 'running']
  const nodes = ['pve-node-01', 'pve-node-02', 'pve-node-03']

  for (let i = 1; i <= 52; i++) {
    const vmid = 100 + i
    const student = mockData.students[(i - 1) % mockData.students.length]
    const node = nodes[i % 3]
    vms.push({
      vmid,
      name: `student-${String(i).padStart(2, '0')}`,
      type: i <= 40 ? 'vm' : 'ct',
      status: statuses[i % statuses.length],
      node,
      cpu_usage: Math.random() * 80 + 10,
      memory_mb: 2048,
      disk_gb: 30,
      ip_address: `192.168.1.${100 + i}`,
      uptime_seconds: Math.random() * 86400 * 30,
      tags: [`classroom-0${(i % 2) + 1}`],
      student_id: student.id,
      os_image: 'debian-12.0_x64'
    })
  }

  return vms
}

mockData.vms = generateVMs()

// 模拟 API
export const mockApi = {
  async get(path, params = {}) {
    await delay()

    if (path.includes('/auth/verify')) {
      return { valid: true, user: mockData.user }
    }

    if (path.includes('/datacenter/summary')) {
      return mockData.cluster
    }

    if (path.includes('/datacenter/nodes')) {
      return { nodes: mockData.nodes }
    }

    if (path.includes('/datacenter/search')) {
      const q = params.q || ''
      const results = [
        ...mockData.vms.filter(vm => vm.name.toLowerCase().includes(q.toLowerCase())),
        ...mockData.nodes.filter(node => node.name.toLowerCase().includes(q.toLowerCase()))
      ]
      return { results: results.map(r => ({ id: r.vmid || r.id, name: r.name, type: r.type || 'node', status: r.status })) }
    }

    if (path.includes('/vms') && !path.includes('/console') && !path.includes('/options') && !path.includes('/snapshots')) {
      const vms = mockData.vms
      if (params.node) {
        return { vms: vms.filter(vm => vm.node === params.node), total: vms.filter(vm => vm.node === params.node).length }
      }
      return { vms, total: vms.length }
    }

    if (path.match(/\/vms\/\d+$/)) {
      const vmid = parseInt(path.split('/').pop())
      const vm = mockData.vms.find(v => v.vmid === vmid)
      if (!vm) return { error: 'VM not found' }
      return {
        ...vm,
        description: `学生 ${vm.student_id || vm.vmid} 的虚拟机`,
        cpu: { cores: 2, sockets: 1, usage_percent: vm.cpu_usage },
        memory: { total_mb: vm.memory_mb, used_mb: vm.memory_mb * 0.6, usage_percent: 60, balloon_mb: vm.memory_mb },
        disks: [{ name: 'scsi0', size_gb: vm.disk_gb, used_gb: 12.5, type: 'qcow2', storage: 'local-lvm' }],
        networks: [{ name: 'net0', model: 'virtio', bridge: 'vmbr0', mac_address: '32:33:34:35:36:37', ip_address: vm.ip_address }],
        os: { image: 'debian-12.0_x64', type: 'linux' }
      }
    }

    if (path.includes('/vms/') && path.includes('/console')) {
      return {
        url: `wss://192.168.1.1:8006/api2/json/nodes/pve-node-01/qemu/${path.split('/')[3]}/vncwebsocket`,
        ticket: 'PVEVNC:mock-ticket-' + Date.now(),
        port: 5900,
        password: 'xyz123'
      }
    }

    if (path.includes('/nodes/') && path.includes('/summary')) {
      const nodeId = path.split('/')[3]
      const node = mockData.nodes.find(n => n.id === nodeId) || mockData.nodes[0]
      return {
        node: node.name,
        status: node.status,
        cpu: { cores: 16, sockets: 1, usage_percent: node.cpu_usage, model: 'Intel Xeon E5-2680 v4' },
        memory: { total_gb: node.memory_total_gb, used_gb: node.memory_used_gb, free_gb: node.memory_total_gb - node.memory_used_gb, usage_percent: (node.memory_used_gb / node.memory_total_gb * 100).toFixed(1) },
        disk: { total_gb: node.disk_total_gb, used_gb: node.disk_used_gb, free_gb: node.disk_total_gb - node.disk_used_gb, usage_percent: (node.disk_used_gb / node.disk_total_gb * 100).toFixed(1) },
        uptime_hours: node.uptime_hours,
        version: node.version,
        ip_address: node.ip_address
      }
    }

    if (path.includes('/nodes/') && path.includes('/tasks')) {
      return { tasks: mockData.tasks }
    }

    if (path.includes('/tasks/running')) {
      return { tasks: mockData.tasks.filter(t => t.status === 'running') }
    }

    if (path.includes('/tasks/') && path.includes('/log')) {
      return {
        logs: [
          '2026-05-22 10:00:00 Starting VM...',
          '2026-05-22 10:00:05 Allocating resources...',
          '2026-05-22 10:00:10 Loading disk image...',
          '2026-05-22 10:00:15 Starting network...',
          '2026-05-22 10:00:20 VM started successfully.'
        ]
      }
    }

    if (path.includes('/storage')) {
      return {
        storage: [
          { id: 'local', name: 'local', type: 'dir', total_gb: 100, used_gb: 45, active: true, content: 'images,backup,iso' },
          { id: 'local-lvm', name: 'local-lvm', type: 'lvmthin', total_gb: 500, used_gb: 320, active: true, content: 'images,rootdir' },
          { id: 'nfs-backup', name: 'nfs-backup', type: 'nfs', total_gb: 2000, used_gb: 850, active: true, content: 'backup,vztmpl,iso' }
        ]
      }
    }

    if (path.includes('/vm/os-images')) {
      return { images: mockData.osImages }
    }

    if (path.includes('/vm/ai-models')) {
      return { models: mockData.aiModels }
    }

    if (path.includes('/firewall/rules')) {
      return { rules: mockData.firewallRules }
    }

    if (path.includes('/settings')) {
      return mockData.settings
    }

    if (path.includes('/users')) {
      return {
        users: [
          { id: 'admin', username: 'admin', role: 'admin', realm: 'pve', email: 'admin@school.local', enabled: true },
          { id: 'teacher01', username: 'teacher01', role: 'teacher', realm: 'pve', email: 'teacher01@school.local', enabled: true },
          { id: 'stu-001', username: 'stu-001', role: 'student', realm: 'pve', enabled: true }
        ]
      }
    }

    if (path.includes('/groups')) {
      return {
        groups: [
          { id: 'admin', name: 'Administrators', comment: '系统管理员组', users: ['admin'] },
          { id: 'teachers', name: 'Teachers', comment: '教师组', users: ['teacher01'] },
          { id: 'students', name: 'Students', comment: '学生组', users: ['stu-001'] }
        ]
      }
    }

    if (path.includes('/roles')) {
      return {
        roles: [
          { id: 'Administrator', privileges: 'Everything' },
          { id: 'PVEAdmin', privileges: 'System.Console, System.Pool.Allocate' },
          { id: 'PVEAuditor', privileges: 'Sys.Audit, Sys.Modify, VM.Audit' },
          { id: 'PVEVMUser', privileges: 'VM.Allocate, VM.Audit, VM.Console' }
        ]
      }
    }

    if (path.includes('/health')) {
      return { status: 'healthy', service: 'cloud-smart-tutor-backend', version: '1.0.0' }
    }

    return {}
  },

  async post(path, data = {}) {
    await delay(800)

    if (path.includes('/auth/login')) {
      if (data.username && data.password) {
        return {
          token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.mock',
          user: { ...mockData.user, username: data.username },
          expires_at: new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString()
        }
      }
      return { error: 'Invalid credentials' }
    }

    if (path.includes('/auth/logout')) {
      return { status: 'logged_out' }
    }

    if (path.includes('/vms') && !path.includes('/console')) {
      const vmid = data.vmid || (100 + mockData.vms.length + 1)
      return {
        task_id: 'UPID:pve-node-01:00012345:00012346:663E1234:vzcreate:' + vmid + ':admin@pve:',
        vmid,
        status: 'creating'
      }
    }

    if (path.includes('/vms/') && (path.includes('/start') || path.includes('/stop') || path.includes('/reboot') || path.includes('/suspend') || path.includes('/resume'))) {
      const vmid = path.split('/')[3]
      const action = path.split('/').pop()
      return {
        task_id: `UPID:pve-node-01:00012345:00012346:663E1234:qm${action}:${vmid}:admin@pve:`,
        status: action + 'ing'
      }
    }

    if (path.includes('/vms/') && path.includes('/clone')) {
      return {
        task_id: 'UPID:pve-node-01:00012345:00012346:663E1234:qmclone:' + data.newid + ':admin@pve:',
        status: 'cloning'
      }
    }

    if (path.includes('/vms/') && path.includes('/migrate')) {
      return {
        task_id: 'UPID:pve-node-01:00012345:00012346:663E1234:qmmigrate:100:admin@pve:',
        status: 'migrating'
      }
    }

    if (path.includes('/tasks/') && path.includes('/stop')) {
      return { status: 'stopping' }
    }

    return { status: 'success' }
  },

  async delete(path) {
    await delay(600)

    if (path.includes('/vms/')) {
      const vmid = path.split('/').pop()
      return {
        task_id: 'UPID:pve-node-01:00012345:00012346:663E1234:qmdestroy:' + vmid + ':admin@pve:',
        status: 'deleting'
      }
    }

    return { status: 'success' }
  },

  async put(path, data = {}) {
    await delay(500)

    if (path.includes('/vms/') && path.includes('/options')) {
      return { status: 'updated' }
    }

    return { status: 'success' }
  }
}

export default mockApi
