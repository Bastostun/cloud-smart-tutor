<template>
  <div class="h-screen bg-gradient-to-br from-[#0a0a0f] via-[#0d1117] to-[#0a0a0f] text-white overflow-hidden flex flex-col">
    
    <!-- 顶部标题栏 -->
    <header class="bg-[#161b22]/80 backdrop-blur-xl border-b border-[#30363d] px-6 py-3 flex items-center justify-between flex-shrink-0">
      <div class="flex items-center space-x-6">
        <div class="flex items-center space-x-3">
          <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-[#58a6ff] to-[#1f6feb] flex items-center justify-center shadow-lg shadow-blue-500/20">
            <svg class="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z"/>
            </svg>
          </div>
          <div>
            <h1 class="text-lg font-bold bg-gradient-to-r from-[#58a6ff] to-[#bc8cff] bg-clip-text text-transparent">云端智师 · 总控大屏</h1>
            <p class="text-[10px] text-gray-500">Cloud Smart Tutor - Command Center v3.0</p>
          </div>
        </div>
        <div class="h-8 w-px bg-[#30363d]"></div>
        <div class="flex items-center space-x-4">
          <div class="flex items-center space-x-2 px-3 py-1.5 bg-[#238636]/10 border border-[#238636]/30 rounded-full">
            <div class="w-2 h-2 rounded-full bg-[#3fb950] animate-pulse"></div>
            <span class="text-[11px] text-[#3fb950] font-medium">系统运行中</span>
          </div>
          <div class="flex items-center space-x-2 text-[11px] text-gray-400">
            <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
            <span>{{ currentTime }}</span>
          </div>
          <div class="flex items-center space-x-2 text-[11px] text-gray-400">
            <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
            </svg>
            <span>运行 {{ uptimeHours }}h</span>
          </div>
        </div>
      </div>
      
      <div class="flex items-center space-x-4">
        <div class="flex items-center space-x-2 px-3 py-1.5 bg-[#1f6feb]/10 border border-[#1f6feb]/30 rounded-full">
          <svg class="w-3 h-3 text-[#58a6ff]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
          </svg>
          <span class="text-[11px] text-[#58a6ff]">{{ onlineStudents }} 名学生在线</span>
        </div>
        <div class="flex items-center space-x-2 px-3 py-1.5 bg-[#d29922]/10 border border-[#d29922]/30 rounded-full">
          <svg class="w-3 h-3 text-[#d29922]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01"/>
          </svg>
          <span class="text-[11px] text-[#d29922]">{{ activeVMs }} 台虚拟机</span>
        </div>
        <button @click="refreshData" class="p-2 rounded-lg hover:bg-[#21262d] transition-colors" :class="{ 'animate-spin': isRefreshing }">
          <svg class="w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
          </svg>
        </button>
      </div>
    </header>

    <!-- 主内容区 -->
    <main class="flex-1 p-4 overflow-hidden">
      <div class="grid grid-cols-12 gap-4 h-full">
        
        <!-- 左侧：系统资源监控 (4列) -->
        <div class="col-span-3 flex flex-col space-y-4">
          <!-- 系统总览卡片 -->
          <div class="bg-[#161b22]/60 backdrop-blur border border-[#30363d] rounded-xl p-4">
            <div class="flex items-center justify-between mb-3">
              <h2 class="text-sm font-semibold text-gray-200 flex items-center space-x-2">
                <svg class="w-4 h-4 text-[#58a6ff]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01"/>
                </svg>
                <span>系统资源总览</span>
              </h2>
              <div class="px-2 py-0.5 bg-[#238636]/20 rounded text-[10px] text-[#3fb950]">218GB 总内存</div>
            </div>
            
            <div class="space-y-3">
              <!-- CPU 使用率 -->
              <div>
                <div class="flex items-center justify-between text-[11px] mb-1.5">
                  <span class="text-gray-400">CPU 使用率</span>
                  <span class="text-white font-mono font-medium">{{ cpuUsage.toFixed(1) }}%</span>
                </div>
                <div class="w-full bg-[#21262d] rounded-full h-2">
                  <div class="h-2 rounded-full transition-all duration-1000 ease-out"
                    :class="cpuUsage > 80 ? 'bg-gradient-to-r from-[#f85149] to-[#da3633]' : cpuUsage > 60 ? 'bg-gradient-to-r from-[#d29922] to-[#bb8009]' : 'bg-gradient-to-r from-[#3fb950] to-[#2ea043]'"
                    :style="{ width: cpuUsage + '%' }"></div>
                </div>
                <div class="flex items-center justify-between text-[10px] text-gray-500 mt-1">
                  <span>16 核 / 32 线程</span>
                  <span>Intel Xeon E5-2686 v4</span>
                </div>
              </div>

              <!-- 内存使用 -->
              <div>
                <div class="flex items-center justify-between text-[11px] mb-1.5">
                  <span class="text-gray-400">内存使用</span>
                  <span class="text-white font-mono font-medium">{{ memoryUsed }}GB / 218GB</span>
                </div>
                <div class="w-full bg-[#21262d] rounded-full h-2">
                  <div class="h-2 rounded-full bg-gradient-to-r from-[#58a6ff] to-[#1f6feb] transition-all duration-1000 ease-out"
                    :style="{ width: memoryPercent + '%' }"></div>
                </div>
                <div class="flex items-center justify-between text-[10px] text-gray-500 mt-1">
                  <span>可用 {{ memoryAvailable }}GB</span>
                  <span>{{ memoryPercent }}%</span>
                </div>
              </div>

              <!-- 存储使用 -->
              <div>
                <div class="flex items-center justify-between text-[11px] mb-1.5">
                  <span class="text-gray-400">ZFS 存储池</span>
                  <span class="text-white font-mono font-medium">{{ storageUsed }}TB / {{ storageTotal }}TB</span>
                </div>
                <div class="w-full bg-[#21262d] rounded-full h-2">
                  <div class="h-2 rounded-full bg-gradient-to-r from-[#bc8cff] to-[#8b5cf6] transition-all duration-1000 ease-out"
                    :style="{ width: storagePercent + '%' }"></div>
                </div>
                <div class="flex items-center justify-between text-[10px] text-gray-500 mt-1">
                  <span>rpool/data</span>
                  <span>{{ storagePercent }}%</span>
                </div>
              </div>
            </div>
          </div>

          <!-- GPU 虚拟化监控 -->
          <div class="bg-[#161b22]/60 backdrop-blur border border-[#30363d] rounded-xl p-4 flex-1">
            <div class="flex items-center justify-between mb-3">
              <h2 class="text-sm font-semibold text-gray-200 flex items-center space-x-2">
                <svg class="w-4 h-4 text-[#f0883e]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
                </svg>
                <span>GPU 虚拟化监控</span>
              </h2>
              <div class="px-2 py-0.5 bg-[#f0883e]/20 rounded text-[10px] text-[#f0883e]">PCIe 直通</div>
            </div>

            <div class="space-y-3">
              <!-- GPU 卡片 -->
              <div class="bg-gradient-to-br from-[#1c1c1f] to-[#242428] rounded-lg p-3 border border-[#30363d]">
                <div class="flex items-center justify-between mb-2">
                  <div class="flex items-center space-x-2">
                    <div class="w-6 h-6 rounded bg-gradient-to-br from-[#76b900] to-[#5a8a00] flex items-center justify-center">
                      <svg class="w-3.5 h-3.5 text-white" fill="currentColor" viewBox="0 0 24 24">
                        <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8z"/>
                      </svg>
                    </div>
                    <div>
                      <div class="text-[11px] font-semibold text-white">NVIDIA RTX A6000</div>
                      <div class="text-[9px] text-gray-500">48GB GDDR6 · PCIe 4.0 x16</div>
                    </div>
                  </div>
                  <div class="w-2 h-2 rounded-full bg-[#3fb950] animate-pulse"></div>
                </div>

                <div class="grid grid-cols-3 gap-2 mb-2">
                  <div class="text-center">
                    <div class="text-lg font-bold text-[#76b900] font-mono">{{ gpuUsage.toFixed(1) }}%</div>
                    <div class="text-[9px] text-gray-500">GPU 使用</div>
                  </div>
                  <div class="text-center">
                    <div class="text-lg font-bold text-[#58a6ff] font-mono">{{ gpuMemoryUsed }}GB</div>
                    <div class="text-[9px] text-gray-500">显存使用</div>
                  </div>
                  <div class="text-center">
                    <div class="text-lg font-bold text-[#f0883e] font-mono">{{ gpuTemp.toFixed(1) }}°C</div>
                    <div class="text-[9px] text-gray-500">温度</div>
                  </div>
                </div>

                <div class="flex items-center justify-between text-[9px] text-gray-500">
                  <span>VFIO 虚拟化: <span class="text-[#3fb950]">已启用</span></span>
                  <span>功率: {{ gpuPower }}W / 300W</span>
                </div>
              </div>

              <!-- GPU 虚拟机分配 -->
              <div class="bg-[#21262d]/50 rounded-lg p-3">
                <div class="text-[11px] font-medium text-gray-300 mb-2">GPU 虚拟机分配</div>
                <div class="space-y-2">
                  <div v-for="vm in gpuVMs" :key="vm.id" class="flex items-center justify-between">
                    <div class="flex items-center space-x-2">
                      <div class="w-1.5 h-1.5 rounded-full" :class="vm.status === 'running' ? 'bg-[#3fb950]' : 'bg-gray-600'"></div>
                      <span class="text-[10px] text-gray-400">{{ vm.name }}</span>
                    </div>
                    <div class="flex items-center space-x-3 text-[10px]">
                      <span class="text-gray-500">vGPU: {{ vm.vgpu }}GB</span>
                      <span class="font-mono text-[#58a6ff]">{{ vm.usage }}%</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 网络流量监控 -->
          <div class="bg-[#161b22]/60 backdrop-blur border border-[#30363d] rounded-xl p-4">
            <div class="flex items-center justify-between mb-3">
              <h2 class="text-sm font-semibold text-gray-200 flex items-center space-x-2">
                <svg class="w-4 h-4 text-[#3fb950]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                </svg>
                <span>网络 I/O</span>
              </h2>
              <div class="flex items-center space-x-2 text-[10px]">
                <span class="text-[#3fb950]">↓ {{ networkDown }} MB/s</span>
                <span class="text-[#58a6ff]">↑ {{ networkUp }} MB/s</span>
              </div>
            </div>
            <div class="h-16 flex items-end space-x-px">
              <div v-for="(bar, i) in networkChart" :key="i" class="flex-1 flex flex-col justify-end space-y-px">
                <div class="bg-[#58a6ff]/40 rounded-t-sm" :style="{ height: (bar.up / 100 * 100) + '%' }"></div>
                <div class="bg-[#3fb950]/40 rounded-t-sm" :style="{ height: (bar.down / 100 * 100) + '%' }"></div>
              </div>
            </div>
          </div>
        </div>

        <!-- 中间：核心区域 (6列) -->
        <div class="col-span-6 flex flex-col space-y-4">
          <!-- AI 部署核心区 -->
          <div class="bg-gradient-to-br from-[#161b22] via-[#1c1c1f] to-[#161b22] border border-[#30363d] rounded-xl p-5 relative overflow-hidden">
            <!-- 背景装饰 -->
            <div class="absolute inset-0 opacity-5">
              <div class="absolute top-0 right-0 w-64 h-64 bg-gradient-to-br from-[#58a6ff] to-[#bc8cff] rounded-full blur-3xl"></div>
              <div class="absolute bottom-0 left-0 w-48 h-48 bg-gradient-to-br from-[#f0883e] to-[#da3633] rounded-full blur-3xl"></div>
            </div>

            <div class="relative">
              <div class="flex items-center justify-between mb-4">
                <div class="flex items-center space-x-3">
                  <div class="w-8 h-8 rounded-lg bg-gradient-to-br from-[#58a6ff] to-[#bc8cff] flex items-center justify-center animate-pulse">
                    <svg class="w-4 h-4 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"/>
                    </svg>
                  </div>
                  <div>
                    <h2 class="text-base font-bold bg-gradient-to-r from-[#58a6ff] to-[#bc8cff] bg-clip-text text-transparent">AI 本地部署环境</h2>
                    <p class="text-[10px] text-gray-500">Local AI Inference Engine · GPU Accelerated</p>
                  </div>
                </div>
                <div class="flex items-center space-x-2">
                  <div class="px-2 py-1 bg-[#238636]/20 border border-[#238636]/30 rounded text-[10px] text-[#3fb950]">运行中</div>
                  <div class="px-2 py-1 bg-[#1f6feb]/20 border border-[#1f6feb]/30 rounded text-[10px] text-[#58a6ff]">v4.0.2</div>
                </div>
              </div>

              <div class="grid grid-cols-3 gap-4 mb-4">
                <!-- DeepSeek 模型状态 -->
                <div class="bg-[#21262d]/60 backdrop-blur rounded-lg p-3 border border-[#30363d]">
                  <div class="flex items-center space-x-2 mb-2">
                    <div class="w-5 h-5 rounded bg-gradient-to-br from-[#1f6feb] to-[#58a6ff] flex items-center justify-center">
                      <span class="text-[8px] font-bold text-white">DS</span>
                    </div>
                    <div>
                      <div class="text-[11px] font-semibold text-white">DeepSeek-V4-Flash</div>
                      <div class="text-[9px] text-gray-500">deepseek-ai 官方模型</div>
                    </div>
                  </div>
                  <div class="space-y-1.5">
                    <div class="flex items-center justify-between text-[10px]">
                      <span class="text-gray-400">模型参数量</span>
                      <span class="text-[#58a6ff] font-mono">236B</span>
                    </div>
                    <div class="flex items-center justify-between text-[10px]">
                      <span class="text-gray-400">上下文窗口</span>
                      <span class="text-[#58a6ff] font-mono">128K</span>
                    </div>
                    <div class="flex items-center justify-between text-[10px]">
                      <span class="text-gray-400">推理延迟</span>
                      <span class="text-[#3fb950] font-mono">{{ aiLatency }}ms</span>
                    </div>
                  </div>
                </div>

                <!-- 推理性能 -->
                <div class="bg-[#21262d]/60 backdrop-blur rounded-lg p-3 border border-[#30363d]">
                  <div class="text-[11px] font-semibold text-white mb-2">推理性能指标</div>
                  <div class="space-y-1.5">
                    <div class="flex items-center justify-between text-[10px]">
                      <span class="text-gray-400">Tokens/s</span>
                      <span class="text-[#76b900] font-mono font-bold">{{ aiTokensPerSec }}</span>
                    </div>
                    <div class="flex items-center justify-between text-[10px]">
                      <span class="text-gray-400">GPU 利用率</span>
                      <span class="text-[#f0883e] font-mono">{{ aiGpuUtil }}%</span>
                    </div>
                    <div class="flex items-center justify-between text-[10px]">
                      <span class="text-gray-400">显存占用</span>
                      <span class="text-[#bc8cff] font-mono">{{ aiVramUsed }}GB</span>
                    </div>
                  </div>
                </div>

                <!-- RAG 知识库 -->
                <div class="bg-[#21262d]/60 backdrop-blur rounded-lg p-3 border border-[#30363d]">
                  <div class="text-[11px] font-semibold text-white mb-2">RAG 知识库</div>
                  <div class="space-y-1.5">
                    <div class="flex items-center justify-between text-[10px]">
                      <span class="text-gray-400">文档向量</span>
                      <span class="text-[#58a6ff] font-mono">{{ ragVectors }}</span>
                    </div>
                    <div class="flex items-center justify-between text-[10px]">
                      <span class="text-gray-400">检索准确率</span>
                      <span class="text-[#3fb950] font-mono">{{ ragAccuracy }}%</span>
                    </div>
                    <div class="flex items-center justify-between text-[10px]">
                      <span class="text-gray-400">平均检索时间</span>
                      <span class="text-[#d29922] font-mono">{{ ragLatency }}ms</span>
                    </div>
                  </div>
                </div>
              </div>

              <!-- AI 功能快捷入口 -->
              <div class="grid grid-cols-4 gap-2">
                <button class="flex flex-col items-center justify-center p-2 bg-[#21262d]/40 hover:bg-[#21262d]/80 rounded-lg border border-[#30363d] hover:border-[#58a6ff]/50 transition-all group">
                  <svg class="w-4 h-4 text-[#58a6ff] mb-1 group-hover:scale-110 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z"/>
                  </svg>
                  <span class="text-[10px] text-gray-400 group-hover:text-white">AI 伴学状态</span>
                </button>
                <button class="flex flex-col items-center justify-center p-2 bg-[#21262d]/40 hover:bg-[#21262d]/80 rounded-lg border border-[#30363d] hover:border-[#bc8cff]/50 transition-all group">
                  <svg class="w-4 h-4 text-[#bc8cff] mb-1 group-hover:scale-110 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
                  </svg>
                  <span class="text-[10px] text-gray-400 group-hover:text-white">课纲知识库</span>
                </button>
                <button class="flex flex-col items-center justify-center p-2 bg-[#21262d]/40 hover:bg-[#21262d]/80 rounded-lg border border-[#30363d] hover:border-[#f0883e]/50 transition-all group">
                  <svg class="w-4 h-4 text-[#f0883e] mb-1 group-hover:scale-110 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
                  </svg>
                  <span class="text-[10px] text-gray-400 group-hover:text-white">自动排错引擎</span>
                </button>
                <button class="flex flex-col items-center justify-center p-2 bg-[#21262d]/40 hover:bg-[#21262d]/80 rounded-lg border border-[#30363d] hover:border-[#3fb950]/50 transition-all group">
                  <svg class="w-4 h-4 text-[#3fb950] mb-1 group-hover:scale-110 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4"/>
                  </svg>
                  <span class="text-[10px] text-gray-400 group-hover:text-white">向量数据库</span>
                </button>
              </div>
            </div>
          </div>

          <!-- 教学错误实时监控 -->
          <div class="bg-[#161b22]/60 backdrop-blur border border-[#30363d] rounded-xl p-4 flex-1 overflow-hidden flex flex-col">
            <div class="flex items-center justify-between mb-3 flex-shrink-0">
              <h2 class="text-sm font-semibold text-gray-200 flex items-center space-x-2">
                <svg class="w-4 h-4 text-[#f85149]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
                </svg>
                <span>高频错误实时监控</span>
              </h2>
              <div class="flex items-center space-x-2">
                <span class="text-[10px] text-gray-500">错误率:</span>
                <span class="text-[11px] font-bold font-mono" :class="errorRate > 20 ? 'text-[#f85149]' : 'text-[#3fb950]'">{{ errorRate.toFixed(1) }}%</span>
              </div>
            </div>

            <!-- 错误排行榜 -->
            <div class="space-y-2 overflow-y-auto flex-1 pr-1">
              <div v-for="(error, index) in topErrors" :key="error.id"
                class="flex items-center justify-between p-2.5 bg-[#21262d]/40 rounded-lg border border-[#30363d] hover:border-[#f85149]/30 transition-colors"
                :class="index === 0 ? 'border-l-2 border-l-[#f85149]' : index === 1 ? 'border-l-2 border-l-[#d29922]' : index === 2 ? 'border-l-2 border-l-[#f0883e]' : ''">
                <div class="flex items-center space-x-3">
                  <div class="w-5 h-5 rounded flex items-center justify-center text-[10px] font-bold"
                    :class="index < 3 ? 'bg-[#f85149]/20 text-[#f85149]' : 'bg-[#21262d] text-gray-500'">
                    {{ index + 1 }}
                  </div>
                  <div>
                    <div class="text-[11px] font-medium text-gray-200">{{ error.type }}</div>
                    <div class="text-[9px] text-gray-500">{{ error.description }}</div>
                  </div>
                </div>
                <div class="flex items-center space-x-3">
                  <div class="w-24 bg-[#21262d] rounded-full h-1.5">
                    <div class="h-1.5 rounded-full transition-all duration-500"
                      :class="error.severity === 'critical' ? 'bg-[#f85149]' : error.severity === 'high' ? 'bg-[#d29922]' : 'bg-[#58a6ff]'"
                      :style="{ width: (error.count / 20 * 100) + '%' }"></div>
                  </div>
                  <div class="text-[11px] font-mono font-medium w-8 text-right"
                    :class="error.severity === 'critical' ? 'text-[#f85149]' : error.severity === 'high' ? 'text-[#d29922]' : 'text-[#58a6ff]'">
                    {{ error.count }}
                  </div>
                </div>
              </div>
            </div>

            <!-- 错误趋势图 -->
            <div class="mt-3 flex-shrink-0">
              <div class="flex items-center justify-between text-[10px] text-gray-500 mb-1">
                <span>错误趋势（30分钟）</span>
                <span>{{ trendData[trendData.length - 1] || 0 }} 次/分</span>
              </div>
              <div class="h-12 flex items-end space-x-px">
                <div v-for="(bar, i) in trendData" :key="i"
                  class="flex-1 bg-gradient-to-t from-[#f85149]/60 to-[#f85149]/20 hover:from-[#f85149]/80 hover:to-[#f85149]/40 transition-all rounded-t-sm"
                  :style="{ height: (bar / 50 * 100) + '%' }"
                  :title="bar + ' 次错误'"></div>
              </div>
            </div>
          </div>
        </div>

        <!-- 右侧：学生状态与实时流 (3列) -->
        <div class="col-span-3 flex flex-col space-y-4">
          <!-- 学生状态概览 -->
          <div class="bg-[#161b22]/60 backdrop-blur border border-[#30363d] rounded-xl p-4">
            <h2 class="text-sm font-semibold text-gray-200 mb-3 flex items-center space-x-2">
              <svg class="w-4 h-4 text-[#58a6ff]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"/>
              </svg>
              <span>学生状态概览</span>
            </h2>

            <div class="grid grid-cols-2 gap-3 mb-3">
              <div class="bg-[#238636]/10 rounded-lg p-2.5 text-center border border-[#238636]/20">
                <div class="text-lg font-bold text-[#3fb950] font-mono">{{ onlineStudents }}</div>
                <div class="text-[9px] text-gray-400">在线学生</div>
              </div>
              <div class="bg-[#1f6feb]/10 rounded-lg p-2.5 text-center border border-[#1f6feb]/20">
                <div class="text-lg font-bold text-[#58a6ff] font-mono">{{ activeVMs }}</div>
                <div class="text-[9px] text-gray-400">活跃虚拟机</div>
              </div>
              <div class="bg-[#d29922]/10 rounded-lg p-2.5 text-center border border-[#d29922]/20">
                <div class="text-lg font-bold text-[#d29922] font-mono">{{ stuckStudents }}</div>
                <div class="text-[9px] text-gray-400">卡壳学生</div>
              </div>
              <div class="bg-[#8b5cf6]/10 rounded-lg p-2.5 text-center border border-[#8b5cf6]/20">
                <div class="text-lg font-bold text-[#bc8cff] font-mono">{{ aiAssistCount }}</div>
                <div class="text-[9px] text-gray-400">AI 介入次数</div>
              </div>
            </div>

            <!-- 卡壳预警 -->
            <div v-if="stuckStudents > 0" class="bg-[#d29922]/10 border border-[#d29922]/30 rounded-lg p-2.5">
              <div class="flex items-center space-x-2 mb-1.5">
                <svg class="w-3 h-3 text-[#d29922] animate-pulse" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
                </svg>
                <span class="text-[11px] font-medium text-[#d29922]">卡壳预警</span>
              </div>
              <div class="text-[10px] text-gray-400">{{ stuckPercentage }}% 学生卡在同类错误</div>
              <div class="text-[10px] text-gray-500 mt-0.5">建议切回大屏集中授课</div>
              <button class="mt-2 w-full py-1.5 text-[10px] bg-[#d29922] hover:bg-[#bb8009] text-white rounded transition-colors">一键切回大屏</button>
            </div>
          </div>

          <!-- 实时错误流 -->
          <div class="bg-[#161b22]/60 backdrop-blur border border-[#30363d] rounded-xl p-4 flex-1 overflow-hidden flex flex-col">
            <h2 class="text-sm font-semibold text-gray-200 mb-3 flex items-center space-x-2 flex-shrink-0">
              <svg class="w-4 h-4 text-[#3fb950]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
              </svg>
              <span>实时错误流</span>
              <div class="w-1.5 h-1.5 rounded-full bg-[#3fb950] animate-pulse ml-auto"></div>
            </h2>

            <div class="space-y-1.5 overflow-y-auto flex-1 pr-1">
              <div v-for="log in errorStream" :key="log.id"
                class="p-2 bg-[#21262d]/30 rounded border-l-2"
                :class="log.level === 'error' ? 'border-l-[#f85149]' : log.level === 'warning' ? 'border-l-[#d29922]' : 'border-l-[#58a6ff]'">
                <div class="flex items-center justify-between mb-0.5">
                  <span class="text-[10px] font-medium text-gray-300">{{ log.student }}</span>
                  <span class="text-[9px] text-gray-600 font-mono">{{ log.time }}</span>
                </div>
                <div class="text-[10px] text-gray-400 truncate">{{ log.message }}</div>
              </div>
            </div>
          </div>

          <!-- AI 教学建议 -->
          <div class="bg-[#161b22]/60 backdrop-blur border border-[#30363d] rounded-xl p-4">
            <h2 class="text-sm font-semibold text-gray-200 mb-3 flex items-center space-x-2">
              <svg class="w-4 h-4 text-[#bc8cff]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"/>
              </svg>
              <span>AI 教学建议</span>
            </h2>
            <div class="space-y-2">
              <div v-for="(suggestion, idx) in aiSuggestions" :key="idx"
                class="p-2.5 bg-[#21262d]/40 rounded-lg border border-[#30363d]">
                <div class="flex items-start space-x-2">
                  <div class="w-5 h-5 rounded flex items-center justify-center flex-shrink-0 mt-0.5"
                    :class="suggestion.type === 'alert' ? 'bg-[#f85149]/20 text-[#f85149]' : suggestion.type === 'info' ? 'bg-[#58a6ff]/20 text-[#58a6ff]' : 'bg-[#3fb950]/20 text-[#3fb950]'">
                    <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path v-if="suggestion.type === 'alert'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01"/>
                      <path v-else-if="suggestion.type === 'info'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                      <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                    </svg>
                  </div>
                  <p class="text-[10px] text-gray-300 leading-relaxed">{{ suggestion.text }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'

// 实时时钟
const currentTime = ref('')
const updateClock = () => {
  const now = new Date()
  currentTime.value = now.toLocaleString('zh-CN', { hour12: false })
}

// 系统运行时间
const uptimeHours = ref(127)

// 刷新状态
const isRefreshing = ref(false)
const refreshData = () => {
  isRefreshing.value = true
  setTimeout(() => isRefreshing.value = false, 1000)
}

// 系统资源
const cpuUsage = ref(42)
const memoryUsed = ref(156)
const memoryTotal = 218
const memoryAvailable = computed(() => (memoryTotal - memoryUsed.value).toFixed(1))
const memoryPercent = computed(() => ((memoryUsed.value / memoryTotal) * 100).toFixed(1))
const storageUsed = ref(3.2)
const storageTotal = ref(8)
const storagePercent = computed(() => ((storageUsed.value / storageTotal.value) * 100).toFixed(1))

// GPU 监控
const gpuUsage = ref(73)
const gpuMemoryUsed = ref(32.4)
const gpuTemp = ref(58)
const gpuPower = ref(187)
const gpuVMs = ref([
  { id: 1, name: 'AI 推理服务 (vm-200)', status: 'running', vgpu: 16, usage: 73 },
  { id: 2, name: 'RAG 向量数据库 (vm-201)', status: 'running', vgpu: 8, usage: 45 },
  { id: 3, name: '错误分析引擎 (vm-202)', status: 'running', vgpu: 12, usage: 62 },
  { id: 4, name: '备用节点 (vm-203)', status: 'stopped', vgpu: 0, usage: 0 }
])

// 网络监控
const networkDown = ref(125)
const networkUp = ref(48)
const networkChart = ref([])

// AI 部署
const aiLatency = ref(89)
const aiTokensPerSec = ref(142)
const aiGpuUtil = ref(73)
const aiVramUsed = ref(32.4)
const ragVectors = ref('24.8K')
const ragAccuracy = ref(94.2)
const ragLatency = ref(12)

// 学生状态
const onlineStudents = ref(52)
const activeVMs = ref(48)
const stuckStudents = ref(7)
const stuckPercentage = ref(28)
const aiAssistCount = ref(156)

// 错误监控
const errorRate = ref(24)
const topErrors = ref([
  { id: 1, type: 'Nginx 配置缺少分号', description: '/etc/nginx/sites-enabled/default 语法错误', count: 15, severity: 'critical' },
  { id: 2, type: '80 端口被占用', description: 'Address already in use: AH00072', count: 9, severity: 'high' },
  { id: 3, type: '权限不足', description: 'Permission denied: 需要 sudo 权限', count: 6, severity: 'medium' },
  { id: 4, type: '服务未安装', description: 'Unit nginx.service not found', count: 4, severity: 'medium' },
  { id: 5, type: '防火墙拦截', description: 'ufw deny 80/tcp', count: 3, severity: 'low' }
])

const trendData = ref([])

const errorStream = ref([
  { id: 1, student: '张同学', time: '14:23:45', level: 'error', message: 'nginx: [emerg] unexpected end of file' },
  { id: 2, student: '李同学', time: '14:23:42', level: 'error', message: 'nginx: [emerg] unexpected end of file' },
  { id: 3, student: '王同学', time: '14:23:38', level: 'warning', message: '端口 80 已被 apache2 占用' },
  { id: 4, student: '刘同学', time: '14:23:35', level: 'error', message: 'nginx: [emerg] unexpected end of file' },
  { id: 5, student: '陈同学', time: '14:23:30', level: 'warning', message: 'Permission denied: mkdir /var/www/html' },
  { id: 6, student: '赵同学', time: '14:23:25', level: 'info', message: 'systemctl restart nginx 成功' },
  { id: 7, student: '孙同学', time: '14:23:20', level: 'error', message: '配置语法检查失败: nginx -t' },
])

const aiSuggestions = ref([
  { type: 'alert', text: '15 名学生在 Nginx 配置上出现同类错误，建议集中讲解。' },
  { type: 'info', text: '端口冲突错误频次上升 40%，可在实验手册中增加排查指南。' },
  { type: 'success', text: '建议将 "Nginx 配置文件语法" 列为课前预习重点。' }
])

// 数据更新
let streamId = 7
const studentNames = ['张同学', '李同学', '王同学', '刘同学', '陈同学', '赵同学', '孙同学', '周同学', '吴同学', '郑同学']
const errorMessages = [
  { level: 'error', message: 'nginx: [emerg] unexpected end of file' },
  { level: 'warning', message: '端口 80 已被占用' },
  { level: 'warning', message: 'Permission denied' },
  { level: 'info', message: '服务启动失败' },
  { level: 'error', message: '配置文件语法错误' },
]

const addErrorStream = () => {
  const student = studentNames[Math.floor(Math.random() * studentNames.length)]
  const error = errorMessages[Math.floor(Math.random() * errorMessages.length)]
  errorStream.value.unshift({ id: ++streamId, student, time: new Date().toTimeString().slice(0, 8), level: error.level, message: error.message })
  if (errorStream.value.length > 30) errorStream.value.pop()
}

const updateMetrics = () => {
  cpuUsage.value = Math.round(Math.min(95, Math.max(20, cpuUsage.value + (Math.random() - 0.5) * 10)) * 10) / 10
  memoryUsed.value = Math.round(Math.min(200, Math.max(120, memoryUsed.value + (Math.random() - 0.5) * 5)) * 10) / 10
  gpuUsage.value = Math.round(Math.min(95, Math.max(40, gpuUsage.value + (Math.random() - 0.5) * 8)) * 10) / 10
  gpuTemp.value = Math.round(Math.min(75, Math.max(45, gpuTemp.value + (Math.random() - 0.5) * 4)) * 10) / 10
  aiLatency.value = Math.round(Math.min(120, Math.max(60, aiLatency.value + (Math.random() - 0.5) * 20)) * 10) / 10
  aiTokensPerSec.value = Math.round(Math.min(180, Math.max(100, aiTokensPerSec.value + (Math.random() - 0.5) * 15)) * 10) / 10
  networkDown.value = Math.round(Math.min(200, Math.max(50, networkDown.value + (Math.random() - 0.5) * 30)) * 10) / 10
  networkUp.value = Math.round(Math.min(80, Math.max(20, networkUp.value + (Math.random() - 0.5) * 15)) * 10) / 10
  onlineStudents.value = Math.min(60, Math.max(40, onlineStudents.value + (Math.random() > 0.7 ? (Math.random() > 0.5 ? 1 : -1) : 0)))
  errorRate.value = Math.round(Math.min(35, Math.max(15, errorRate.value + (Math.random() - 0.5) * 3)) * 10) / 10
  
  // 网络图表更新
  networkChart.value.shift()
  networkChart.value.push({
    down: Math.random() * 80 + 20,
    up: Math.random() * 40 + 10
  })
}

onMounted(() => {
  updateClock()
  setInterval(updateClock, 1000)
  
  // 初始化趋势数据
  for (let i = 0; i < 30; i++) {
    trendData.value.push(10 + Math.random() * 40)
    networkChart.value.push({
      down: Math.random() * 80 + 20,
      up: Math.random() * 40 + 10
    })
  }
  
  // 定时更新
  setInterval(updateMetrics, 2000)
  setInterval(addErrorStream, 3000)
  setInterval(() => {
    trendData.value.shift()
    trendData.value.push(10 + Math.random() * 40)
  }, 3000)
})

onBeforeUnmount(() => {
  // 清理定时器
})
</script>

<style scoped>
/* 自定义滚动条 */
::-webkit-scrollbar {
  width: 4px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 2px;
}

::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.2);
}

/* 动画 */
@keyframes pulse-glow {
  0%, 100% { box-shadow: 0 0 0 0 rgba(88, 166, 255, 0.4); }
  50% { box-shadow: 0 0 0 8px rgba(88, 166, 255, 0); }
}

.animate-pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}
</style>
