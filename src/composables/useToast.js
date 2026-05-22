import { ref } from 'vue'

const toasts = ref([])

let toastId = 0

export function useToast() {
  const success = (message, duration = 3000) => {
    addToast(message, 'success', duration)
  }

  const error = (message, duration = 5000) => {
    addToast(message, 'error', duration)
  }

  const warning = (message, duration = 4000) => {
    addToast(message, 'warning', duration)
  }

  const info = (message, duration = 3000) => {
    addToast(message, 'info', duration)
  }

  const remove = (id) => {
    toasts.value = toasts.value.filter(t => t.id !== id)
  }

  return { toasts, success, error, warning, info, remove }
}

function addToast(message, type = 'info', duration = 3000) {
  const id = ++toastId
  const toast = { id, message, type }
  toasts.value.push(toast)

  if (duration > 0) {
    setTimeout(() => {
      remove(id)
    }, duration)
  }
}

function remove(id) {
  toasts.value = toasts.value.filter(t => t.id !== id)
}
