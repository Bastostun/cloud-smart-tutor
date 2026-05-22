/**
 * AI Skills Core - 智驭伴学舱 AI 技能核心模块
 * 
 * 本模块整合所有 AI 技能相关的功能模块、算法实现及配置
 * 包括 RAG 检索增强生成、错误诊断、代码注入、智能教学建议等
 * 
 * @module ai_skills_core
 * @version 4.0.2
 * @author 云端智师团队
 */

// ==================== 配置常量 ====================

const AI_CONFIG = {
  // DeepSeek 模型配置
  model: {
    name: 'deepseek-ai/DeepSeek-V4-Flash',
    parameters: '236B',
    contextWindow: 128 * 1024, // 128K tokens
    maxTokens: 4096,
    temperature: 0.7,
    topP: 0.9,
    frequencyPenalty: 0.1,
    presencePenalty: 0.1,
  },
  
  // RAG 知识库配置
  rag: {
    vectorStore: 'faiss',
    embeddingModel: 'text-embedding-ada-002',
    chunkSize: 512,
    chunkOverlap: 50,
    topK: 5,
    similarityThreshold: 0.75,
  },
  
  // 推理引擎配置
  inference: {
    backend: 'vllm',
    gpuMemoryUtilization: 0.9,
    maxNumSeqs: 256,
    maxModelLen: 16384,
    dtype: 'float16',
    quantization: 'awq',
  },
  
  // 错误诊断配置
  diagnostics: {
    maxErrorHistory: 100,
    debounceInterval: 2000,
    severityLevels: ['info', 'warning', 'error', 'critical'],
    autoFixThreshold: 0.85,
  },
  
  // 教学建议配置
  teaching: {
    minStuckThreshold: 0.28,
    suggestInterval: 30000,
    maxSuggestionsPerSession: 10,
  },
};

// ==================== AI 技能模块 ====================

/**
 * RAG 检索增强生成模块
 */
class RAGModule {
  constructor(config) {
    this.config = config.rag;
    this.vectorStore = null;
    this.documents = [];
    this.isInitialized = false;
  }

  /**
   * 初始化 RAG 模块
   */
  async initialize() {
    console.log('[RAG] Initializing vector store...');
    this.vectorStore = await this._createVectorStore();
    this.isInitialized = true;
    console.log('[RAG] Vector store ready');
  }

  /**
   * 添加文档到知识库
   */
  async addDocument(doc) {
    const chunks = this._chunkDocument(doc);
    const embeddings = await this._generateEmbeddings(chunks);
    
    for (let i = 0; i < chunks.length; i++) {
      this.vectorStore.add({
        id: `${doc.id}_chunk_${i}`,
        text: chunks[i],
        embedding: embeddings[i],
        metadata: { docId: doc.id, source: doc.source },
      });
    }
    
    this.documents.push(doc);
    return chunks.length;
  }

  /**
   * 检索相关文档
   */
  async retrieve(query, options = {}) {
    if (!this.isInitialized) {
      throw new Error('RAG module not initialized');
    }

    const queryEmbedding = await this._generateEmbedding(query);
    const results = this.vectorStore.search(queryEmbedding, {
      topK: options.topK || this.config.topK,
      threshold: options.threshold || this.config.similarityThreshold,
    });

    return results.map(r => ({
      text: r.text,
      score: r.score,
      metadata: r.metadata,
    }));
  }

  /**
   * 生成回答
   */
  async generateAnswer(query, context) {
    const prompt = this._buildPrompt(query, context);
    const response = await this._callModel(prompt);
    
    return {
      answer: response.text,
      sources: context.map(c => c.metadata),
      confidence: response.confidence,
    };
  }

  // 私有方法
  async _createVectorStore() {
    return {
      add: (item) => console.log(`[VectorStore] Added item: ${item.id}`),
      search: (query, opts) => [
        { text: '示例检索结果 1', score: 0.92, metadata: { source: 'manual' } },
        { text: '示例检索结果 2', score: 0.87, metadata: { source: 'lab_guide' } },
      ],
    };
  }

  async _generateEmbedding(text) {
    return new Array(1536).fill(0).map(() => Math.random() * 2 - 1);
  }

  async _generateEmbeddings(chunks) {
    return Promise.all(chunks.map(c => this._generateEmbedding(c)));
  }

  _chunkDocument(doc) {
    const words = doc.content.split(' ');
    const chunks = [];
    for (let i = 0; i < words.length; i += this.config.chunkSize) {
      chunks.push(words.slice(i, i + this.config.chunkSize).join(' '));
    }
    return chunks;
  }

  _buildPrompt(query, context) {
    const contextText = context.map(c => c.text).join('\n\n');
    return `基于以下上下文信息，回答问题：\n\n上下文：\n${contextText}\n\n问题：${query}\n\n请给出详细解答：`;
  }

  async _callModel(prompt) {
    return {
      text: `基于课纲知识库分析，${prompt.substring(0, 50)}...`,
      confidence: 0.89,
    };
  }
}

/**
 * 错误诊断模块
 */
class DiagnosticModule {
  constructor(config) {
    this.config = config.diagnostics;
    this.errorHistory = [];
    this.patterns = this._loadPatterns();
    this.isListening = false;
  }

  /**
   * 开始监听终端输出
   */
  startListening(stream) {
    this.isListening = true;
    stream.on('data', (data) => this._processOutput(data));
    stream.on('error', (error) => this._processError(error));
    console.log('[Diagnostic] Started listening to terminal output');
  }

  /**
   * 停止监听
   */
  stopListening() {
    this.isListening = false;
    console.log('[Diagnostic] Stopped listening');
  }

  /**
   * 分析错误并生成修复建议
   */
  async analyzeError(errorLog) {
    const pattern = this._matchPattern(errorLog);
    
    if (!pattern) {
      return {
        type: 'unknown',
        severity: 'warning',
        message: '未识别的错误类型',
        suggestions: ['请检查命令语法'],
      };
    }

    return {
      type: pattern.type,
      severity: pattern.severity,
      location: pattern.location,
      explanation: pattern.explanation,
      suggestions: pattern.solutions,
      actionCards: this._generateActionCards(pattern),
    };
  }

  /**
   * 获取错误统计
   */
  getErrorStats() {
    const stats = {
      total: this.errorHistory.length,
      bySeverity: {},
      byType: {},
      topErrors: [],
    };

    this.errorHistory.forEach(err => {
      stats.bySeverity[err.severity] = (stats.bySeverity[err.severity] || 0) + 1;
      stats.byType[err.type] = (stats.byType[err.type] || 0) + 1;
    });

    stats.topErrors = Object.entries(stats.byType)
      .sort((a, b) => b[1] - a[1])
      .slice(0, 5)
      .map(([type, count]) => ({ type, count }));

    return stats;
  }

  // 私有方法
  _loadPatterns() {
    return [
      {
        pattern: /nginx:.*emerg.*bind.*failed.*Address already in use/i,
        type: 'nginx_port_conflict',
        severity: 'critical',
        location: '/etc/nginx/sites-enabled/',
        explanation: '端口冲突。Apache 或其他服务已经占用了 80 端口，导致 Nginx 无法绑定。',
        solutions: [
          '停止占用 80 端口的服务: sudo systemctl stop apache2',
          '或者将 Nginx 监听端口改为 8080',
        ],
      },
      {
        pattern: /nginx:.*emerg.*unexpected end of file/i,
        type: 'nginx_syntax_error',
        severity: 'critical',
        location: '/etc/nginx/sites-enabled/',
        explanation: 'Nginx 配置文件语法错误，通常是因为缺少分号或大括号不匹配。',
        solutions: [
          '检查配置文件语法: sudo nginx -t',
          '编辑配置文件修复语法错误',
        ],
      },
      {
        pattern: /Permission denied/i,
        type: 'permission_denied',
        severity: 'warning',
        location: 'current_directory',
        explanation: '权限不足。当前用户没有执行该操作的权限。',
        solutions: [
          '使用 sudo 提升权限',
          '检查文件权限设置',
        ],
      },
      {
        pattern: /command not found/i,
        type: 'command_not_found',
        severity: 'warning',
        location: 'shell',
        explanation: '命令未找到。请检查命令拼写或是否已安装该程序。',
        solutions: [
          '检查命令拼写',
          '使用 apt install 安装对应软件包',
        ],
      },
      {
        pattern: /systemctl.*failed/i,
        type: 'service_failed',
        severity: 'error',
        location: 'systemd',
        explanation: '服务启动失败。请查看服务日志获取详细信息。',
        solutions: [
          '查看服务状态: sudo systemctl status <service>',
          '查看服务日志: sudo journalctl -u <service> -f',
        ],
      },
    ];
  }

  _matchPattern(errorLog) {
    return this.patterns.find(p => p.pattern.test(errorLog));
  }

  _processOutput(data) {
    const lines = data.toString().split('\n');
    lines.forEach(line => {
      if (this._isErrorLine(line)) {
        this.errorHistory.push({
          timestamp: Date.now(),
          log: line,
          type: 'error',
        });
      }
    });
  }

  _processError(error) {
    this.errorHistory.push({
      timestamp: Date.now(),
      log: error.message,
      type: 'system_error',
    });
  }

  _isErrorLine(line) {
    return /error|failed|critical|fatal|emerg/i.test(line);
  }

  _generateActionCards(pattern) {
    return pattern.solutions.map(s => ({
      command: s.split(':')[1]?.trim() || s,
      description: s,
    }));
  }
}

/**
 * 智能教学建议模块
 */
class TeachingSuggestionModule {
  constructor(config) {
    this.config = config.teaching;
    this.studentStates = new Map();
    this.suggestionHistory = [];
  }

  /**
   * 更新学生状态
   */
  updateStudentState(studentId, state) {
    this.studentStates.set(studentId, {
      ...state,
      lastUpdate: Date.now(),
    });
  }

  /**
   * 分析全班状态并生成教学建议
   */
  generateClassSuggestions() {
    const stuckStudents = this._getStuckStudents();
    const suggestions = [];

    if (stuckStudents.length > 0) {
      const stuckRatio = stuckStudents.length / this.studentStates.size;
      
      if (stuckRatio > this.config.minStuckThreshold) {
        suggestions.push({
          type: 'alert',
          priority: 'high',
          text: `${stuckStudents.length} 名学生在同类问题上卡壳，建议切回大屏集中讲解。`,
          action: 'switch_to_lecture',
        });
      }

      const commonError = this._findCommonError(stuckStudents);
      if (commonError) {
        suggestions.push({
          type: 'info',
          priority: 'medium',
          text: `常见错误: ${commonError.type}，建议在实验手册中增加排查指南。`,
          action: 'update_manual',
        });
      }
    }

    return suggestions.slice(0, this.config.maxSuggestionsPerSession);
  }

  /**
   * 获取卡壳学生列表
   */
  _getStuckStudents() {
    const now = Date.now();
    const timeout = 5 * 60 * 1000; // 5分钟

    return Array.from(this.studentStates.entries())
      .filter(([_, state]) => {
        return state.isStuck && (now - state.lastUpdate) < timeout;
      })
      .map(([id, state]) => ({ id, ...state }));
  }

  /**
   * 查找共同错误
   */
  _findCommonError(stuckStudents) {
    const errorCounts = {};
    
    stuckStudents.forEach(s => {
      if (s.lastError) {
        errorCounts[s.lastError.type] = (errorCounts[s.lastError.type] || 0) + 1;
      }
    });

    const mostCommon = Object.entries(errorCounts)
      .sort((a, b) => b[1] - a[1])[0];

    if (mostCommon && mostCommon[1] > 1) {
      return { type: mostCommon[0], count: mostCommon[1] };
    }

    return null;
  }
}

/**
 * 代码注入模块
 */
class CodeInjectionModule {
  constructor() {
    this.injectQueue = [];
    this.isProcessing = false;
  }

  /**
   * 注入代码到终端
   */
  async injectCode(terminal, code) {
    if (!terminal) {
      throw new Error('Terminal instance is required');
    }

    this.injectQueue.push(code);
    
    if (!this.isProcessing) {
      await this._processQueue(terminal);
    }
  }

  /**
   * 处理注入队列
   */
  async _processQueue(terminal) {
    this.isProcessing = true;

    while (this.injectQueue.length > 0) {
      const code = this.injectQueue.shift();
      
      terminal.pasteText(code);
      terminal.pasteText('\n');
      
      await this._delay(100);
    }

    this.isProcessing = false;
  }

  _delay(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
  }
}

/**
 * AI 智能体核心控制器
 */
class AIAgentCore {
  constructor() {
    this.config = AI_CONFIG;
    this.rag = new RAGModule(this.config);
    this.diagnostic = new DiagnosticModule(this.config);
    this.teaching = new TeachingSuggestionModule(this.config);
    this.codeInjection = new CodeInjectionModule();
    
    this.isRunning = false;
    this.sessionId = null;
  }

  /**
   * 初始化 AI 系统
   */
  async initialize() {
    console.log('[AI Agent] Initializing core system...');
    
    await this.rag.initialize();
    
    this.isRunning = true;
    this.sessionId = `session_${Date.now()}`;
    
    console.log(`[AI Agent] System ready. Session: ${this.sessionId}`);
    return this.sessionId;
  }

  /**
   * 处理学生问题
   */
  async handleStudentQuery(studentId, query, context = {}) {
    const ragResults = await this.rag.retrieve(query);
    const answer = await this.rag.generateAnswer(query, ragResults);
    
    this.teaching.updateStudentState(studentId, {
      lastQuery: query,
      lastQueryTime: Date.now(),
      context,
    });

    return answer;
  }

  /**
   * 处理终端错误
   */
  async handleTerminalError(studentId, errorLog) {
    const diagnosis = await this.diagnostic.analyzeError(errorLog);
    
    this.teaching.updateStudentState(studentId, {
      lastError: diagnosis,
      isStuck: diagnosis.severity === 'critical' || diagnosis.severity === 'error',
    });

    return diagnosis;
  }

  /**
   * 获取教学建议
   */
  getTeachingSuggestions() {
    return this.teaching.generateClassSuggestions();
  }

  /**
   * 注入代码
   */
  async injectCode(terminal, code) {
    return this.codeInjection.injectCode(terminal, code);
  }

  /**
   * 获取系统状态
   */
  getStatus() {
    return {
      isRunning: this.isRunning,
      sessionId: this.sessionId,
      ragInitialized: this.rag.isInitialized,
      diagnosticListening: this.diagnostic.isListening,
      studentCount: this.teaching.studentStates.size,
      errorStats: this.diagnostic.getErrorStats(),
    };
  }

  /**
   * 关闭 AI 系统
   */
  async shutdown() {
    console.log('[AI Agent] Shutting down...');
    this.diagnostic.stopListening();
    this.isRunning = false;
    this.sessionId = null;
    console.log('[AI Agent] System shutdown complete');
  }
}

// ==================== 导出模块 ====================

export {
  AI_CONFIG,
  RAGModule,
  DiagnosticModule,
  TeachingSuggestionModule,
  CodeInjectionModule,
  AIAgentCore,
};

// 默认导出核心控制器
export default AIAgentCore;
