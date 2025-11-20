import axios from 'axios'
import Cookies from 'js-cookie'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

apiClient.interceptors.request.use((config) => {
  const token = Cookies.get('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export const authAPI = {
  register: (email: string, password: string, name: string) =>
    apiClient.post('/auth/register', { email, password, name }),

  login: (email: string, password: string) =>
    apiClient.post('/auth/login', { email, password }),
}

export const templateAPI = {
  getAll: () => apiClient.get('/templates'),
  getById: (id: string) => apiClient.get(`/templates/${id}`),
  create: (data: any) => apiClient.post('/templates', data),
  update: (id: string, data: any) => apiClient.put(`/templates/${id}`, data),
  delete: (id: string) => apiClient.delete(`/templates/${id}`),
}

export const documentAPI = {
  list: () => apiClient.get('/documents'),
  create: (data: any) => apiClient.post('/documents', data),
  getById: (id: string) => apiClient.get(`/documents/${id}`),
  update: (id: string, data: any) => apiClient.put(`/documents/${id}`, data),
  delete: (id: string) => apiClient.delete(`/documents/${id}`),
}

export const generateAPI = {
  generateDocument: (templateId: string, variables: any, title: string) =>
    apiClient.post('/generate', {
      template_id: templateId,
      variables,
      title,
    }),
}

export default apiClient
