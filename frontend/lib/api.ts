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

  googleLogin: (idToken: string) =>
    apiClient.post('/auth/google', { id_token: idToken }),
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
  listSharedWithMe: () => apiClient.get('/documents/shared-with-me'),
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

export const signatureAPI = {
  sign: (documentId: string, signerName: string, signerEmail: string, signatureData: string) =>
    apiClient.post(`/documents/${documentId}/sign`, {
      signer_name: signerName,
      signer_email: signerEmail,
      signature_data: signatureData,
    }),

  getSignatures: (documentId: string) =>
    apiClient.get(`/documents/${documentId}/signatures`),

  deleteSignature: (signatureId: string) =>
    apiClient.delete(`/signatures/${signatureId}`),
}

export const collaborationAPI = {
  shareDocument: (documentId: string, sharedWithEmail: string, permission: string) =>
    apiClient.post(`/documents/${documentId}/share`, {
      shared_with_email: sharedWithEmail,
      permission,
    }),

  getDocumentShares: (documentId: string) =>
    apiClient.get(`/documents/${documentId}/shares`),

  deleteShare: (shareId: string) =>
    apiClient.delete(`/shares/${shareId}`),

  getDocumentComments: (documentId: string) =>
    apiClient.get(`/documents/${documentId}/comments`),

  createComment: (documentId: string, content: string) =>
    apiClient.post(`/documents/${documentId}/comments`, { content }),

  deleteComment: (commentId: string) =>
    apiClient.delete(`/comments/${commentId}`),

  getDocumentVersions: (documentId: string) =>
    apiClient.get(`/documents/${documentId}/versions`),
}

export default apiClient
