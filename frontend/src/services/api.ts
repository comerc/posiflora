const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

export interface TelegramIntegration {
  id: number
  shop_id: number
  bot_token: string
  chat_id: string
  enabled: boolean
  created_at: string
  updated_at: string
}

export interface TelegramStatus {
  enabled: boolean
  chat_id: string
  last_sent_at?: string
  sent_count: number
  failed_count: number
}

export interface ConnectTelegramRequest {
  botToken: string
  chatId: string
  enabled: boolean
}

async function request<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`

  try {
    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    })

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: response.statusText }))
      throw new Error(error.error || `HTTP error! status: ${response.status}`)
    }

    return response.json()
  } catch (error) {
    // Обработка сетевых ошибок (CORS, connection refused и т.д.)
    if (error instanceof TypeError && error.message.includes('fetch')) {
      throw new Error(
        `Не удалось подключиться к серверу. Убедитесь, что backend запущен на ${API_BASE_URL}`,
      )
    }
    throw error
  }
}

export const api = {
  connectTelegram: (shopId: number, data: ConnectTelegramRequest): Promise<TelegramIntegration> => {
    return request<TelegramIntegration>(`/shops/${shopId}/telegram/connect`, {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },

  getTelegramStatus: (shopId: number, showFullChatId): Promise<TelegramStatus> => {
    const url = showFullChatId
      ? `/shops/${shopId}/telegram/status?mask=disabled`
      : `/shops/${shopId}/telegram/status`
    return request<TelegramStatus>(url)
  },
}
