import axios, { AxiosError, AxiosResponse, InternalAxiosRequestConfig } from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_URL || '/api/v1'

// Create axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor to add auth token
api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('authToken')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error: AxiosError) => {
    return Promise.reject(error)
  }
)

// Response interceptor for error handling and retry logic
api.interceptors.response.use(
  (response: AxiosResponse) => {
    return response
  },
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & {
      headers: Record<string, string>
    }

    // Handle 401 errors (unauthorized)
    if (error.response?.status === 401) {
      localStorage.removeItem('authToken')
      localStorage.removeItem('user')
      window.location.href = '/login'
      return Promise.reject(error)
    }

    // Retry logic for network errors
    if (
      !error.response &&
      originalRequest &&
      !originalRequest.headers['X-Retry-Count']
    ) {
      originalRequest.headers['X-Retry-Count'] = '1'
      return api(originalRequest)
    }

    return Promise.reject(error)
  }
)

export default api

// API Response types
export interface ApiResponse<T = any> {
  success: boolean
  data?: T
  error?: {
    code: string
    message: string
  }
}

export interface PaginatedResponse<T = any> {
  data: T[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}
