'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import Sidebar from '@/components/Sidebar'
import { useAuthStore, useDocumentStore } from '@/lib/store'
import { documentAPI } from '@/lib/api'
import { useTranslation } from '@/hooks/useTranslation'
import { IconAdd, IconDocumentFile } from '@/components/Icons'

export default function DocumentsPage() {
  const router = useRouter()
  const { token, user, logout } = useAuthStore()
  const { documents, setDocuments } = useDocumentStore()
  const { t } = useTranslation()
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!token) {
      router.push('/login')
      return
    }

    const fetchDocuments = async () => {
      try {
        const response = await documentAPI.list()
        setDocuments(response.data)
      } catch (error) {
        console.error('Failed to fetch documents:', error)
      } finally {
        setLoading(false)
      }
    }

    fetchDocuments()
  }, [token, router, setDocuments])

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
        <div className="container mx-auto px-4 py-6 max-w-6xl">
          <div className="flex justify-between items-center mb-8">
            <div>
              <h2 className="heading-2 mb-2 mt-6">{t('documents.title')}</h2>
              <p className="text-gray-600">{t('documents.subtitle')}</p>
            </div>
            <Link href="/documents/create" className="btn btn-primary gap-2">
              <IconAdd />
              {t('documents.create')}
            </Link>
          </div>

          {!documents || documents.length === 0 ? (
            <div className="card text-center py-16">
              <div className="text-6xl mb-4"><IconDocumentFile /></div>
              <p className="text-gray-600 mb-6 text-lg">{t('documents.noDocuments')}</p>
              <p className="text-gray-500 mb-8">{t('documents.noDocumentsDesc')}</p>
              <Link href="/documents/create" className="btn btn-primary gap-2">
                <IconAdd />
                {t('documents.create')}
              </Link>
            </div>
          ) : (
            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
              {documents.map((doc: any) => (
                <Link key={doc.id} href={`/documents/${doc.id}`}>
                  <div className="card cursor-pointer group overflow-hidden hover:shadow-lg transition-shadow">
                    <div className="flex items-start justify-between mb-4">
                      <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center group-hover:scale-110 transition-transform">
                        <span className="text-2xl"><IconDocumentFile /></span>
                      </div>
                      <span className="px-3 py-1 bg-blue-100 text-blue-700 rounded-full text-xs font-semibold">{doc.type}</span>
                    </div>
                    <h3 className="heading-3 mb-2 group-hover:text-blue-600 transition-colors">{doc.title}</h3>
                    <p className="text-sm text-gray-500 mb-4">
                      {t('documents.created')}: {new Date(doc.created_at).toLocaleDateString('es-CO', {
                        year: 'numeric',
                        month: 'short',
                        day: 'numeric'
                      })}
                    </p>
                    <div className="flex items-center text-blue-600 text-sm font-semibold group-hover:gap-2 transition-all gap-1">
                      {t('documents.view')} <span>→</span>
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
