import { useLanguage, Language } from '@/lib/i18n'
import en from '@/locales/en.json'
import es from '@/locales/es.json'
import pt from '@/locales/pt.json'

type Translations = {
  [key: string]: any
}

const translations: { [key in Language]: Translations } = {
  en,
  es,
  pt,
}

const loadTranslations = (language: Language): Translations => {
  return translations[language]
}

export const useTranslation = () => {
  const language = useLanguage((state) => state.language)
  const translations = loadTranslations(language)

  const t = (key: string, defaultValue?: string): string => {
    const keys = key.split('.')
    let value: any = translations

    for (const k of keys) {
      value = value?.[k]
    }

    return typeof value === 'string' ? value : defaultValue || key
  }

  const interpolate = (text: string, variables: { [key: string]: string }): string => {
    return text.replace(/\{\{(\w+)\}\}/g, (_, key) => variables[key] || `{{${key}}}`)
  }

  return { t, interpolate, language }
}
