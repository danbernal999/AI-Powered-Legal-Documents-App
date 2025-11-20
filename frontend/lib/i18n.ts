import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export type Language = 'en' | 'es' | 'pt'

interface I18nStore {
  language: Language
  setLanguage: (language: Language) => void
}

export const useLanguage = create<I18nStore>()(
  persist(
    (set) => ({
      language: 'en',
      setLanguage: (language: Language) => set({ language }),
    }),
    {
      name: 'language-storage',
    }
  )
)
