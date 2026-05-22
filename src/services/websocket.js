const WS_BASE = import.meta.env.VITE_WS_BASE || 'ws://localhost:8080/ws'

class WebSocketService {
  constructor() {
    this.ws = null
    this.reconnectTimer = null
    this.listeners = {}
    this.isConnecting = false
    this.shouldReconnect = true
    this.reconnectAttempts = 0
    this.maxReconnectAttempts = 5
    this.reconnectDelay = 3000
  }

  connect(token) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      return
    }

    if (this.isConnecting) {
      return
    }

    this.isConnecting = true
    const url = `${WS_BASE}?token=${token}`
    this.ws = new WebSocket(url)

    this.ws.onopen = () => {
      this.isConnecting = false
      this.reconnectAttempts = 0
      this.emit('connected')
    }

    this.ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        this.emit(data.type || 'message', data)
      } catch (err) {
        console.error('WebSocket message parse error:', err)
      }
    }

    this.ws.onerror = (error) => {
      this.isConnecting = false
      this.emit('error', error)
    }

    this.ws.onclose = () => {
      this.isConnecting = false
      this.emit('disconnected')

      if (this.shouldReconnect && this.reconnectAttempts < this.maxReconnectAttempts) {
        this.reconnectAttempts++
        this.reconnectTimer = setTimeout(() => {
          this.connect(token)
        }, this.reconnectDelay)
      }
    }
  }

  disconnect() {
    this.shouldReconnect = false
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
    this.listeners = {}
  }

  on(event, callback) {
    if (!this.listeners[event]) {
      this.listeners[event] = []
    }
    this.listeners[event].push(callback)
  }

  off(event, callback) {
    if (this.listeners[event]) {
      this.listeners[event] = this.listeners[event].filter(cb => cb !== callback)
    }
  }

  emit(event, data) {
    if (this.listeners[event]) {
      this.listeners[event].forEach(callback => {
        try {
          callback(data)
        } catch (err) {
          console.error(`WebSocket listener error for event "${event}":`, err)
        }
      })
    }
  }

  send(type, data) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({ type, ...data }))
    } else {
      console.warn('WebSocket is not connected')
    }
  }

  subscribeVMStatus(vmid) {
    this.send('subscribe', { type: 'vm_status', vmid })
  }

  subscribeNodeStatus(nodeId) {
    this.send('subscribe', { type: 'node_status', nodeId })
  }

  subscribeTaskStatus() {
    this.send('subscribe', { type: 'task_status' })
  }

  subscribeAll() {
    this.send('subscribe', { type: 'all' })
  }
}

export const wsService = new WebSocketService()
export default wsService
