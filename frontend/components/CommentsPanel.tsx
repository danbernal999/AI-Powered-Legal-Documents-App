'use client'

import { useState, useEffect } from 'react'
import { collaborationAPI } from '@/lib/api'
import { useAuthStore } from '@/lib/store'

interface Comment {
  id: string
  document_id: string
  user_id: string
  user_name: string
  content: string
  created_at: string
  updated_at: string
}

interface CommentsPanelProps {
  documentId: string
  onCommentAdded?: () => void
}

export default function CommentsPanel({
  documentId,
  onCommentAdded,
}: CommentsPanelProps) {
  const { user } = useAuthStore()
  const [comments, setComments] = useState<Comment[]>([])
  const [newComment, setNewComment] = useState('')
  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    fetchComments()
  }, [documentId])

  const fetchComments = async () => {
    setLoading(true)
    try {
      const response = await collaborationAPI.getDocumentComments(documentId)
      setComments(response.data || [])
    } catch (error) {
      console.error('Error cargando comentarios:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleAddComment = async () => {
    if (!newComment.trim()) {
      return
    }

    setSubmitting(true)
    try {
      await collaborationAPI.createComment(documentId, newComment)
      setNewComment('')
      await fetchComments()
      onCommentAdded?.()
    } catch (error) {
      console.error('Error creando comentario:', error)
      alert('No se pudo crear el comentario')
    } finally {
      setSubmitting(false)
    }
  }

  const handleDeleteComment = async (commentId: string) => {
    if (!confirm('¿Estás seguro de que deseas eliminar este comentario?')) return

    try {
      await collaborationAPI.deleteComment(commentId)
      await fetchComments()
    } catch (error) {
      console.error('Error eliminando comentario:', error)
      alert('No se pudo eliminar el comentario')
    }
  }

  if (loading) {
    return (
      <div className="p-4 text-center text-gray-500">
        Cargando comentarios...
      </div>
    )
  }

  return (
    <div className="space-y-4">
      <div className="space-y-3">
        {comments.length === 0 ? (
          <p className="text-sm text-gray-500 text-center py-4">
            No hay comentarios aún
          </p>
        ) : (
          comments.map((comment) => (
            <div
              key={comment.id}
              className="p-4 bg-gray-50 rounded-lg border border-gray-200"
            >
              <div className="flex items-start justify-between mb-2">
                <div>
                  <p className="text-sm font-semibold text-gray-900">
                    {comment.user_name || 'Usuario'}
                  </p>
                  <p className="text-xs text-gray-500">
                    {new Date(comment.created_at).toLocaleString()}
                  </p>
                </div>
                {user?.id === comment.user_id && (
                  <button
                    onClick={() => handleDeleteComment(comment.id)}
                    className="p-1.5 text-gray-500 hover:text-red-600 hover:bg-red-50 rounded transition"
                    title="Eliminar comentario"
                  >
                    <i className="ri-trash-2-line text-sm" />
                  </button>
                )}
              </div>
              <p className="text-sm text-gray-700 leading-relaxed">
                {comment.content}
              </p>
            </div>
          ))
        )}
      </div>

      <div className="border-t border-gray-200 pt-4">
        <label className="block text-sm font-semibold text-gray-900 mb-2">
          Agregar comentario
        </label>
        <textarea
          value={newComment}
          onChange={(e) => setNewComment(e.target.value)}
          placeholder="Escribe tu comentario..."
          rows={3}
          className="w-full px-4 py-2 border-2 border-gray-200 rounded-lg focus:border-indigo-600 focus:outline-none resize-none"
          disabled={submitting}
        />
        <button
          onClick={handleAddComment}
          disabled={!newComment.trim() || submitting}
          className="mt-2 px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition disabled:opacity-50 font-semibold text-sm"
        >
          {submitting ? 'Enviando...' : 'Enviar comentario'}
        </button>
      </div>
    </div>
  )
}
