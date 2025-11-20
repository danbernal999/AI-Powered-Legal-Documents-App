import { create } from 'zustand'
import Cookies from 'js-cookie'

interface User {
  id: string
  email: string
  name: string
  created_at?: string
}

interface AuthStore {
  token: string | null
  user: User | null
  isLoading: boolean
  error: string | null
  login: (token: string, user: User) => void
  logout: () => void
  setError: (error: string | null) => void
  setLoading: (loading: boolean) => void
  initializeAuth: () => void
}

export const useAuthStore = create<AuthStore>((set) => ({
  token: null,
  user: null,
  isLoading: false,
  error: null,

  login: (token: string, user: User) => {
    Cookies.set('token', token, { expires: 7 })
    set({ token, user, error: null })
  },

  logout: () => {
    Cookies.remove('token')
    set({ token: null, user: null })
  },

  setError: (error: string | null) => {
    set({ error })
  },

  setLoading: (loading: boolean) => {
    set({ isLoading: loading })
  },

  initializeAuth: () => {
    const token = Cookies.get('token')
    if (token) {
      set({ token })
    }
  },
}))

interface DocumentStore {
  documents: any[]
  currentDocument: any | null
  isLoading: boolean
  error: string | null
  setDocuments: (documents: any[]) => void
  setCurrentDocument: (document: any) => void
  setLoading: (loading: boolean) => void
  setError: (error: string | null) => void
}

export const useDocumentStore = create<DocumentStore>((set) => ({
  documents: [],
  currentDocument: null,
  isLoading: false,
  error: null,

  setDocuments: (documents: any[]) => {
    set({ documents, error: null })
  },

  setCurrentDocument: (document: any) => {
    set({ currentDocument: document })
  },

  setLoading: (loading: boolean) => {
    set({ isLoading: loading })
  },

  setError: (error: string | null) => {
    set({ error })
  },
}))

interface TemplateStore {
  templates: any[]
  isLoading: boolean
  error: string | null
  setTemplates: (templates: any[]) => void
  setLoading: (loading: boolean) => void
  setError: (error: string | null) => void
}

export const useTemplateStore = create<TemplateStore>((set) => ({
  templates: [],
  isLoading: false,
  error: null,

  setTemplates: (templates: any[]) => {
    set({ templates, error: null })
  },

  setLoading: (loading: boolean) => {
    set({ isLoading: loading })
  },

  setError: (error: string | null) => {
    set({ error })
  },
}))
