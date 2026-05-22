<template>
  <div class="h-screen flex flex-col bg-[#000000] text-white overflow-hidden"
    :class="{ 'cursor-none': isPointerLocked }">
    
    <!-- 标题栏 -->
    <header class="bg-[#1d1d1f] h-[48px] flex items-center px-4 flex-shrink-0 border-b border-[#333]">
      <!-- 左侧：课程名称 + 网络状态 -->
      <div class="flex items-center space-x-3 flex-1 min-w-0">
        <div class="flex items-center space-x-2">
          <svg class="w-5 h-5 text-[#0066cc]" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/>
          </svg>
          <span class="text-[14px] font-semibold truncate">{{ courseName }}</span>
          <span class="text-[12px] text-gray-400 truncate">· {{ labChapter }}</span>
        </div>
        <!-- 网络延迟指示器 -->
        <div class="flex items-center space-x-1 px-2 py-0.5 rounded-full bg-[#272729]">
          <div class="w-1.5 h-1.5 rounded-full"
            :class="wsLatency < 50 ? 'bg-green-500' : wsLatency < 100 ? 'bg-yellow-500' : 'bg-red-500'"></div>
          <span class="text-[10px] text-gray-400 font-mono">{{ wsLatency }}ms</span>
        </div>
      </div>

      <!-- 右侧：功能按钮 -->
      <div class="flex items-center space-x-2">
        <!-- 全局搜索 -->
        <button @click="showSearch = true" class="p-1.5 rounded-lg hover:bg-[#272729] transition-colors" title="全局搜索">
          <svg class="w-4 h-4 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
          </svg>
        </button>

        <!-- 硬件联调视图 -->
        <button @click="showHardwareModal = true" class="p-1.5 rounded-lg hover:bg-[#272729] transition-colors" title="硬件联调视图">
          <svg class="w-4 h-4 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z"/>
          </svg>
        </button>

        <!-- 时光机（快照回滚） -->
        <div class="relative">
          <button @click="showSnapshotDropdown = !showSnapshotDropdown" class="p-1.5 rounded-lg hover:bg-[#272729] transition-colors" title="时光机">
            <svg class="w-4 h-4 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </button>
          <!-- 快照下拉面板 -->
          <div v-if="showSnapshotDropdown" class="absolute right-0 top-full mt-1 w-72 bg-[#1d1d1f] border border-[#333] rounded-lg shadow-lg overflow-hidden z-50">
            <div class="p-3 border-b border-[#333]">
              <div class="text-[12px] font-semibold text-gray-300">快照管理</div>
            </div>
            <div class="max-h-64 overflow-y-auto">
              <div v-for="snap in snapshots" :key="snap.name" class="px-3 py-2 hover:bg-[#272729] cursor-pointer flex items-center justify-between">
                <div>
                  <div class="text-[12px] font-medium">{{ snap.name }}</div>
                  <div class="text-[10px] text-gray-500">{{ snap.description }} · {{ snap.time }}</div>
                </div>
                <button @click.stop="rollbackSnapshot(snap.name)" class="text-[10px] text-[#0066cc] hover:underline">恢复</button>
              </div>
            </div>
            <div class="p-2 border-t border-[#333]">
              <button @click="createSnapshot" class="w-full py-1 text-[11px] text-[#0066cc] hover:underline">创建新快照</button>
            </div>
          </div>
        </div>

        <!-- 全屏专注 -->
        <button @click="toggleFullscreen" class="p-1.5 rounded-lg hover:bg-[#272729] transition-colors" :title="isFullscreen ? '退出全屏' : '全屏专注'">
          <svg class="w-4 h-4 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path v-if="!isFullscreen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8V4m0 0h4M4 16v4m0 0h4m16-12V4m0 0h-4m4 12v4m0 0h-4"/>
            <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4h4V4m0 0v4m12-4h-4V4m0 0v4M4 16v4m0 0h4m-4 0h4m12-4v4m0 0h-4m4 0h-4"/>
          </svg>
        </button>

        <!-- 退出/下课 -->
        <button @click="exitClass" class="px-3 py-1 text-[12px] bg-[#272729] hover:bg-[#333] rounded-lg transition-colors">退出/下课</button>
      </div>
    </header>

    <!-- 主内容区：三面板布局 -->
    <div class="flex flex-1 min-h-0 relative">
      
      <!-- 左侧：实验手册面板 -->
      <aside class="flex flex-col bg-[#f5f5f7] text-[#1d1d1f] flex-shrink-0 overflow-hidden"
        :class="isFullscreen ? 'w-0 overflow-hidden' : ''"
        :style="{ width: isFullscreen ? '0' : leftPanelWidth + 'px' }">
        <!-- 大纲导航 -->
        <div class="px-3 py-2 border-b border-[#e0e0e0] flex-shrink-0">
          <div class="text-[12px] font-semibold mb-2">实验步骤</div>
          <div class="space-y-1">
            <div v-for="(step, idx) in labSteps" :key="idx"
              @click="scrollToStep(idx)"
              class="text-[11px] px-2 py-1 rounded cursor-pointer hover:bg-[#e0e0e0] transition-colors flex items-center"
              :class="activeStep === idx ? 'bg-[#0066cc] text-white' : 'text-[#7a7a7a]'">
              <span class="mr-1.5">{{ idx + 1 }}.</span>
              <span class="truncate">{{ step.title }}</span>
            </div>
          </div>
        </div>

        <!-- Markdown 内容区 -->
        <div class="flex-1 overflow-y-auto px-4 py-3" ref="manualContentRef">
          <div class="prose prose-sm max-w-none" v-html="renderedMarkdown"></div>
        </div>

        <!-- 拖拽分割线 -->
        <div class="absolute left-0 top-0 bottom-0 w-1 cursor-col-resize hover:bg-[#0066cc]/30 transition-colors z-10"
          @mousedown="startResizeLeft"></div>
      </aside>

      <!-- 中间：内容控制台 -->
      <main class="flex-1 flex flex-col min-w-0 bg-[#000000] overflow-hidden">
        <!-- 多标签页 -->
        <div class="flex items-center bg-[#1d1d1f] border-b border-[#333] flex-shrink-0">
          <div class="flex items-center flex-1 overflow-x-auto">
            <div v-for="(tab, idx) in consoleTabs" :key="tab.id"
              @click="activeTab = tab.id"
              class="px-3 py-1.5 text-[11px] cursor-pointer border-r border-[#333] flex items-center space-x-1.5 whitespace-nowrap"
              :class="activeTab === tab.id ? 'bg-[#000000] text-white border-b-2 border-b-[#0066cc]' : 'text-gray-400 hover:bg-[#272729]'">
              <div class="w-1.5 h-1.5 rounded-full"
                :class="tab.status === 'connected' ? 'bg-green-500' : tab.status === 'connecting' ? 'bg-yellow-500 animate-pulse' : 'bg-red-500'"></div>
              <span>{{ tab.label }}</span>
              <span class="text-[9px] text-gray-500">{{ tab.type }}</span>
            </div>
          </div>
          <!-- 断开/重连指示器 -->
          <div v-if="connectionStatus === 'disconnected'" class="px-3 py-1 bg-red-500/10 text-red-400 text-[10px] flex items-center space-x-1">
            <svg class="w-3 h-3 animate-spin" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span>重连中...</span>
          </div>
        </div>

        <!-- 终端容器 -->
        <div class="flex-1 relative overflow-hidden group" ref="terminalContainerRef">
          <!-- Shell 模式 -->
          <div v-if="activeTabType === 'shell'" class="absolute inset-0">
            <div ref="xtermRef" class="w-full h-full"></div>
            <!-- 终端工具栏 -->
            <div class="absolute top-2 right-2 flex items-center space-x-1.5 bg-[#1d1d1f]/90 backdrop-blur rounded-lg px-2 py-1 opacity-0 hover:opacity-100 transition-opacity group-hover:opacity-100 z-20">
              <button @click="pasteFromClipboard" class="p-1 rounded hover:bg-[#272729]" title="粘贴剪贴板内容">
                <svg class="w-3.5 h-3.5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
                </svg>
              </button>
              <button @click="terminal?.clear()" class="p-1 rounded hover:bg-[#272729]" title="清屏">
                <svg class="w-3.5 h-3.5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                </svg>
              </button>
              <button @click="increaseFontSize" class="p-1 rounded hover:bg-[#272729]" title="增大字体">
                <svg class="w-3.5 h-3.5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
                </svg>
              </button>
              <button @click="decreaseFontSize" class="p-1 rounded hover:bg-[#272729]" title="减小字体">
                <svg class="w-3.5 h-3.5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 12H12"/>
                </svg>
              </button>
            </div>
          </div>

          <!-- Desktop 模式 (noVNC) -->
          <div v-else-if="activeTabType === 'desktop'" class="absolute inset-0 bg-black flex items-center justify-center">
            <div class="text-center text-gray-500">
              <svg class="w-16 h-16 mx-auto mb-4 opacity-50" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
              </svg>
              <p class="text-[12px]">noVNC 桌面模式加载中...</p>
            </div>
          </div>

          <!-- 全屏模式退出提示 -->
          <div v-if="isFullscreen" class="absolute bottom-4 left-1/2 -translate-x-1/2 px-3 py-1.5 bg-[#1d1d1f]/80 backdrop-blur rounded-full text-[10px] text-gray-400 opacity-0 hover:opacity-100 transition-opacity">
            按 ESC 退出全屏
          </div>
        </div>
      </main>

      <!-- 右侧：智驭伴学面板 -->
      <aside class="flex flex-col bg-[#1d1d1f] border-l border-[#333] flex-shrink-0 overflow-hidden"
        :class="isFullscreen ? 'w-0 overflow-hidden' : ''"
        :style="{ width: isFullscreen ? '0' : rightPanelWidth + 'px' }">
        <!-- 面板头部 -->
        <div class="px-3 py-2 border-b border-[#333] flex items-center justify-between flex-shrink-0">
          <div class="flex items-center space-x-2">
            <div class="w-2 h-2 rounded-full bg-[#0066cc]"></div>
            <span class="text-[12px] font-semibold">智驭伴学</span>
          </div>
          <div class="flex items-center space-x-1">
            <button @click="clearChat" class="p-1 rounded hover:bg-[#272729]" title="清空对话">
              <svg class="w-3 h-3 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
              </svg>
            </button>
          </div>
        </div>

        <!-- AI 对话流 + 错误预警 -->
        <div class="flex-1 overflow-y-auto px-3 py-2 space-y-3" ref="chatContainerRef">
          <!-- 系统欢迎消息 -->
          <div class="flex items-start space-x-2">
            <div class="w-6 h-6 rounded-full bg-[#0066cc] flex items-center justify-center flex-shrink-0">
              <svg class="w-3.5 h-3.5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"/>
              </svg>
            </div>
            <div class="flex-1">
              <div class="bg-[#272729] rounded-lg px-3 py-2 text-[12px] leading-relaxed">
                <p class="font-medium mb-1">欢迎使用智驭伴学舱 👋</p>
                <p class="text-gray-400">我是你的 AI 实训助手，已加载《{{ courseName }}》课纲知识库。我会实时监控你的终端操作，发现错误时主动介入。</p>
              </div>
            </div>
          </div>

          <!-- 错误预警卡片（静默嗅探触发） -->
          <div v-for="alert in errorAlerts" :key="alert.id" class="flex items-start space-x-2">
            <div class="w-6 h-6 rounded-full bg-red-500/20 flex items-center justify-center flex-shrink-0 mt-0.5">
              <svg class="w-3.5 h-3.5 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
              </svg>
            </div>
            <div class="flex-1">
              <div class="bg-red-500/10 border border-red-500/20 rounded-lg px-3 py-2">
                <div class="text-[11px] font-medium text-red-400 mb-1">🚨 检测到错误</div>
                <div class="text-[10px] font-mono text-gray-400 mb-2 bg-black/30 rounded px-2 py-1">{{ alert.rawLog }}</div>
                <div class="text-[11px] text-gray-300 space-y-1">
                  <p><span class="text-yellow-400">错误定位：</span>{{ alert.location }}</p>
                  <p><span class="text-yellow-400">原理解析：</span>{{ alert.explanation }}</p>
                  <p><span class="text-yellow-400">行动建议：</span>{{ alert.suggestion }}</p>
                </div>
                <!-- 一键修复指令卡 -->
                <div v-if="alert.actionCards?.length" class="mt-2 space-y-1.5">
                  <div v-for="(card, cardIdx) in alert.actionCards" :key="cardIdx" class="bg-[#0a0a0a] rounded px-2.5 py-1.5 flex items-center justify-between">
                    <code class="text-[10px] font-mono text-[#0066cc]">{{ card.command }}</code>
                    <button @click="executeFixCommand(card.command)" class="text-[10px] text-white bg-[#0066cc] hover:bg-[#0071e3] px-2 py-0.5 rounded transition-colors">
                      {{ card.executed ? '已执行 ✓' : '执行修复' }}
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- AI 回复消息 -->
          <div v-for="(msg, idx) in chatMessages" :key="idx" class="flex items-start space-x-2" :class="msg.role === 'user' ? 'flex-row-reverse space-x-reverse' : ''">
            <div v-if="msg.role === 'ai'" class="w-6 h-6 rounded-full bg-[#0066cc] flex items-center justify-center flex-shrink-0">
              <svg class="w-3.5 h-3.5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"/>
              </svg>
            </div>
            <div v-else class="w-6 h-6 rounded-full bg-[#333] flex items-center justify-center flex-shrink-0">
              <span class="text-[10px] font-semibold">我</span>
            </div>
            <div class="flex-1" :class="msg.role === 'user' ? 'text-right' : ''">
              <div class="inline-block max-w-[90%] rounded-lg px-3 py-2 text-[12px] leading-relaxed"
                :class="msg.role === 'user' ? 'bg-[#0066cc] text-white' : 'bg-[#272729] text-gray-200'">
                <div v-html="msg.content"></div>
              </div>
            </div>
          </div>

          <!-- AI 正在输入指示器 -->
          <div v-if="aiTyping" class="flex items-start space-x-2">
            <div class="w-6 h-6 rounded-full bg-[#0066cc] flex items-center justify-center flex-shrink-0">
              <svg class="w-3.5 h-3.5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"/>
              </svg>
            </div>
            <div class="bg-[#272729] rounded-lg px-3 py-2 flex items-center space-x-1">
              <div class="w-1.5 h-1.5 rounded-full bg-gray-400 animate-bounce" style="animation-delay: 0ms"></div>
              <div class="w-1.5 h-1.5 rounded-full bg-gray-400 animate-bounce" style="animation-delay: 150ms"></div>
              <div class="w-1.5 h-1.5 rounded-full bg-gray-400 animate-bounce" style="animation-delay: 300ms"></div>
            </div>
          </div>
        </div>

        <!-- 输入框 -->
        <div class="px-3 py-2 border-t border-[#333] flex-shrink-0">
          <div class="flex items-center space-x-2">
            <input v-model="userInput" @keyup.enter="sendMessage"
              class="flex-1 bg-[#272729] border border-[#333] rounded-lg px-3 py-1.5 text-[12px] outline-none focus:border-[#0066cc] transition-colors placeholder:text-gray-500"
              placeholder="输入问题或指令..." />
            <button @click="sendMessage" :disabled="!userInput.trim() || aiTyping"
              class="p-1.5 rounded-lg bg-[#0066cc] hover:bg-[#0071e3] disabled:opacity-50 disabled:cursor-not-allowed transition-colors">
              <svg class="w-4 h-4 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"/>
              </svg>
            </button>
          </div>
          <div class="text-[9px] text-gray-500 mt-1 text-center">AI 基于《{{ courseName }}》课纲进行增强解析</div>
        </div>

        <!-- 拖拽分割线 -->
        <div class="absolute right-0 top-0 bottom-0 w-1 cursor-col-resize hover:bg-[#0066cc]/30 transition-colors z-10"
          @mousedown="startResizeRight"></div>
      </aside>
    </div>

    <!-- 全局搜索模态框 -->
    <transition name="apple-fade">
      <div v-if="showSearch" class="fixed inset-0 z-50 flex items-start justify-center pt-20 bg-black/60 backdrop-blur-sm" @click.self="showSearch = false">
        <div class="w-full max-w-xl bg-[#1d1d1f] border border-[#333] rounded-xl shadow-2xl overflow-hidden">
          <div class="flex items-center px-4 py-3 border-b border-[#333]">
            <svg class="w-4 h-4 text-gray-400 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
            </svg>
            <input v-model="searchQuery" ref="searchInputRef"
              class="flex-1 bg-transparent text-[14px] outline-none placeholder:text-gray-500"
              placeholder="搜索实验手册、命令字典、AI 排错记录..."
              @input="handleSearch" />
            <button @click="showSearch = false" class="text-[11px] text-gray-400 hover:text-white">ESC</button>
          </div>
          <div class="max-h-80 overflow-y-auto">
            <div v-if="searchResults.length === 0 && searchQuery" class="px-4 py-8 text-center text-[12px] text-gray-500">
              未找到匹配结果
            </div>
            <div v-for="result in searchResults" :key="result.id"
              @click="selectSearchResult(result)"
              class="px-4 py-2.5 hover:bg-[#272729] cursor-pointer flex items-center justify-between">
              <div>
                <div class="text-[12px] font-medium">{{ result.title }}</div>
                <div class="text-[10px] text-gray-500 mt-0.5 truncate">{{ result.preview }}</div>
              </div>
              <span class="text-[9px] px-1.5 py-0.5 rounded bg-[#272729] text-gray-400">{{ result.type }}</span>
            </div>
          </div>
        </div>
      </div>
    </transition>

    <!-- 硬件联调视图模态框 -->
    <transition name="apple-fade">
      <div v-if="showHardwareModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm" @click.self="showHardwareModal = false">
        <div class="w-full max-w-3xl bg-[#1d1d1f] border border-[#333] rounded-xl shadow-2xl overflow-hidden max-h-[80vh]">
          <div class="flex items-center justify-between px-4 py-3 border-b border-[#333]">
            <h3 class="text-[14px] font-semibold">物理资源映射拓扑</h3>
            <button @click="showHardwareModal = false" class="p-1 rounded hover:bg-[#272729]">
              <svg class="w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>
          <div class="p-6 overflow-y-auto">
            <div class="space-y-4">
              <!-- 网络映射 -->
              <div class="bg-[#272729] rounded-lg p-4">
                <h4 class="text-[12px] font-semibold mb-3">网络映射</h4>
                <div class="space-y-2 text-[11px]">
                  <div class="flex items-center justify-between py-1 border-b border-[#333]">
                    <span class="text-gray-400">vnet0</span>
                    <svg class="w-4 h-4 text-[#0066cc]" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3"/></svg>
                    <span class="text-white font-mono">物理交换机 05 号端口</span>
                  </div>
                  <div class="flex items-center justify-between py-1 border-b border-[#333]">
                    <span class="text-gray-400">vnet1</span>
                    <svg class="w-4 h-4 text-[#0066cc]" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3"/></svg>
                    <span class="text-white font-mono">vmbr0 (NAT)</span>
                  </div>
                </div>
              </div>
              <!-- USB/PCIe 直通 -->
              <div class="bg-[#272729] rounded-lg p-4">
                <h4 class="text-[12px] font-semibold mb-3">USB / PCIe 直通</h4>
                <div class="space-y-2 text-[11px]">
                  <div class="flex items-center space-x-2 py-1">
                    <div class="w-1.5 h-1.5 rounded-full bg-green-500"></div>
                    <span class="text-gray-400">USB 加密狗</span>
                    <span class="text-white font-mono">0bda:8153 (已挂载)</span>
                  </div>
                  <div class="flex items-center space-x-2 py-1">
                    <div class="w-1.5 h-1.5 rounded-full bg-gray-500"></div>
                    <span class="text-gray-400">GPU 直通</span>
                    <span class="text-gray-500">未配置</span>
                  </div>
                </div>
              </div>
              <!-- 存储映射 -->
              <div class="bg-[#272729] rounded-lg p-4">
                <h4 class="text-[12px] font-semibold mb-3">ZFS 存储映射</h4>
                <div class="space-y-2 text-[11px]">
                  <div class="flex items-center justify-between py-1 border-b border-[#333]">
                    <span class="text-gray-400">rpool/data/vm-100-disk-0</span>
                    <span class="text-white font-mono">30 GB (linked-clone)</span>
                  </div>
                  <div class="flex items-center justify-between py-1">
                    <span class="text-gray-400">rpool/data/vm-100-disk-1</span>
                    <span class="text-white font-mono">100 GB (独立磁盘)</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </transition>

    <!-- 退出确认对话框 -->
    <transition name="apple-fade">
      <div v-if="showExitDialog" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
        <div class="w-full max-w-sm bg-[#1d1d1f] border border-[#333] rounded-xl shadow-2xl p-6">
          <h3 class="text-[16px] font-semibold mb-2">退出实训</h3>
          <p class="text-[12px] text-gray-400 mb-4">虚拟机将挂起保存当前状态，所有未保存的工作可能会丢失。确定要退出吗？</p>
          <div class="flex items-center justify-end space-x-2">
            <button @click="showExitDialog = false" class="px-4 py-1.5 text-[12px] bg-[#272729] hover:bg-[#333] rounded-lg transition-colors">取消</button>
            <button @click="confirmExit" class="px-4 py-1.5 text-[12px] bg-red-600 hover:bg-red-700 rounded-lg transition-colors">确认退出</button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, nextTick, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import { WebglAddon } from 'xterm-addon-webgl'
import { marked } from 'marked'
import hljs from 'highlight.js/lib/core'
import bash from 'highlight.js/lib/languages/bash'
import 'xterm/css/xterm.css'
import 'highlight.js/styles/github-dark.css'

hljs.registerLanguage('bash', bash)

const route = useRoute()
const router = useRouter()

// 课程信息
const courseName = ref('网络安全攻防技术')
const labChapter = ref('任务 3.2：配置 Nginx 反向代理')
const studentId = computed(() => route.params.studentId)

// 面板宽度
const leftPanelWidth = ref(380)
const rightPanelWidth = ref(360)

// 全屏状态
const isFullscreen = ref(false)

// 网络延迟
const wsLatency = ref(12)

// 实验步骤
const labSteps = ref([
  { title: '环境准备', content: `## 1. 环境准备\n\n首先确认系统环境满足要求：\n\n\`\`\`bash\n# 检查系统版本\ncat /etc/os-release\n\n# 检查网络连接\nping -c 3 10.0.0.1\n\n# 更新软件包列表\nsudo apt update\n\`\`\`\n\n确认所有依赖已安装完成。` },
  { title: '安装 Nginx', content: `## 2. 安装 Nginx\n\n使用包管理器安装 Nginx：\n\n\`\`\`bash\nsudo apt install -y nginx\n\n# 启动 Nginx 服务\nsudo systemctl start nginx\n\n# 设置开机自启\nsudo systemctl enable nginx\n\`\`\`\n\n安装完成后，通过 \`systemctl status nginx\` 验证服务状态。` },
  { title: '配置反向代理', content: `## 3. 配置反向代理\n\n编辑 Nginx 配置文件：\n\n\`\`\`bash\nsudo vim /etc/nginx/sites-available/reverse-proxy\n\`\`\`\n\n添加以下配置：\n\n\`\`\`nginx\nserver {\n    listen 80;\n    server_name proxy.example.com;\n\n    location / {\n        proxy_pass http://10.0.0.5:8080;\n        proxy_set_header Host \\$host;\n        proxy_set_header X-Real-IP \\$remote_addr;\n    }\n}\n\`\`\`\n\n启用配置：\n\n\`\`\`bash\nsudo ln -s /etc/nginx/sites-available/reverse-proxy /etc/nginx/sites-enabled/\nsudo nginx -t\nsudo systemctl reload nginx\n\`\`\`` },
  { title: '测试与验证', content: `## 4. 测试与验证\n\n使用 curl 测试反向代理是否生效：\n\n\`\`\`bash\ncurl -I http://localhost\n\n# 应该看到 200 OK 响应\n# 检查响应头中是否包含后端服务器信息\n\`\`\`\n\n也可以通过浏览器访问 \`http://虚拟机IP\` 进行可视化验证。` }
])

const activeStep = ref(0)
const manualContentRef = ref(null)

// 一键注入代码到终端
const injectCodeToTerminal = (code) => {
  if (terminal && connectionStatus.value === 'connected') {
    terminal.pasteText(code)
    terminal.pasteText('\n')
  }
}

// 从剪贴板粘贴到终端
const pasteFromClipboard = async () => {
  try {
    const text = await navigator.clipboard.readText()
    if (text && terminal) {
      terminal.pasteText(text)
    }
  } catch (err) {
    // 浏览器不支持 clipboard API，使用降级方案
    terminal?.pasteText('')
  }
}

// 终端字体大小控制
const terminalFontSize = ref(13)

const increaseFontSize = () => {
  if (terminalFontSize.value < 20) {
    terminalFontSize.value++
    terminal.options.fontSize = terminalFontSize.value
    nextTick(() => fitAddon?.fit())
  }
}

const decreaseFontSize = () => {
  if (terminalFontSize.value > 10) {
    terminalFontSize.value--
    terminal.options.fontSize = terminalFontSize.value
    nextTick(() => fitAddon?.fit())
  }
}

// Markdown 渲染（带一键注入按钮）
const renderedMarkdown = computed(() => {
  let content = labSteps.value[activeStep.value]?.content || ''
  
  // 先渲染 Markdown
  let html = marked(content, {
    highlight: (code, lang) => {
      if (lang && hljs.getLanguage(lang)) {
        return hljs.highlight(code, { language: lang }).value
      }
      return hljs.highlightAuto(code).value
    }
  })
  
  // 为代码块添加一键注入按钮
  const codeBlockRegex = /<pre><code[^>]*>([\s\S]*?)<\/code><\/pre>/g
  html = html.replace(codeBlockRegex, (match, codeContent) => {
    // 转义 HTML 实体以便存储在 data 属性中
    const escapedCode = codeContent
      .replace(/&amp;/g, '&')
      .replace(/&lt;/g, '<')
      .replace(/&gt;/g, '>')
      .replace(/&quot;/g, '"')
      .replace(/&#39;/g, "'")
      .replace(/\n$/g, '')
    
    return `<div class="code-block-wrapper"><pre><code${match.replace(/<pre><code|<\/code><\/pre>/g, '')}</code></pre><button class="inject-code-btn" data-code="${encodeURIComponent(escapedCode)}" title="一键注入到终端">⚡ 注入</button></div>`
  })
  
  return html
})

// 控制台标签页
const consoleTabs = ref([
  { id: 'vm-100', label: '攻击机 (10.0.0.100)', type: 'shell', status: 'connected' },
  { id: 'vm-101', label: '靶机 (10.0.0.101)', type: 'shell', status: 'connected' },
  { id: 'vm-102', label: 'Web 服务器', type: 'desktop', status: 'connecting' }
])
const activeTab = ref('vm-100')
const activeTabType = computed(() => {
  const tab = consoleTabs.value.find(t => t.id === activeTab.value)
  return tab?.type || 'shell'
})

// 连接状态
const connectionStatus = ref('connected')

// xterm 相关
let terminal = null
let fitAddon = null
const xtermRef = ref(null)
const terminalContainerRef = ref(null)

// 快照
const snapshots = ref([
  { name: 'initial-env', description: '初始环境（系统自动创建）', time: '2026-05-23 08:00:00', isAuto: true },
  { name: 'before-nginx', description: '安装 Nginx 前', time: '2026-05-23 09:15:00', isAuto: false },
  { name: 'config-backup', description: '配置文件备份点', time: '2026-05-23 10:30:00', isAuto: false }
])
const showSnapshotDropdown = ref(false)

// 搜索
const showSearch = ref(false)
const searchQuery = ref('')
const searchResults = ref([])
const searchInputRef = ref(null)

// 硬件模态框
const showHardwareModal = ref(false)

// 退出对话框
const showExitDialog = ref(false)

// AI 聊天
const chatMessages = ref([])
const userInput = ref('')
const aiTyping = ref(false)

// 错误预警
const errorAlerts = ref([])

// 指针锁定
const isPointerLocked = ref(false)

// 拖拽
let isResizingLeft = false
let isResizingRight = false

// 方法
const scrollToStep = (idx) => {
  activeStep.value = idx
  if (manualContentRef.value) {
    manualContentRef.value.scrollTop = 0
  }
}

const initTerminal = () => {
  if (!xtermRef.value) return

  terminal = new Terminal({
    cursorBlink: true,
    fontSize: terminalFontSize.value,
    fontFamily: '"JetBrains Mono", "SF Mono", "Cascadia Code", Menlo, Monaco, "Courier New", monospace',
    theme: {
      background: '#000000',
      foreground: '#ffffff',
      cursor: '#0066cc',
      selection: 'rgba(0, 102, 204, 0.3)',
      black: '#000000',
      red: '#ff6b6b',
      green: '#69db7c',
      yellow: '#ffd43b',
      blue: '#0066cc',
      magenta: '#cc5de8',
      cyan: '#3bc9db',
      white: '#ffffff',
    },
    allowProposedApi: true
  })

  fitAddon = new FitAddon()
  terminal.loadAddon(fitAddon)

  try {
    const webglAddon = new WebglAddon()
    terminal.loadAddon(webglAddon)
    webglAddon.onContextLoss(() => {
      webglAddon.dispose()
    })
  } catch (e) {
    // WebGL 不可用，使用 Canvas 渲染
  }

  terminal.open(xtermRef.value)
  fitAddon.fit()

  terminal.writeln('\x1b[1;32mWelcome to Debian 12.0 (智驭伴学舱 实训环境)\x1b[0m')
  terminal.writeln('\x1b[37mLast login: ' + new Date().toLocaleString() + '\x1b[0m')
  terminal.writeln('')

  terminal.onData((data) => {
    // 发送数据到后端 WebSocket
    if (connectionStatus.value === 'connected') {
      // wsService.sendTermData(activeTab.value, data)
      terminal.write(data)
    }
  })

  // 支持 Ctrl+V 粘贴快捷键
  terminal.attachCustomKeyEventHandler((ev) => {
    if (ev.ctrlKey && ev.key === 'v' && ev.type === 'keydown') {
      navigator.clipboard.readText().then((text) => {
        terminal.pasteText(text)
      })
      return false
    }
    // 支持 Ctrl+C 复制
    if (ev.ctrlKey && ev.key === 'c' && ev.type === 'keydown') {
      const selection = terminal.getSelection()
      if (selection) {
        navigator.clipboard.writeText(selection)
      }
      return false
    }
    return true
  })

  window.addEventListener('resize', () => {
    if (fitAddon && terminal) {
      fitAddon.fit()
    }
  })
}

const toggleFullscreen = () => {
  isFullscreen.value = !isFullscreen.value
  if (isFullscreen.value) {
    document.documentElement.requestFullscreen?.()
  } else {
    document.exitFullscreen?.()
  }
  nextTick(() => {
    if (fitAddon) {
      fitAddon.fit()
    }
  })
}

const createSnapshot = () => {
  const name = 'snapshot-' + Date.now()
  snapshots.value.push({
    name,
    description: '手动创建的快照',
    time: new Date().toLocaleString(),
    isAuto: false
  })
  showSnapshotDropdown.value = false
}

const rollbackSnapshot = (name) => {
  if (confirm(`确定要回滚到快照 "${name}" 吗？当前未保存的更改将丢失。`)) {
    // 执行快照回滚 API
    showSnapshotDropdown.value = false
    terminal?.writeln(`\r\n\x1b[33m正在回滚到快照: ${name}...\x1b[0m\r\n`)
  }
}

const handleSearch = () => {
  if (!searchQuery.value.trim()) {
    searchResults.value = []
    return
  }
  
  const query = searchQuery.value.toLowerCase()
  searchResults.value = [
    ...labSteps.value.flatMap((step, idx) => {
      const matches = step.content.toLowerCase().includes(query)
      return matches ? [{
        id: `step-${idx}`,
        title: step.title,
        preview: step.content.substring(0, 80) + '...',
        type: '实验手册'
      }] : []
    }),
    { id: 'cmd-1', title: 'systemctl restart nginx', preview: '重启 Nginx 服务', type: '命令' },
    { id: 'cmd-2', title: 'nginx -t', preview: '测试 Nginx 配置语法', type: '命令' },
  ].slice(0, 8)
}

const selectSearchResult = (result) => {
  if (result.type === '实验手册') {
    const idx = parseInt(result.id.split('-')[1])
    activeStep.value = idx
  } else if (result.type === '命令') {
    terminal?.pasteText(result.title)
  }
  showSearch.value = false
  searchQuery.value = ''
}

const sendMessage = async () => {
  if (!userInput.value.trim() || aiTyping.value) return

  const message = userInput.value.trim()
  chatMessages.value.push({ role: 'user', content: message })
  userInput.value = ''
  aiTyping.value = true

  // 模拟 AI 响应
  setTimeout(() => {
    aiTyping.value = false
    chatMessages.value.push({
      role: 'ai',
      content: `根据你的实验进度，当前需要配置 Nginx 反向代理。你可以参考左侧实验手册的第 3 步，或者告诉我具体遇到的问题，我会帮你分析。`
    })
    nextTick(() => {
      const container = document.querySelector('.overflow-y-auto')
      if (container) container.scrollTop = container.scrollHeight
    })
  }, 1500)
}

const clearChat = () => {
  chatMessages.value = []
  errorAlerts.value = []
}

const executeFixCommand = (command) => {
  terminal?.pasteText(command + '\n')
  // 标记已执行
  errorAlerts.value.forEach(alert => {
    alert.actionCards?.forEach(card => {
      if (card.command === command) {
        card.executed = true
      }
    })
  })
}

const exitClass = () => {
  showExitDialog.value = true
}

const confirmExit = () => {
  // 挂起虚拟机
  // await api.vms.suspend(vmid)
  terminal?.dispose()
  router.push({ name: 'StudentLogin' })
}

// 面板拖拽调整
const startResizeLeft = (e) => {
  isResizingLeft = true
  const startX = e.clientX
  const startWidth = leftPanelWidth.value
  const onMouseMove = (e) => {
    if (!isResizingLeft) return
    leftPanelWidth.value = Math.max(250, Math.min(600, startWidth + (e.clientX - startX)))
    nextTick(() => fitAddon?.fit())
  }
  const onMouseUp = () => {
    isResizingLeft = false
    document.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseup', onMouseUp)
  }
  document.addEventListener('mousemove', onMouseMove)
  document.addEventListener('mouseup', onMouseUp)
}

const startResizeRight = (e) => {
  isResizingRight = true
  const startX = e.clientX
  const startWidth = rightPanelWidth.value
  const onMouseMove = (e) => {
    if (!isResizingRight) return
    rightPanelWidth.value = Math.max(250, Math.min(600, startWidth - (e.clientX - startX)))
  }
  const onMouseUp = () => {
    isResizingRight = false
    document.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseup', onMouseUp)
  }
  document.addEventListener('mousemove', onMouseMove)
  document.addEventListener('mouseup', onMouseUp)
}

// 模拟错误嗅探
const simulateErrorDetection = () => {
  setTimeout(() => {
    errorAlerts.value.push({
      id: Date.now(),
      rawLog: 'nginx: [emerg] bind() to 0.0.0.0:80 failed (98: Address already in use)',
      location: '/etc/nginx/sites-enabled/default:23',
      explanation: '端口冲突。Apache 或其他服务已经占用了 80 端口，导致 Nginx 无法绑定。',
      suggestion: '停止占用 80 端口的服务，或者将 Nginx 监听端口改为 8080。',
      actionCards: [
        { command: 'sudo systemctl stop apache2', executed: false },
        { command: 'sudo sed -i "s/listen 80/listen 8080/" /etc/nginx/sites-enabled/default && sudo nginx -t && sudo systemctl reload nginx', executed: false }
      ]
    })
  }, 5000)
}

// 全屏状态监听
const handleFullscreenChange = () => {
  isFullscreen.value = !!document.fullscreenElement
  nextTick(() => fitAddon?.fit())
}

// WebSocket 延迟模拟
let latencyInterval = null
const updateLatency = () => {
  wsLatency.value = Math.max(5, Math.min(150, wsLatency.value + (Math.random() - 0.5) * 10))
}

onMounted(async () => {
  await nextTick()
  initTerminal()
  simulateErrorDetection()

  document.addEventListener('fullscreenchange', handleFullscreenChange)
  
  latencyInterval = setInterval(updateLatency, 2000)

  // ESC 退出全屏
  document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape' && isFullscreen.value) {
      toggleFullscreen()
    }
    // Ctrl+K 打开搜索
    if (e.key === 'k' && (e.ctrlKey || e.metaKey)) {
      e.preventDefault()
      showSearch.value = true
      nextTick(() => searchInputRef.value?.focus())
    }
  })

  // 点击外部关闭快照下拉
  document.addEventListener('click', (e) => {
    if (!e.target.closest('.relative')) {
      showSnapshotDropdown.value = false
    }
  })

  // 实验手册内容区域的事件委托 - 一键注入按钮点击
  if (manualContentRef.value) {
    manualContentRef.value.addEventListener('click', (e) => {
      const btn = e.target.closest('.inject-code-btn')
      if (btn) {
        e.preventDefault()
        const encodedCode = btn.getAttribute('data-code')
        if (encodedCode) {
          const code = decodeURIComponent(encodedCode)
          injectCodeToTerminal(code)
        }
      }
    })
  }
})

onBeforeUnmount(() => {
  terminal?.dispose()
  document.removeEventListener('fullscreenchange', handleFullscreenChange)
  if (latencyInterval) clearInterval(latencyInterval)
})
</script>

<style scoped>
.prose h2 {
  font-size: 16px;
  font-weight: 600;
  margin-top: 16px;
  margin-bottom: 8px;
  color: #1d1d1f;
}

.prose p {
  font-size: 13px;
  line-height: 1.6;
  margin-bottom: 12px;
  color: #333;
}

.prose code {
  font-family: 'JetBrains Mono', 'SF Mono', monospace;
  font-size: 11px;
  background: #e0e0e0;
  padding: 2px 6px;
  border-radius: 4px;
  color: #1d1d1f;
}

.prose pre {
  background: #1d1d1f;
  border-radius: 8px;
  padding: 12px;
  margin: 12px 0;
  overflow-x: auto;
}

.prose pre code {
  background: transparent;
  padding: 0;
  color: #ffffff;
  font-size: 11px;
  line-height: 1.6;
}

/* 自定义滚动条 */
::-webkit-scrollbar {
  width: 6px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.15);
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.25);
}

/* 终端容器滚动条 */
.overflow-y-auto::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.15);
}

.overflow-y-auto::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.25);
}

/* 代码块容器 */
.code-block-wrapper {
  position: relative;
  margin: 12px 0;
}

.code-block-wrapper pre {
  margin: 0 !important;
  border-radius: 8px !important;
}

.inject-code-btn {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 10;
  padding: 4px 10px;
  font-size: 10px;
  font-weight: 600;
  color: #ffffff;
  background: rgba(0, 102, 204, 0.85);
  border: none;
  border-radius: 6px;
  cursor: pointer;
  opacity: 0;
  transition: opacity 0.2s ease, background 0.2s ease;
  backdrop-filter: blur(8px);
}

.code-block-wrapper:hover .inject-code-btn {
  opacity: 1;
}

.inject-code-btn:hover {
  background: rgba(0, 113, 227, 1);
}

.inject-code-btn:active {
  transform: scale(0.95);
}

/* xterm 样式覆盖 */
:deep(.xterm) {
  height: 100%;
}

:deep(.xterm-viewport) {
  overflow-y: auto;
}

/* 淡入淡出动画 */
.apple-fade-enter-active,
.apple-fade-leave-active {
  transition: opacity 0.2s ease;
}

.apple-fade-enter-from,
.apple-fade-leave-to {
  opacity: 0;
}
</style>
