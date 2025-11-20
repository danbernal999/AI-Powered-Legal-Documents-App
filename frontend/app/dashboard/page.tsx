'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import Sidebar from '@/components/Sidebar'
import { useAuthStore, useDocumentStore, useTemplateStore } from '@/lib/store'
import { documentAPI, templateAPI } from '@/lib/api'
import { useTranslation } from '@/hooks/useTranslation'
import { IconDocumentFile } from '@/components/Icons'

const getStatusColor = (status: string) => {
  switch (status?.toLowerCase()) {
    case 'draft':
      return { bg: 'bg-yellow-100', text: 'text-yellow-700', label: 'Draft' }
    case 'signed':
      return { bg: 'bg-green-100', text: 'text-green-700', label: 'Signed' }
    case 'in review':
      return { bg: 'bg-blue-100', text: 'text-blue-700', label: 'In Review' }
    case 'completed':
      return { bg: 'bg-purple-100', text: 'text-purple-700', label: 'Completed' }
    default:
      return { bg: 'bg-gray-100', text: 'text-gray-700', label: status || 'Draft' }
  }
}

export default function DashboardPage() {
  const router = useRouter()
  const { token, user, logout } = useAuthStore()
  const { documents, setDocuments } = useDocumentStore()
  const { templates, setTemplates } = useTemplateStore()
  const { t, interpolate } = useTranslation()
  const [loading, setLoading] = useState(true)
  const [recentDocs, setRecentDocs] = useState<any[]>([])
  const [searchQuery, setSearchQuery] = useState('')

  useEffect(() => {
    if (!token) {
      router.push('/login')
      return
    }

    const fetchData = async () => {
      try {
        const [docsResponse, templatesResponse] = await Promise.all([
          documentAPI.list(),
          templateAPI.getAll()
        ])
        setDocuments(docsResponse.data)
        setRecentDocs(docsResponse.data.slice(0, 4))
        setTemplates(templatesResponse.data)
      } catch (error) {
        console.error('Failed to fetch data:', error)
      } finally {
        setLoading(false)
      }
    }

    fetchData()
  }, [token, router, setDocuments, setTemplates])

  const handleLogout = () => {
    logout()
    router.push('/')
  }

  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchQuery(e.target.value)
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
        <div className="container mx-auto px-4 py-6 max-w-6xl">
          <div className="mb-8">
            <div className="flex flex-col md:flex-row items-start md:items-center justify-between gap-4 mb-6">
              <div>
                <h1 className="text-4xl font-bold text-gray-900 mb-2 mt-6">
                  {interpolate(t('dashboard.welcome'), { name: user?.name?.split(' ')[0] || '' })}
                </h1>
                <p className="text-gray-600">{t('dashboard.whatWouldYouLike')}</p>
              </div>
              <Link href="/documents/create" className="btn btn-primary gap-2 whitespace-nowrap mt-6 md:mt-0">
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                </svg>
                {t('dashboard.createNewDocument')}
              </Link>
            </div>

            <div className="relative max-w-md">
              <svg className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
              <input
                type="text"
                placeholder={t('dashboard.searchDocuments')}
                value={searchQuery}
                onChange={handleSearch}
                className="w-full pl-10 pr-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>

          <div className="mb-12">
            <h2 className="heading-2 mb-6">{t('dashboard.quickStart')}</h2>
            <div className="grid md:grid-cols-3 gap-6">
              <Link href="/documents/create" className="card hover:shadow-lg transition-shadow group cursor-pointer">
                <div className="w-14 h-14 bg-teal-100 rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
                  <svg className="w-8 h-8 text-teal-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M13 10V3L4 14h7v7l9-11h-7z" />
                  </svg>
                </div>
                <h3 className="heading-3 mb-2 group-hover:text-blue-600 transition-colors">{t('dashboard.generateContract')}</h3>
                <p className="text-gray-600 text-sm">{t('dashboard.generateContractDesc')}</p>
              </Link>

              <div className="card hover:shadow-lg transition-shadow cursor-pointer">
                <div className="w-14 h-14 bg-teal-100 rounded-lg flex items-center justify-center mb-4 hover:scale-110 transition-transform">
                  <svg className="w-8 h-8 text-teal-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
                <h3 className="heading-3 mb-2 hover:text-blue-600 transition-colors">{t('dashboard.analyzeDocument')}</h3>
                <p className="text-gray-600 text-sm">{t('dashboard.analyzeDocumentDesc')}</p>
              </div>

              <Link href="/templates" className="card hover:shadow-lg transition-shadow group cursor-pointer">
                <div className="w-14 h-14 bg-teal-100 rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
                  <svg className="w-8 h-8 text-teal-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4" />
                  </svg>
                </div>
                <h3 className="heading-3 mb-2 group-hover:text-blue-600 transition-colors">{t('dashboard.startTemplate')}</h3>
                <p className="text-gray-600 text-sm">{t('dashboard.startTemplateDesc')}</p>
              </Link>
            </div>
          </div>

          {recentDocs.length > 0 && (
            <div className="mb-12">
              <div className="flex items-center justify-between mb-6">
                <h2 className="heading-2">{t('dashboard.recentDocuments')}</h2>
                <Link href="/documents" className="text-blue-600 hover:text-blue-700 font-semibold text-sm">
                  {t('dashboard.viewAll')} →
                </Link>
              </div>
              
              <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
                {recentDocs.map((doc: any) => {
                  const status = getStatusColor(doc.status)
                  return (
                    <Link key={doc.id} href={`/documents/${doc.id}`}>
                      <div className="card cursor-pointer group overflow-hidden hover:shadow-lg transition-shadow h-full flex flex-col">
                        <div className="flex items-start justify-between mb-4">
                          <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center group-hover:scale-110 transition-transform">
                            <span className="text-2xl"><IconDocumentFile /></span>
                          </div>
                        </div>
                        <h3 className="heading-3 mb-2 group-hover:text-blue-600 transition-colors line-clamp-2">{doc.title}</h3>
                        <p className="text-xs text-gray-500 mb-3">{doc.type}</p>
                        <div className="flex items-center justify-between mt-auto pt-3 border-t border-gray-100">
                          <p className="text-xs text-gray-500">
                            {t('dashboard.modified')} {new Date(doc.updated_at || doc.created_at).toLocaleDateString('es-CO', {
                              year: 'numeric',
                              month: 'short',
                              day: 'numeric'
                            })}
                          </p>
                          <span className={`px-3 py-1 ${status.bg} ${status.text} rounded-full text-xs font-semibold`}>
                            {status.label}
                          </span>
                        </div>
                      </div>
                    </Link>
                  )
                })}
              </div>
            </div>
          )}

          {templates && templates.length > 0 && (
            <div>
              <div className="flex items-center justify-between mb-6">
                <h2 className="heading-2">{t('dashboard.popularTemplates')}</h2>
                <Link href="/templates" className="text-blue-600 hover:text-blue-700 font-semibold text-sm">
                  {t('dashboard.viewAll')} →
                </Link>
              </div>
              
              <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
                {templates.slice(0, 4).map((template: any) => (
                  <Link key={template.id} href={`/documents/create?templateId=${template.id}`}>
                    <div className="card cursor-pointer group overflow-hidden hover:shadow-lg transition-shadow h-full flex flex-col">
                      <div className="w-12 h-12 bg-slate-100 rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
                        <svg className="w-6 h-6 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                        </svg>
                      </div>
                      <h3 className="heading-3 mb-2 group-hover:text-blue-600 transition-colors">{template.name}</h3>
                      <p className="text-sm text-gray-600 line-clamp-2 flex-1 mb-4">{template.description}</p>
                      <div className="flex items-center text-blue-600 text-sm font-semibold group-hover:gap-2 transition-all gap-1">
                        {t('dashboard.useTemplate')} <span>→</span>
                      </div>
                    </div>
                  </Link>
                ))}
              </div>
            </div>
          )}
        </div>
      </main>
    </div>
  )
}
