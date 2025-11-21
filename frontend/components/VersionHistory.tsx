'use client'

import { useState, useEffect } from 'react'
import { collaborationAPI } from '@/lib/api'

interface Version {
  id: string
  document_id: string
  user_id: string
  user_name: string
  previous_content: string
  current_content: string
  change_summary: string
  created_at: string
}

interface VersionHistoryProps {
  documentId: string
}

export default function VersionHistory({ documentId }: VersionHistoryProps) {
  const [versions, setVersions] = useState<Version[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedVersion, setSelectedVersion] = useState<string | null>(null)

  useEffect(() => {
    fetchVersions()
  }, [documentId])

  const fetchVersions = async () => {
    setLoading(true)
    try {
      const response = await collaborationAPI.getDocumentVersions(documentId)
      setVersions(response.data || [])
    } catch (error) {
      console.error('Error cargando historial de versiones:', error)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="p-4 text-center text-gray-500">
        Cargando historial...
      </div>
    )
  }

  if (versions.length === 0) {
    return (
      <div className="p-4 text-center text-gray-500">
        No hay historial de cambios aún
      </div>
    )
  }

  return (
    <div className="space-y-4">
      <div className="space-y-2">
        {versions.map((version, index) => (
          <button
            key={version.id}
            onClick={() =>
              setSelectedVersion(
                selectedVersion === version.id ? null : version.id
              )
            }
            className="w-full text-left p-4 bg-gray-50 rounded-lg border border-gray-200 hover:border-indigo-300 hover:bg-indigo-50 transition"
          >
            <div className="flex items-start justify-between">
              <div className="flex-1">
                <div className="flex items-center gap-2 mb-1">
                  <span className="inline-block px-2 py-1 bg-indigo-100 text-indigo-700 rounded text-xs font-semibold">
                    v{versions.length - index}
                  </span>
                  <p className="text-sm font-semibold text-gray-900">
                    {version.user_name || 'Usuario'} • 
                    {new Date(version.created_at).toLocaleString()}
                  </p>
                </div>
                {version.change_summary && (
                  <p className="text-sm text-gray-600">
                    {version.change_summary}
                  </p>
                )}
              </div>
              <i
                className={`ri-chevron-down-line text-gray-400 transition ${
                  selectedVersion === version.id ? 'rotate-180' : ''
                }`}
              />
            </div>

            {selectedVersion === version.id && (
              <div className="mt-4 pt-4 border-t border-gray-200 space-y-3">
                <div>
                  <h4 className="text-xs font-semibold text-gray-900 mb-2">
                    Contenido anterior (primeros 300 caracteres):
                  </h4>
                  <div className="p-3 bg-red-50 rounded border border-red-200 text-xs text-gray-700 overflow-auto max-h-32 whitespace-pre-wrap">
                    {version.previous_content
                      ? version.previous_content.substring(0, 300) + '...'
                      : 'N/A'}
                  </div>
                </div>

                <div>
                  <h4 className="text-xs font-semibold text-gray-900 mb-2">
                    Contenido nuevo (primeros 300 caracteres):
                  </h4>
                  <div className="p-3 bg-green-50 rounded border border-green-200 text-xs text-gray-700 overflow-auto max-h-32 whitespace-pre-wrap">
                    {version.current_content.substring(0, 300) + '...'}
                  </div>
                </div>
              </div>
            )}
          </button>
        ))}
      </div>
    </div>
  )
}
