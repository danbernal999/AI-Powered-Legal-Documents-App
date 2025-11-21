'use client'

import { useState } from 'react'
import { collaborationAPI } from '@/lib/api'

interface ShareDocumentModalProps {
  documentId: string
  isOpen: boolean
  onClose: () => void
  onShareSuccess: () => void
}

export default function ShareDocumentModal({
  documentId,
  isOpen,
  onClose,
  onShareSuccess,
}: ShareDocumentModalProps) {
  const [email, setEmail] = useState('')
  const [permission, setPermission] = useState('read')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleShare = async () => {
    if (!email.trim()) {
      setError('Por favor ingresa un email')
      return
    }

    setLoading(true)
    setError('')
    try {
      await collaborationAPI.shareDocument(documentId, email, permission)
      setEmail('')
      setPermission('read')
      onShareSuccess()
      onClose()
    } catch (err) {
      setError('No se pudo compartir el documento')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  if (!isOpen) return null

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-white rounded-xl shadow-2xl max-w-md w-full mx-4">
        <div className="p-6 border-b border-gray-200">
          <h2 className="heading-2 mb-2">Compartir Documento</h2>
          <p className="text-sm text-gray-600">
            Comparte este documento con otros usuarios
          </p>
        </div>

        <div className="p-6 space-y-4">
          {error && (
            <div className="p-3 bg-red-50 border-l-4 border-red-500 rounded">
              <p className="text-sm text-red-700">{error}</p>
            </div>
          )}

          <div>
            <label className="block text-sm font-semibold text-gray-900 mb-2">
              Email del usuario *
            </label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="usuario@ejemplo.com"
              className="w-full px-4 py-2 border-2 border-gray-200 rounded-lg focus:border-indigo-600 focus:outline-none"
              disabled={loading}
            />
          </div>

          <div>
            <label className="block text-sm font-semibold text-gray-900 mb-2">
              Permiso *
            </label>
            <select
              value={permission}
              onChange={(e) => setPermission(e.target.value)}
              className="w-full px-4 py-2 border-2 border-gray-200 rounded-lg focus:border-indigo-600 focus:outline-none"
              disabled={loading}
            >
              <option value="read">Solo lectura</option>
              <option value="edit">Puede editar</option>
              <option value="admin">Administrador</option>
            </select>
          </div>

          <div className="flex gap-2 pt-4">
            <button
              onClick={onClose}
              disabled={loading}
              className="flex-1 px-4 py-2 text-sm font-semibold text-gray-700 border border-gray-300 rounded-lg hover:bg-gray-50 transition disabled:opacity-50"
            >
              Cancelar
            </button>
            <button
              onClick={handleShare}
              disabled={loading}
              className="flex-1 px-4 py-2 text-sm font-semibold text-white bg-indigo-600 rounded-lg hover:bg-indigo-700 transition disabled:opacity-50"
            >
              {loading ? 'Compartiendo...' : 'Compartir'}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
