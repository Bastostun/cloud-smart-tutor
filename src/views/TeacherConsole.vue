<template>
  <div class="h-screen flex flex-col bg-apple-canvas text-apple-ink overflow-hidden">
    <!-- 标题栏 -->
    <header class="bg-apple-ink text-white h-[48px] flex items-center px-4 flex-shrink-0">
      <!-- Logo & 版本 -->
      <div class="flex items-center space-x-apple-sm mr-apple-lg">
        <svg class="w-5 h-5 text-white" viewBox="0 0 24 24" fill="currentColor">
          <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/>
        </svg>
        <span class="apple-caption text-white/90">云端智师</span>
        <span class="apple-fine-print text-white/60">v8.2.4</span>
      </div>

      <!-- 搜索栏 -->
      <div class="flex-1 max-w-md mx-apple-lg">
        <div class="relative">
          <svg class="absolute left-2 top-1/2 -translate-y-1/2 w-4 h-4 text-white/50" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
          </svg>
          <input v-model="searchQuery" @input="handleSearch"
            class="w-full bg-white/10 text-white apple-caption rounded-apple-sm pl-8 pr-apple-md py-1 outline-none focus:bg-white/20 transition-colors"
            placeholder="搜索虚拟机、容器、节点、存储..." />
          <!-- 搜索结果下拉 -->
          <div v-if="searchResults.length > 0" class="absolute top-full left-0 right-0 mt-1 bg-white rounded-apple-sm shadow-apple-subtle overflow-hidden z-50">
            <div v-for="item in searchResults" :key="item.id"
              @click="selectSearchResult(item)"
              class="px-apple-md py-apple-xs apple-caption hover:bg-apple-parchment cursor-pointer flex items-center justify-between">
              <span>{{ item.name }}</span>
              <span class="apple-badge apple-badge-info text-[10px]">{{ item.type }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 用户信息 -->
      <div class="flex items-center space-x-apple-sm ml-auto">
        <span class="apple-caption text-white/80">{{ user.username }}</span>
        <button @click="showSettings = true" class="w-6 h-6 flex items-center justify-center text-white/60 hover:text-white transition-colors">
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/>
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
          </svg>
        </button>
        <button @click="showHelp = true" class="w-6 h-6 flex items-center justify-center text-white/60 hover:text-white transition-colors">
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
        </button>
        <button @click="showCreateVM = true" class="apple-btn-utility text-white/80 hover:text-white border-white/20 text-[11px]">创建虚拟机</button>
        <button @click="showCreateCT = true" class="apple-btn-utility text-white/80 hover:text-white border-white/20 text-[11px]">创建容器</button>
        <button @click="logout" class="w-6 h-6 flex items-center justify-center text-white/60 hover:text-red-400 transition-colors">
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/>
          </svg>
        </button>
      </div>
    </header>

    <!-- 主内容区 -->
    <div class="flex flex-1 min-h-0">
      <!-- 资源树 -->
      <aside class="w-64 bg-apple-parchment border-r border-apple-hairline flex flex-col flex-shrink-0 overflow-hidden"
        :style="{ width: treeWidth + 'px' }">
        <!-- 视图选择 -->
        <div class="px-apple-sm py-apple-xxs border-b border-apple-divider">
          <select v-model="treeView" class="w-full bg-transparent apple-caption text-apple-ink outline-none">
            <option value="server">服务器视图</option>
            <option value="folder">文件夹视图</option>
            <option value="storage">存储视图</option>
            <option value="pool">资源池视图</option>
          </select>
        </div>

        <!-- 资源树列表 -->
        <div class="flex-1 overflow-y-auto">
          <!-- 数据中心 -->
          <div class="tree-node">
            <div class="tree-node-header flex items-center space-x-apple-xxs px-apple-sm py-apple-xxs cursor-pointer hover:bg-apple-divider"
              :class="{ 'bg-apple-blue/10': selectedItem?.id === 'datacenter' }"
              @click="toggleNode('datacenter'); selectItem('datacenter')">
              <svg class="w-4 h-4 text-apple-blue flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"/>
              </svg>
              <span class="apple-caption-strong text-apple-ink flex-1 truncate">数据中心</span>
              <svg class="w-3 h-3 text-apple-ink-muted-48 flex-shrink-0 transition-transform"
                :class="{ 'rotate-90': expandedNodes.includes('datacenter') }" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
              </svg>
            </div>
            <div v-show="expandedNodes.includes('datacenter')" class="ml-apple-lg">
              <!-- 节点列表 -->
              <div v-for="node in nodes" :key="node.id" class="tree-node">
                <div class="tree-node-header flex items-center space-x-apple-xxs px-apple-sm py-apple-xxs cursor-pointer hover:bg-apple-divider"
                  :class="{ 'bg-apple-blue/10': selectedItem?.id === node.id }"
                  @click="toggleNode(node.id); selectItem(node)">
                  <svg class="w-4 h-4 text-apple-ink-muted-48 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01"/>
                  </svg>
                  <span class="apple-caption text-apple-ink flex-1 truncate">{{ node.name }}</span>
                  <div class="w-1.5 h-1.5 rounded-full flex-shrink-0"
                    :class="node.status === 'online' ? 'bg-green-500' : 'bg-red-500'"></div>
                </div>
                <div v-show="expandedNodes.includes(node.id)" class="ml-apple-lg">
                  <!-- 节点下的 VM -->
                  <div v-for="vm in getVMsByNode(node.name)" :key="vm.vmid"
                    class="tree-node-header flex items-center space-x-apple-xxs px-apple-sm py-apple-xxs cursor-pointer hover:bg-apple-divider"
                    :class="{ 'bg-apple-blue/10': selectedItem?.vmid === vm.vmid }"
                    @click="selectItem(vm)">
                    <svg class="w-4 h-4 flex-shrink-0" :class="vm.type === 'vm' ? 'text-orange-500' : 'text-blue-500'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path v-if="vm.type === 'vm'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
                      <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/>
                    </svg>
                    <span class="apple-caption text-apple-ink flex-1 truncate">{{ vm.name }}</span>
                    <div class="w-1.5 h-1.5 rounded-full flex-shrink-0"
                      :class="vm.status === 'running' ? 'bg-green-500' : vm.status === 'stopped' ? 'bg-orange-500' : 'bg-red-500'"></div>
                  </div>
                </div>
              </div>

              <!-- 存储 -->
              <div v-for="storage in storages" :key="storage.id"
                class="tree-node-header flex items-center space-x-apple-xxs px-apple-sm py-apple-xxs cursor-pointer hover:bg-apple-divider"
                :class="{ 'bg-apple-blue/10': selectedItem?.id === storage.id }"
                @click="selectItem(storage)">
                <svg class="w-4 h-4 text-apple-ink-muted-48 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4"/>
                </svg>
                <span class="apple-caption text-apple-ink flex-1 truncate">{{ storage.name }}</span>
                <div class="w-1.5 h-1.5 rounded-full flex-shrink-0 bg-green-500"></div>
              </div>
            </div>
          </div>
        </div>
      </aside>

      <!-- 内容面板 -->
      <main class="flex-1 flex flex-col min-w-0 overflow-hidden">
        <!-- 操作按钮栏 -->
        <div v-if="selectedItem" class="bg-apple-pearl border-b border-apple-divider px-apple-md py-apple-xxs flex items-center space-x-apple-xs flex-shrink-0">
          <!-- VM 操作按钮 -->
          <template v-if="selectedItem.type === 'vm' || selectedItem.vmid">
            <button @click="startVM" :disabled="selectedItem.status === 'running'"
              class="apple-btn-utility text-apple-caption" :class="selectedItem.status !== 'running' ? 'text-green-600' : 'opacity-50'">
              <span class="flex items-center space-x-apple-xxs">
                <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 24 24"><path d="M8 5v14l11-7z"/></svg>
                <span>启动</span>
              </span>
            </button>
            <div class="relative" ref="stopButtonRef">
              <button @click="showStopOptions = !showStopOptions" :disabled="selectedItem.status !== 'running'"
                class="apple-btn-utility text-apple-caption" :class="selectedItem.status === 'running' ? 'text-orange-600' : 'opacity-50'">
                <span class="flex items-center space-x-apple-xxs">
                  <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 24 24"><path d="M6 6h12v12H6z"/></svg>
                  <span>关机</span>
                </span>
              </button>
              <div v-if="showStopOptions" class="absolute top-full left-0 mt-1 bg-white rounded-apple-sm shadow-apple-subtle overflow-hidden z-40 min-w-32">
                <button @click="shutdownVM" class="w-full text-left px-apple-md py-apple-xxs apple-caption hover:bg-apple-parchment">关机 (ACPI)</button>
                <button @click="stopVM" class="w-full text-left px-apple-md py-apple-xxs apple-caption hover:bg-apple-parchment text-red-600">停止 (强制)</button>
              </div>
            </div>
            <button @click="resetVM" :disabled="selectedItem.status !== 'running'"
              class="apple-btn-utility text-apple-caption" :class="selectedItem.status === 'running' ? 'text-orange-600' : 'opacity-50'">重置</button>
            <button @click="migrateVM" class="apple-btn-utility text-apple-caption">迁移</button>
            <div class="relative" ref="consoleButtonRef">
              <button @click="showConsoleOptions = !showConsoleOptions"
                class="apple-btn-utility text-apple-caption text-apple-blue">
                <span class="flex items-center space-x-apple-xxs">
                  <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
                  </svg>
                  <span>控制台</span>
                </span>
              </button>
              <div v-if="showConsoleOptions" class="absolute top-full left-0 mt-1 bg-white rounded-apple-sm shadow-apple-subtle overflow-hidden z-40 min-w-32">
                <button @click="openConsole('novnc')" class="w-full text-left px-apple-md py-apple-xxs apple-caption hover:bg-apple-parchment">noVNC</button>
                <button @click="openConsole('spice')" class="w-full text-left px-apple-md py-apple-xxs apple-caption hover:bg-apple-parchment">SPICE</button>
                <button @click="openConsole('xterm')" class="w-full text-left px-apple-md py-apple-xxs apple-caption hover:bg-apple-parchment">xterm.js</button>
              </div>
            </div>
          </template>
          <!-- 节点操作按钮 -->
          <template v-else-if="selectedItem.id?.startsWith('node')">
            <button @click="rebootNode" class="apple-btn-utility text-apple-caption text-orange-600">Reboot</button>
            <button @click="shutdownNode" class="apple-btn-utility text-apple-caption text-red-600">Shutdown</button>
            <button @click="openShell" class="apple-btn-utility text-apple-caption">Shell</button>
            <button @click="bulkActions" class="apple-btn-utility text-apple-caption">Bulk Actions</button>
          </template>
        </div>

        <!-- 内容面板标签 -->
        <div v-if="selectedItem" class="flex items-center border-b border-apple-divider px-apple-md flex-shrink-0">
          <button v-for="tab in contentTabs" :key="tab"
            @click="activeContentTab = tab"
            class="px-apple-md py-apple-sm apple-caption-strong border-b-2 transition-colors"
            :class="activeContentTab === tab ? 'border-apple-blue text-apple-blue' : 'border-transparent text-apple-ink-muted-48 hover:text-apple-ink'">
            {{ tab }}
          </button>
        </div>

        <!-- 内容面板主体 -->
        <div class="flex-1 overflow-y-auto p-apple-lg">
          <!-- 加载状态 -->
          <div v-if="isLoading" class="flex items-center justify-center h-full">
            <div class="text-center">
              <svg class="animate-spin w-8 h-8 mx-auto mb-apple-md text-apple-blue" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <p class="apple-caption text-apple-ink-muted-48">加载中...</p>
            </div>
          </div>
          
          <!-- 数据中心内容 -->
          <div v-else-if="selectedItem?.id === 'datacenter'" class="space-y-apple-lg">
            <div v-if="activeContentTab === '概要'" class="space-y-apple-lg">
              <div class="grid grid-cols-2 md:grid-cols-4 gap-apple-lg">
                <div class="apple-card text-center">
                  <div class="apple-headline-md text-apple-blue">{{ summary.nodes }}</div>
                  <div class="apple-caption text-apple-ink-muted-48 mt-apple-xxs">节点</div>
                </div>
                <div class="apple-card text-center">
                  <div class="apple-headline-md text-green-600">{{ summary.vms }}</div>
                  <div class="apple-caption text-apple-ink-muted-48 mt-apple-xxs">虚拟机</div>
                </div>
                <div class="apple-card text-center">
                  <div class="apple-headline-md text-orange-500">{{ summary.containers }}</div>
                  <div class="apple-caption text-apple-ink-muted-48 mt-apple-xxs">容器</div>
                </div>
                <div class="apple-card text-center">
                  <div class="apple-headline-md text-purple-600">{{ summary.storage_pools }}</div>
                  <div class="apple-caption text-apple-ink-muted-48 mt-apple-xxs">存储池</div>
                </div>
              </div>
              <div class="grid grid-cols-1 md:grid-cols-3 gap-apple-lg">
                <div class="apple-card">
                  <h3 class="apple-tagline mb-apple-md">CPU 使用率</h3>
                  <div class="apple-progress h-3">
                    <div class="apple-progress-fill" :style="{ width: summary.cpu_usage_percent + '%' }"></div>
                  </div>
                  <div class="apple-caption text-apple-ink-muted-48 mt-apple-xxs">{{ summary.cpu_usage_percent }}%</div>
                </div>
                <div class="apple-card">
                  <h3 class="apple-tagline mb-apple-md">内存使用率</h3>
                  <div class="apple-progress h-3">
                    <div class="apple-progress-fill warning" :style="{ width: summary.memory_usage_percent + '%' }"></div>
                  </div>
                  <div class="apple-caption text-apple-ink-muted-48 mt-apple-xxs">{{ summary.memory_usage_percent }}%</div>
                </div>
                <div class="apple-card">
                  <h3 class="apple-tagline mb-apple-md">KSM 节省率</h3>
                  <div class="apple-progress h-3">
                    <div class="apple-progress-fill success" :style="{ width: summary.ksm_saving_percent + '%' }"></div>
                  </div>
                  <div class="apple-caption text-apple-ink-muted-48 mt-apple-xxs">{{ summary.ksm_saving_percent }}%</div>
                </div>
              </div>
            </div>
            <div v-else-if="activeContentTab === '存储'" class="apple-card">
              <div class="flex items-center justify-between mb-apple-lg">
                <h2 class="apple-tagline">存储管理</h2>
                <button class="apple-btn-primary">添加存储</button>
              </div>
              <table class="apple-table w-full">
                <thead>
                  <tr>
                    <th class="apple-caption-strong text-left">名称</th>
                    <th class="apple-caption-strong text-left">类型</th>
                    <th class="apple-caption-strong text-right">容量</th>
                    <th class="apple-caption-strong text-right">已用</th>
                    <th class="apple-caption-strong text-center">状态</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="s in storages" :key="s.id" class="border-t border-apple-divider">
                    <td class="py-apple-md apple-caption-strong">{{ s.name }}</td>
                    <td class="apple-caption">{{ s.type }}</td>
                    <td class="apple-caption text-right">{{ s.total_gb }} GB</td>
                    <td class="apple-caption text-right">{{ s.used_gb }} GB</td>
                    <td class="text-center">
                      <span class="apple-badge apple-badge-success">
                        <div class="w-1.5 h-1.5 rounded-full bg-current"></div>
                        <span>活跃</span>
                      </span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div v-else-if="activeContentTab === '权限'" class="apple-card">
              <h2 class="apple-tagline mb-apple-lg">权限管理</h2>
              <div class="space-y-apple-lg">
                <div>
                  <h3 class="apple-caption-strong mb-apple-md">用户</h3>
                  <table class="apple-table w-full">
                    <thead>
                      <tr>
                        <th class="apple-caption-strong text-left">用户名</th>
                        <th class="apple-caption-strong text-left">角色</th>
                        <th class="apple-caption-strong text-left">认证域</th>
                        <th class="apple-caption-strong text-center">状态</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="u in users" :key="u.id" class="border-t border-apple-divider">
                        <td class="py-apple-md apple-caption-strong">{{ u.username }}</td>
                        <td class="apple-caption">{{ u.role }}</td>
                        <td class="apple-caption">{{ u.realm }}</td>
                        <td class="text-center">
                          <span class="apple-badge" :class="u.enabled ? 'apple-badge-success' : 'apple-badge-error'">
                            {{ u.enabled ? '启用' : '禁用' }}
                          </span>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <div>
                  <h3 class="apple-caption-strong mb-apple-md">认证方式</h3>
                  <div class="grid grid-cols-2 gap-apple-md">
                    <div class="apple-card-parchment p-apple-md">
                      <div class="apple-body-strong">Proxmox VE 认证</div>
                      <div class="apple-caption text-apple-ink-muted-48 mt-apple-xxs">本地用户数据库</div>
                    </div>
                    <div class="apple-card-parchment p-apple-md">
                      <div class="apple-body-strong">Linux PAM</div>
                      <div class="apple-caption text-apple-ink-muted-48 mt-apple-xxs">系统用户认证</div>
                    </div>
                    <div class="apple-card-parchment p-apple-md">
                      <div class="apple-body-strong">LDAP</div>
                      <div class="apple-caption text-apple-ink-muted-48 mt-apple-xxs">目录服务认证</div>
                    </div>
                    <div class="apple-card-parchment p-apple-md">
                      <div class="apple-body-strong">微软 AD</div>
                      <div class="apple-caption text-apple-ink-muted-48 mt-apple-xxs">活动目录认证</div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 节点内容 -->
          <div v-else-if="selectedItem?.id?.startsWith('node')" class="space-y-apple-lg">
            <div v-if="activeContentTab === '概要'" class="space-y-apple-lg">
              <div class="grid grid-cols-2 md:grid-cols-3 gap-apple-lg">
                <div class="apple-card">
                  <h3 class="apple-tagline mb-apple-md">CPU</h3>
                  <div class="apple-headline-md">{{ nodeSummary.cpu?.cores }} 核</div>
                  <div class="apple-caption text-apple-ink-muted-48 mt-apple-xxs">{{ nodeSummary.cpu?.model }}</div>
                  <div class="apple-progress mt-apple-md h-2">
                    <div class="apple-progress-fill" :style="{ width: (nodeSummary.cpu?.usage_percent || 0) + '%' }"></div>
                  </div>
                  <div class="apple-caption text-right mt-apple-xxs">{{ nodeSummary.cpu?.usage_percent }}%</div>
                </div>
                <div class="apple-card">
                  <h3 class="apple-tagline mb-apple-md">内存</h3>
                  <div class="apple-headline-md">{{ nodeSummary.memory?.used_gb }} / {{ nodeSummary.memory?.total_gb }} GB</div>
                  <div class="apple-progress mt-apple-md h-2">
                    <div class="apple-progress-fill warning" :style="{ width: (nodeSummary.memory?.usage_percent || 0) + '%' }"></div>
                  </div>
                  <div class="apple-caption text-right mt-apple-xxs">{{ nodeSummary.memory?.usage_percent }}%</div>
                </div>
                <div class="apple-card">
                  <h3 class="apple-tagline mb-apple-md">磁盘</h3>
                  <div class="apple-headline-md">{{ nodeSummary.disk?.used_gb }} / {{ nodeSummary.disk?.total_gb }} GB</div>
                  <div class="apple-progress mt-apple-md h-2">
                    <div class="apple-progress-fill" :style="{ width: (nodeSummary.disk?.usage_percent || 0) + '%' }"></div>
                  </div>
                  <div class="apple-caption text-right mt-apple-xxs">{{ nodeSummary.disk?.usage_percent }}%</div>
                </div>
              </div>
              <div class="apple-card">
                <h3 class="apple-tagline mb-apple-md">节点信息</h3>
                <div class="grid grid-cols-2 gap-apple-md apple-caption">
                  <div><span class="text-apple-ink-muted-48">IP 地址:</span> <span class="font-mono ml-apple-sm">{{ nodeSummary.ip_address }}</span></div>
                  <div><span class="text-apple-ink-muted-48">运行时间:</span> {{ nodeSummary.uptime_hours }} 小时</div>
                  <div><span class="text-apple-ink-muted-48">版本:</span> {{ nodeSummary.version }}</div>
                  <div><span class="text-apple-ink-muted-48">状态:</span> <span class="text-green-600">{{ nodeSummary.status }}</span></div>
                </div>
              </div>
            </div>
          </div>

          <!-- VM 内容 -->
          <div v-else-if="selectedItem?.type === 'vm' || selectedItem?.vmid" class="space-y-apple-lg">
            <div v-if="activeContentTab === '概要'" class="space-y-apple-lg">
              <div class="grid grid-cols-2 md:grid-cols-4 gap-apple-lg">
                <div class="apple-card">
                  <div class="apple-caption text-apple-ink-muted-48">状态</div>
                  <div class="apple-badge mt-apple-sm" :class="{ 'apple-badge-success': selectedItem.status === 'running', 'apple-badge-warning': selectedItem.status === 'stopped' }">
                    <div class="w-1.5 h-1.5 rounded-full bg-current"></div>
                    <span>{{ selectedItem.status === 'running' ? '运行中' : '已停止' }}</span>
                  </div>
                </div>
                <div class="apple-card">
                  <div class="apple-caption text-apple-ink-muted-48">节点</div>
                  <div class="apple-body-strong mt-apple-sm">{{ selectedItem.node }}</div>
                </div>
                <div class="apple-card">
                  <div class="apple-caption text-apple-ink-muted-48">IP 地址</div>
                  <div class="apple-body-strong mt-apple-sm font-mono">{{ selectedItem.ip_address }}</div>
                </div>
                <div class="apple-card">
                  <div class="apple-caption text-apple-ink-muted-48">运行时间</div>
                  <div class="apple-body-strong mt-apple-sm">{{ formatUptime(selectedItem.uptime_seconds) }}</div>
                </div>
              </div>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-apple-lg">
                <div class="apple-card">
                  <h3 class="apple-tagline mb-apple-md">硬件配置</h3>
                  <div class="space-y-apple-sm apple-caption">
                    <div class="flex justify-between"><span class="text-apple-ink-muted-48">CPU</span><span>{{ selectedItem.cpu?.cores || 2 }} 核</span></div>
                    <div class="flex justify-between"><span class="text-apple-ink-muted-48">内存</span><span>{{ selectedItem.memory_mb || selectedItem.memory?.total_mb || 2048 }} MB</span></div>
                    <div class="flex justify-between"><span class="text-apple-ink-muted-48">磁盘</span><span>{{ selectedItem.disk_gb || selectedItem.disks?.[0]?.size_gb || 30 }} GB</span></div>
                    <div class="flex justify-between"><span class="text-apple-ink-muted-48">网络</span><span>virtio, vmbr0</span></div>
                  </div>
                </div>
                <div class="apple-card">
                  <h3 class="apple-tagline mb-apple-md">操作系统</h3>
                  <div class="space-y-apple-sm apple-caption">
                    <div class="flex justify-between"><span class="text-apple-ink-muted-48">镜像</span><span>{{ selectedItem.os_image || selectedItem.os?.image || 'debian-12.0_x64' }}</span></div>
                    <div class="flex justify-between"><span class="text-apple-ink-muted-48">类型</span><span>Linux</span></div>
                  </div>
                </div>
              </div>
            </div>
            <div v-else-if="activeContentTab === '控制台'" class="apple-card">
              <div class="flex items-center justify-between mb-apple-lg">
                <h2 class="apple-tagline">虚拟机控制台</h2>
                <div class="flex items-center space-x-apple-sm">
                  <button @click="openConsole('novnc')" class="apple-btn-primary">noVNC</button>
                  <button @click="openConsole('spice')" class="apple-btn-secondary">SPICE</button>
                </div>
              </div>
              <div class="aspect-video bg-apple-ink rounded-apple-lg flex items-center justify-center">
                <div class="text-center text-white/60">
                  <svg class="w-16 h-16 mx-auto mb-apple-md" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
                  </svg>
                  <p class="apple-tagline mb-apple-sm">点击"noVNC"或"SPICE"按钮打开控制台</p>
                  <p class="apple-caption">或者使用全屏模式以获得更好的体验</p>
                </div>
              </div>
            </div>
            <div v-else-if="activeContentTab === '快照'" class="apple-card">
              <div class="flex items-center justify-between mb-apple-lg">
                <h2 class="apple-tagline">快照管理</h2>
                <button @click="createSnapshot" class="apple-btn-primary">创建快照</button>
              </div>
              <div class="space-y-apple-md">
                <div v-for="snap in snapshots" :key="snap.name" class="apple-card-parchment p-apple-md flex items-center justify-between">
                  <div>
                    <div class="apple-body-strong">{{ snap.name }}</div>
                    <div class="apple-caption text-apple-ink-muted-48 mt-apple-xxs">{{ snap.description }} · {{ snap.date }}</div>
                  </div>
                  <div class="flex items-center space-x-apple-sm">
                    <button @click="rollbackSnapshot(snap.name)" class="apple-btn-utility">回滚</button>
                    <button @click="deleteSnapshot(snap.name)" class="apple-btn-utility text-red-500">删除</button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>

    <!-- 日志面板 -->
    <div class="border-t border-apple-divider flex flex-col flex-shrink-0 transition-all duration-200"
      :class="settings.logHidden ? 'h-0 overflow-hidden' : ''"
      :style="{ height: settings.logHidden ? '0' : logHeight + 'px' }">
      <!-- 日志面板头部 -->
      <div class="bg-apple-pearl px-apple-md py-apple-xxs flex items-center justify-between flex-shrink-0">
        <div class="flex items-center space-x-apple-sm">
          <button v-for="logType in ['任务日志', '集群日志']" :key="logType"
            @click="activeLogType = logType"
            class="apple-fine-print px-apple-sm py-apple-xxs rounded-apple-xs transition-colors"
            :class="activeLogType === logType ? 'bg-apple-blue text-white' : 'text-apple-ink-muted-48 hover:text-apple-ink'">
            {{ logType }}
          </button>
        </div>
        <div class="flex items-center space-x-apple-xs">
          <button @click="logHeight = Math.max(100, Math.min(400, logHeight - 20))" class="w-5 h-5 flex items-center justify-center text-apple-ink-muted-48 hover:text-apple-ink">
            <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7"/></svg>
          </button>
          <button @click="logHeight = Math.max(100, Math.min(400, logHeight + 20))" class="w-5 h-5 flex items-center justify-center text-apple-ink-muted-48 hover:text-apple-ink">
            <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/></svg>
          </button>
          <button @click="settings.logHidden = !settings.logHidden" class="w-5 h-5 flex items-center justify-center text-apple-ink-muted-48 hover:text-apple-ink">
            <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
          </button>
        </div>
      </div>

      <!-- 日志列表 -->
      <div class="flex-1 overflow-y-auto min-h-0">
        <table class="apple-table w-full">
          <thead>
            <tr>
              <th class="apple-fine-print text-left w-48">任务 ID</th>
              <th class="apple-fine-print text-left w-24">类型</th>
              <th class="apple-fine-print text-left w-20">状态</th>
              <th class="apple-fine-print text-left w-24">节点</th>
              <th class="apple-fine-print text-left">描述</th>
              <th class="apple-fine-print text-right w-16">进度</th>
              <th class="apple-fine-print text-right w-20">时间</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="task in tasks" :key="task.id" class="border-t border-apple-divider hover:bg-apple-parchment cursor-pointer"
              @dblclick="showTaskLog(task)">
              <td class="py-apple-xxs apple-fine-print font-mono truncate">{{ task.id }}</td>
              <td class="apple-fine-print">{{ task.type }}</td>
              <td class="apple-fine-print">
                <span class="inline-block w-2 h-2 rounded-full mr-apple-xxs"
                  :class="task.status === 'running' ? 'bg-blue-500 apple-pulse' : task.status === 'completed' ? 'bg-green-500' : 'bg-red-500'"></span>
                {{ task.status === 'running' ? '运行中' : task.status === 'completed' ? '已完成' : '失败' }}
              </td>
              <td class="apple-fine-print">{{ task.node }}</td>
              <td class="apple-fine-print truncate">{{ task.description }}</td>
              <td class="apple-fine-print text-right">
                <span v-if="task.status === 'running'" class="text-apple-blue">{{ task.progress }}%</span>
                <span v-else-if="task.status === 'completed'" class="text-green-600">100%</span>
              </td>
              <td class="apple-fine-print text-right">{{ formatTime(task.started_at) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 底部拖拽条 -->
    <div v-if="!settings.logHidden"
      class="h-1 bg-apple-divider hover:bg-apple-blue/30 cursor-ns-resize transition-colors"
      @mousedown="startResizeLog">
    </div>

    <!-- 设置对话框 -->
    <transition name="apple-fade">
      <div v-if="showSettings" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
        <div class="apple-card w-full max-w-2xl mx-4 max-h-[90vh] overflow-y-auto">
          <div class="flex items-center justify-between mb-apple-lg">
            <h2 class="apple-tagline">我的设置</h2>
            <button @click="showSettings = false" class="apple-btn-utility">关闭</button>
          </div>
          <div class="space-y-apple-xl">
            <!-- 存储设置 -->
            <div>
              <h3 class="apple-caption-strong mb-apple-md">存储设置</h3>
              <div class="space-y-apple-sm">
                <div v-for="s in storages" :key="s.id" class="flex items-center justify-between apple-card-parchment p-apple-md">
                  <div>
                    <div class="apple-body-strong">{{ s.name }}</div>
                    <div class="apple-caption text-apple-ink-muted-48">{{ s.type }} · {{ s.used_gb }}/{{ s.total_gb }} GB</div>
                  </div>
                  <label class="relative inline-flex items-center cursor-pointer">
                    <input type="checkbox" v-model="s.active" class="sr-only peer" />
                    <div class="w-9 h-5 bg-apple-divider peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-apple-blue"></div>
                  </label>
                </div>
              </div>
            </div>

            <!-- xterm.js 设置 -->
            <div>
              <h3 class="apple-caption-strong mb-apple-md">xterm.js 终端设置</h3>
              <div class="grid grid-cols-2 gap-apple-md">
                <div>
                  <label class="apple-caption-strong block mb-apple-xxs">字体</label>
                  <input v-model="settings.xterm.fontFamily" class="apple-input w-full" />
                </div>
                <div>
                  <label class="apple-caption-strong block mb-apple-xxs">字体大小</label>
                  <input v-model.number="settings.xterm.fontSize" type="number" class="apple-input w-full" />
                </div>
                <div>
                  <label class="apple-caption-strong block mb-apple-xxs">字符间距</label>
                  <input v-model.number="settings.xterm.letterSpacing" type="number" step="0.1" class="apple-input w-full" />
                </div>
                <div>
                  <label class="apple-caption-strong block mb-apple-xxs">行高</label>
                  <input v-model.number="settings.xterm.lineHeight" type="number" step="0.1" class="apple-input w-full" />
                </div>
              </div>
            </div>

            <!-- 重置按钮 -->
            <div class="flex items-center space-x-apple-sm">
              <button @click="clearSavedUser" class="apple-btn-secondary">清除保存的用户名</button>
              <button @click="resetLayout" class="apple-btn-secondary">重置布局</button>
            </div>
          </div>
        </div>
      </div>
    </transition>

    <!-- 创建虚拟机对话框 -->
    <transition name="apple-fade">
      <div v-if="showCreateVM" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
        <div class="apple-card w-full max-w-lg mx-4">
          <div class="flex items-center justify-between mb-apple-lg">
            <h2 class="apple-tagline">创建虚拟机</h2>
            <button @click="showCreateVM = false" class="apple-btn-utility">关闭</button>
          </div>
          <div class="space-y-apple-lg">
            <div>
              <label class="apple-caption-strong block mb-apple-xxs">节点</label>
              <select v-model="createVMForm.node" class="apple-input w-full">
                <option v-for="node in nodes" :key="node.id" :value="node.name">{{ node.name }}</option>
              </select>
            </div>
            <div>
              <label class="apple-caption-strong block mb-apple-xxs">虚拟机名称</label>
              <input v-model="createVMForm.name" class="apple-input w-full" placeholder="student-XX" />
            </div>
            <div>
              <label class="apple-caption-strong block mb-apple-xxs">操作系统镜像</label>
              <select v-model="createVMForm.os_image" class="apple-input w-full">
                <option v-for="img in osImages" :key="img.id" :value="img.id">{{ img.name }}（{{ img.description }}）</option>
              </select>
            </div>
            <div class="grid grid-cols-3 gap-apple-md">
              <div>
                <label class="apple-caption-strong block mb-apple-xxs">CPU 核心</label>
                <input v-model.number="createVMForm.cpu_cores" type="number" class="apple-input w-full" />
              </div>
              <div>
                <label class="apple-caption-strong block mb-apple-xxs">内存 (MB)</label>
                <input v-model.number="createVMForm.memory_mb" type="number" class="apple-input w-full" />
              </div>
              <div>
                <label class="apple-caption-strong block mb-apple-xxs">磁盘 (GB)</label>
                <input v-model.number="createVMForm.disk_gb" type="number" class="apple-input w-full" />
              </div>
            </div>
          </div>
          <div class="flex items-center justify-end space-x-apple-sm mt-apple-xl">
            <button @click="showCreateVM = false" class="apple-btn-secondary">取消</button>
            <button @click="createVM" class="apple-btn-primary">创建</button>
          </div>
        </div>
      </div>
    </transition>

    <!-- 创建容器对话框 -->
    <transition name="apple-fade">
      <div v-if="showCreateCT" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
        <div class="apple-card w-full max-w-lg mx-4">
          <div class="flex items-center justify-between mb-apple-lg">
            <h2 class="apple-tagline">创建容器</h2>
            <button @click="showCreateCT = false" class="apple-btn-utility">关闭</button>
          </div>
          <p class="apple-caption text-apple-ink-muted-48">容器创建向导（与虚拟机类似，但使用 LXC 技术）</p>
          <div class="flex items-center justify-end space-x-apple-sm mt-apple-xl">
            <button @click="showCreateCT = false" class="apple-btn-secondary">取消</button>
            <button class="apple-btn-primary">创建</button>
          </div>
        </div>
      </div>
    </transition>

    <!-- 任务日志详情对话框 -->
    <transition name="apple-fade">
      <div v-if="showTaskLogDialog && selectedTask" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
        <div class="apple-card w-full max-w-2xl mx-4 max-h-[80vh]">
          <div class="flex items-center justify-between mb-apple-lg">
            <h2 class="apple-tagline">任务日志 - {{ selectedTask.type }}</h2>
            <div class="flex items-center space-x-apple-sm">
              <button v-if="selectedTask.status === 'running'" @click="stopTask" class="apple-btn-utility text-red-500">中止任务</button>
              <button @click="showTaskLogDialog = false" class="apple-btn-utility">关闭</button>
            </div>
          </div>
          <div class="apple-terminal h-[400px] overflow-y-auto">
            <div v-for="(line, i) in taskLog" :key="i" class="apple-terminal-line apple-terminal-output">{{ line }}</div>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '@/services/api'
import { wsService } from '@/services/websocket'

const user = reactive({
  username: 'admin',
  role: 'admin'
})

const searchQuery = ref('')
const searchResults = ref([])
const selectedItem = ref(null)
const activeContentTab = ref('概要')
const treeView = ref('server')
const expandedNodes = ref(['datacenter'])
const showStopOptions = ref(false)
const showConsoleOptions = ref(false)
const showSettings = ref(false)
const showCreateVM = ref(false)
const showCreateCT = ref(false)
const showTaskLogDialog = ref(false)
const selectedTask = ref(null)
const activeLogType = ref('任务日志')
const isLoading = ref(false)
const errorMsg = ref('')

const treeWidth = ref(250)
const logHeight = ref(200)
const settings = reactive({
  xterm: {
    fontFamily: 'Menlo, Monaco, "Courier New", monospace',
    fontSize: 13,
    letterSpacing: 0,
    lineHeight: 1.2
  },
  logHidden: false
})

const summary = ref({})
const nodes = ref([])
const vms = ref([])
const storages = ref([])
const tasks = ref([])
const users = ref([])
const osImages = ref([])
const snapshots = ref([
  { name: 'pre-update', description: '更新前快照', date: '2026-05-20 10:00:00' },
  { name: 'clean', description: '干净系统', date: '2026-05-15 09:00:00' }
])
const taskLog = ref([])

const nodeSummary = ref({})

const createVMForm = reactive({
  node: '',
  name: '',
  os_image: 'debian-12.0_x64',
  cpu_cores: 2,
  memory_mb: 2048,
  disk_gb: 30
})

const contentTabs = computed(() => {
  if (!selectedItem.value) return []
  if (selectedItem.value.id === 'datacenter') return ['概要', '集群', '存储', '备份', '权限', '防火墙']
  if (selectedItem.value.id?.startsWith('node')) return ['概要', 'Shell', 'System', 'Disks', '任务记录']
  if (selectedItem.value.type === 'vm' || selectedItem.value.vmid) return ['概要', '控制台', '硬件', '选项', '快照', '防火墙', '权限']
  if (selectedItem.value.id) return ['概要', '内容', '权限']
  return ['概要']
})

const getVMsByNode = (nodeName) => {
  return vms.value.filter(vm => vm.node === nodeName)
}

const toggleNode = (nodeId) => {
  const index = expandedNodes.value.indexOf(nodeId)
  if (index > -1) expandedNodes.value.splice(index, 1)
  else expandedNodes.value.push(nodeId)
}

const selectItem = (item) => {
  selectedItem.value = item
  showStopOptions.value = false
  showConsoleOptions.value = false

  if (item.id === 'datacenter') {
    activeContentTab.value = '概要'
  } else if (item.id?.startsWith('node')) {
    activeContentTab.value = '概要'
    loadNodeSummary(item.id)
  } else if (item.type === 'vm' || item.vmid) {
    activeContentTab.value = '概要'
  }
}

const handleSearch = async () => {
  if (searchQuery.value.length < 2) {
    searchResults.value = []
    return
  }
  try {
    const res = await api.datacenter.search(searchQuery.value)
    searchResults.value = res.results || []
  } catch (err) {
    console.error('Search failed:', err)
  }
}

const selectSearchResult = (item) => {
  searchQuery.value = ''
  searchResults.value = []
  const vm = vms.value.find(v => v.vmid === item.id)
  if (vm) selectItem(vm)
  else {
    const node = nodes.value.find(n => n.id === item.id)
    if (node) selectItem(node)
  }
}

const formatUptime = (seconds) => {
  if (!seconds) return '0 秒'
  const days = Math.floor(seconds / 86400)
  const hours = Math.floor((seconds % 86400) / 3600)
  if (days > 0) return `${days} 天 ${hours} 小时`
  return `${hours} 小时`
}

const formatTime = (date) => {
  return new Date(date).toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
}

const loadSummary = async () => {
  try {
    summary.value = await api.datacenter.getSummary()
  } catch (err) { console.error(err) }
}

const loadNodes = async () => {
  try {
    const res = await api.datacenter.getNodes()
    nodes.value = res.nodes || []
  } catch (err) { console.error(err) }
}

const loadVMs = async () => {
  try {
    const res = await api.vms.getList()
    vms.value = res.vms || []
  } catch (err) { console.error(err) }
}

const loadStorage = async () => {
  try {
    const res = await api.storage.getList()
    storages.value = res.storage || []
  } catch (err) { console.error(err) }
}

const loadTasks = async () => {
  try {
    const res = await api.tasks.getRunning()
    tasks.value = res.tasks || []
  } catch (err) { console.error(err) }
}

const loadUsers = async () => {
  try {
    const res = await api.permissions.getUsers()
    users.value = res.users || []
  } catch (err) { console.error(err) }
}

const loadOSImages = async () => {
  try {
    const res = await api.vms.getOSImages()
    osImages.value = res.images || []
  } catch (err) { console.error(err) }
}

const loadNodeSummary = async (nodeId) => {
  try {
    nodeSummary.value = await api.nodes.getSummary(nodeId)
  } catch (err) { console.error(err) }
}

const startVM = async () => {
  if (!selectedItem.value?.vmid) return
  try {
    await api.vms.start(selectedItem.value.vmid)
    await loadVMs()
  } catch (err) { console.error(err) }
}

const shutdownVM = async () => {
  if (!selectedItem.value?.vmid) return
  try {
    await api.vms.stop(selectedItem.value.vmid)
    showStopOptions.value = false
    await loadVMs()
  } catch (err) { console.error(err) }
}

const stopVM = async () => {
  if (!selectedItem.value?.vmid) return
  try {
    await api.vms.stop(selectedItem.value.vmid)
    showStopOptions.value = false
    await loadVMs()
  } catch (err) { console.error(err) }
}

const resetVM = async () => {
  if (!selectedItem.value?.vmid) return
  try {
    await api.vms.reboot(selectedItem.value.vmid)
    await loadVMs()
  } catch (err) { console.error(err) }
}

const migrateVM = () => {
  alert('迁移虚拟机功能')
}

const openConsole = async (type) => {
  if (!selectedItem.value?.vmid) return
  showConsoleOptions.value = false
  try {
    const res = await api.vms.getConsole(selectedItem.value.vmid)
    window.open(res.url, '_blank', 'width=1280,height=720')
  } catch (err) { console.error(err) }
}

const createVM = async () => {
  try {
    const res = await api.vms.create(createVMForm)
    showCreateVM.value = false
    await loadVMs()
    alert('虚拟机创建中，任务 ID: ' + res.task_id)
  } catch (err) { console.error(err) }
}

const createSnapshot = () => {
  alert('创建快照功能')
}

const rollbackSnapshot = (name) => {
  if (confirm('确定要回滚到此快照吗？')) {
    alert('回滚快照: ' + name)
  }
}

const deleteSnapshot = (name) => {
  if (confirm('确定要删除此快照吗？')) {
    snapshots.value = snapshots.value.filter(s => s.name !== name)
  }
}

const showTaskLog = async (task) => {
  selectedTask.value = task
  showTaskLogDialog.value = true
  try {
    const res = await api.tasks.getLog(task.id)
    taskLog.value = res.logs || []
  } catch (err) { console.error(err) }
}

const stopTask = async () => {
  if (!selectedTask.value) return
  try {
    await api.tasks.stop(selectedTask.value.id)
    showTaskLogDialog.value = false
    await loadTasks()
  } catch (err) { console.error(err) }
}

const rebootNode = () => { alert('重启节点') }
const shutdownNode = () => { alert('关闭节点') }
const openShell = () => { alert('打开 Shell') }
const bulkActions = () => { alert('批量操作') }
const clearSavedUser = () => { alert('已清除保存的用户名') }
const resetLayout = () => { treeWidth.value = 250; logHeight.value = 200; settings.logHidden = false }
const logout = async () => {
  try { await api.auth.logout() } catch (err) {}
  localStorage.removeItem('pve_token')
  window.location.href = '/login'
}
const showHelp = () => { alert('帮助文档') }

const startResizeLog = (e) => {
  const startY = e.clientY
  const startHeight = logHeight.value
  const onMouseMove = (e) => {
    logHeight.value = Math.max(100, Math.min(400, startHeight - (e.clientY - startY)))
  }
  const onMouseUp = () => {
    document.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseup', onMouseUp)
  }
  document.addEventListener('mousemove', onMouseMove)
  document.addEventListener('mouseup', onMouseUp)
}

let refreshInterval = null

const handleVMStatusUpdate = (data) => {
  const vmIndex = vms.value.findIndex(v => v.vmid === data.vmid)
  if (vmIndex !== -1) {
    vms.value[vmIndex] = { ...vms.value[vmIndex], ...data }
  }
}

const handleTaskUpdate = (data) => {
  const taskIndex = tasks.value.findIndex(t => t.id === data.task_id)
  if (taskIndex !== -1) {
    tasks.value[taskIndex] = { ...tasks.value[taskIndex], ...data }
  } else {
    tasks.value.unshift(data)
  }
}

const handleNodeUpdate = (data) => {
  const nodeIndex = nodes.value.findIndex(n => n.id === data.node_id)
  if (nodeIndex !== -1) {
    nodes.value[nodeIndex] = { ...nodes.value[nodeIndex], ...data }
  }
}

onMounted(async () => {
  isLoading.value = true
  try {
    await Promise.all([
      loadSummary(),
      loadNodes(),
      loadVMs(),
      loadStorage(),
      loadTasks(),
      loadUsers(),
      loadOSImages()
    ])
  } finally {
    isLoading.value = false
  }

  const token = localStorage.getItem('pve_token')
  if (token) {
    wsService.connect(token)
    wsService.subscribeAll()
    
    wsService.on('vm_status', handleVMStatusUpdate)
    wsService.on('task_status', handleTaskUpdate)
    wsService.on('node_status', handleNodeUpdate)
  }

  refreshInterval = setInterval(() => {
    loadSummary()
    loadTasks()
  }, 5000)
})

onBeforeUnmount(() => {
  if (refreshInterval) clearInterval(refreshInterval)
  wsService.off('vm_status', handleVMStatusUpdate)
  wsService.off('task_status', handleTaskUpdate)
  wsService.off('node_status', handleNodeUpdate)
})
</script>

<style scoped>
.tree-node-header:hover {
  background-color: rgba(0, 102, 204, 0.05);
}
</style>
