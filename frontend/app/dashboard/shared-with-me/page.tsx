'use client'

import { useEffect, useState } from 'react'
import { documentAPI } from '@/lib/api'
import Link from 'next/link'

interface SharedDocument {
  share: {
    id: string
    document_id: string
    permission: string
    created_at: string
  }
  document: {
    id: string
    title: string
    content: string
    created_at: string
    updated_at: string
  }
  shared_by: {
    id: string
    name: string
    email: string
  }
}

export default function SharedWithMePage() {
  const [sharedDocs, setSharedDocs] = useState<SharedDocument[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    loadSharedDocuments()
  }, [])

  const loadSharedDocuments = async () => {
    try {
      setLoading(true)
      const response = await documentAPI.listSharedWithMe()
      setSharedDocs(response.data || [])
      setError('')
    } catch (err) {
      console.error('Error loading shared documents:', err)
      setError('No se pudieron cargar los documentos compartidos')
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Cargando documentos compartidos...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">Documentos Compartidos Conmigo</h1>
          <p className="text-gray-600">Documentes que otros usuarios han compartido contigo</p>
        </div>

        {error && (
          <div className="mb-6 p-4 bg-red-50 border-l-4 border-red-500 rounded">
            <p className="text-red-700">{error}</p>
          </div>
        )}

        {sharedDocs.length === 0 ? (
          <div className="bg-white rounded-lg shadow p-8 text-center">
            <div className="text-gray-400 mb-4">
              <svg
                className="mx-auto h-12 w-12"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                />
              </svg>
            </div>
            <h3 className="text-lg font-semibold text-gray-900 mb-2">No hay documentos compartidos</h3>
            <p className="text-gray-600">Cuando alguien comparta un documento contigo, aparecerá aquí</p>
          </div>
        ) : (
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {sharedDocs.map((item) => (
              <Link key={item.share.id} href={`/documents/${item.document.id}`}>
                <div className="bg-white rounded-lg shadow hover:shadow-lg transition p-6 cursor-pointer">
                  <div className="mb-4">
                    <h3 className="text-lg font-semibold text-gray-900 truncate mb-2">
                      {item.document.title}
                    </h3>
                    <p className="text-sm text-gray-600">
                      Compartido por <span className="font-semibold">{item.shared_by.name}</span>
                    </p>
                    <p className="text-xs text-gray-500 mt-1">{item.shared_by.email}</p>
                  </div>

                  <div className="flex items-center justify-between pt-4 border-t border-gray-200">
                    <div>
                      <span className={`inline-block px-3 py-1 rounded-full text-xs font-semibold ${
                        item.share.permission === 'edit'
                          ? 'bg-blue-100 text-blue-800'
                          : item.share.permission === 'admin'
                          ? 'bg-purple-100 text-purple-800'
                          : 'bg-gray-100 text-gray-800'
                      }`}>
                        {item.share.permission === 'edit'
                          ? 'Puedo editar'
                          : item.share.permission === 'admin'
                          ? 'Administrador'
                          : 'Solo lectura'}
                      </span>
                    </div>
                    <span className="text-xs text-gray-500">
                      {new Date(item.share.created_at).toLocaleDateString()}
                    </span>
                  </div>
                </div>
              </Link>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
