'use client'

import { useState, useEffect } from 'react'
import { collaborationAPI } from '@/lib/api'

interface Collaborator {
  id: string
  document_id: string
  shared_by_user_id: string
  shared_with_user_id: string
  permission: string
  created_at: string
}

interface CollaboratorsPanelProps {
  documentId: string
  isLoading?: boolean
}

export default function CollaboratorsPanel({
  documentId,
  isLoading = false,
}: CollaboratorsPanelProps) {
  const [collaborators, setCollaborators] = useState<Collaborator[]>([])
  const [loading, setLoading] = useState(isLoading)

  useEffect(() => {
    fetchCollaborators()
  }, [documentId])

  const fetchCollaborators = async () => {
    setLoading(true)
    try {
      const response = await collaborationAPI.getDocumentShares(documentId)
      setCollaborators(response.data || [])
    } catch (error) {
      console.error('Error cargando colaboradores:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleRemoveCollaborator = async (shareId: string) => {
    if (!confirm('¿Estás seguro de que deseas revocar el acceso?')) return

    try {
      await collaborationAPI.deleteShare(shareId)
      await fetchCollaborators()
    } catch (error) {
      console.error('Error removiendo colaborador:', error)
    }
  }

  const getPermissionLabel = (permission: string) => {
    const labels: Record<string, string> = {
      read: 'Solo lectura',
      edit: 'Puede editar',
      admin: 'Administrador',
    }
    return labels[permission] || permission
  }

  const getPermissionColor = (permission: string) => {
    const colors: Record<string, string> = {
      read: 'bg-blue-100 text-blue-800',
      edit: 'bg-green-100 text-green-800',
      admin: 'bg-purple-100 text-purple-800',
    }
    return colors[permission] || 'bg-gray-100 text-gray-800'
  }

  if (loading) {
    return (
      <div className="p-4 text-center text-gray-500">
        Cargando colaboradores...
      </div>
    )
  }

  if (collaborators.length === 0) {
    return (
      <div className="p-4 text-center text-gray-500">
        Este documento aún no ha sido compartido con nadie
      </div>
    )
  }

  return (
    <div className="space-y-3">
      {collaborators.map((collab) => (
        <div
          key={collab.id}
          className="flex items-center justify-between p-4 bg-gray-50 rounded-lg border border-gray-200 hover:border-gray-300 transition"
        >
          <div className="flex-1">
            <p className="text-sm font-semibold text-gray-900">
              {collab.shared_with_user_id}
            </p>
            <p className="text-xs text-gray-500 mt-1">
              Compartido el {new Date(collab.created_at).toLocaleDateString()}
            </p>
          </div>

          <div className="flex items-center gap-3">
            <span
              className={`inline-block px-3 py-1 rounded-full text-xs font-semibold ${getPermissionColor(
                collab.permission
              )}`}
            >
              {getPermissionLabel(collab.permission)}
            </span>

            <button
              onClick={() => handleRemoveCollaborator(collab.id)}
              className="p-2 text-gray-500 hover:text-red-600 hover:bg-red-50 rounded-lg transition"
              title="Revocar acceso"
            >
              <i className="ri-close-line" />
            </button>
          </div>
        </div>
      ))}
    </div>
  )
}
