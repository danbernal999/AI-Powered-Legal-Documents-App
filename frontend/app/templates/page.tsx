'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import Sidebar from '@/components/Sidebar'
import { useAuthStore, useTemplateStore } from '@/lib/store'
import { templateAPI } from '@/lib/api'
import { useTranslation } from '@/hooks/useTranslation'
import { IconTemplates, IconDocumentFile, IconAdd } from '@/components/Icons'

export default function TemplatesPage() {
  const router = useRouter()
  const { token, user, logout } = useAuthStore()
  const { templates, setTemplates } = useTemplateStore()
  const { t } = useTranslation()
  const [loading, setLoading] = useState(true)

  const getTemplateTypeLabel = (type: string) => {
    switch(type?.toLowerCase()) {
      case 'custom': return t('templates.typeCustom')
      case 'nda': return t('templates.typeNDA')
      case 'employment': return t('templates.typeEmployment')
      case 'rental': return t('templates.typeRental')
      case 'freelance': return t('templates.typeFreelance')
      default: return type
    }
  }

  useEffect(() => {
    if (!token) {
      router.push('/login')
      return
    }

    const fetchTemplates = async () => {
      try {
        const response = await templateAPI.getAll()
        setTemplates(response.data)
      } catch (error) {
        console.error('Failed to fetch templates:', error)
      } finally {
        setLoading(false)
      }
    }

    fetchTemplates()
  }, [token, router, setTemplates])

  const handleLogout = () => {
    logout()
    router.push('/')
  }

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p className="text-gray-600">{t('common.loading')}</p>
      </div>
    )
  }

  return (
    <div className="flex min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
      <Sidebar user={user} onLogout={handleLogout} />

      <main className="flex-1 lg:ml-0">
        <div className="container mx-auto px-4 py-6 max-w-5xl">
          <div className="flex justify-between items-start mb-12">
            <div>
              <h2 className="heading-2 mb-2 mt-6">{t('templates.title')}</h2>
              <p className="text-gray-600">{t('templates.description')}</p>
            </div>
            <Link href="/templates/create" className="btn btn-primary gap-2 whitespace-nowrap mt-6">
              <IconAdd />
              {t('templates.createTemplate')}
            </Link>
          </div>

          {!templates || templates.length === 0 ? (
            <div className="card text-center py-16">
              <div className="text-6xl mb-4"><IconTemplates /></div>
              <p className="text-gray-600 mb-6 text-lg">{t('templates.noTemplates')}</p>
              <p className="text-gray-500">{t('templates.checkBackSoon')}</p>
            </div>
          ) : (
            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
              {templates.map((template: any) => (
                <Link key={template.id} href={`/documents/create?templateId=${template.id}`}>
                  <div className="card cursor-pointer group overflow-hidden hover:shadow-lg transition-shadow">
                    <div className="flex items-start justify-between mb-4">
                      <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center group-hover:scale-110 transition-transform">
                        <span className="text-2xl"><IconDocumentFile /></span>
                      </div>
                      <span className="px-3 py-1 bg-blue-100 text-blue-700 rounded-full text-xs font-semibold">
                        {getTemplateTypeLabel(template.type)}
                      </span>
                    </div>
                    <h3 className="heading-3 mb-2 group-hover:text-blue-600 transition-colors">
                      {template.name}
                    </h3>
                    <p className="text-sm text-gray-500 mb-4 line-clamp-2">
                      {template.description}
                    </p>
                    <div className="flex items-center text-blue-600 text-sm font-semibold group-hover:gap-2 transition-all gap-1">
                      {t('templates.useTemplate')} <span>→</span>
                    </div>
                  </div>
                </Link>
              ))}
            </div>
          )}
        </div>
      </main>
    </div>
  )
}
