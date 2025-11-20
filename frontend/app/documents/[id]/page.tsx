'use client'

import { useState, useEffect, useCallback } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import Sidebar from '@/components/Sidebar'
import { useAuthStore, useDocumentStore } from '@/lib/store'
import { documentAPI } from '@/lib/api'
import { IconEdit, IconSave, IconDownload, IconLoading, IconDocument, IconCheck, IconVersion, IconDate, IconDelete } from '@/components/Icons'

export default function DocumentPage({ params }: { params: { id: string } }) {
  const router = useRouter()
  const { token, user, logout } = useAuthStore()
  const { currentDocument, setCurrentDocument } = useDocumentStore()
  const [loading, setLoading] = useState(true)
  const [editing, setEditing] = useState(false)
  const [content, setContent] = useState('')
  const [saving, setSaving] = useState(false)
  const [saveSuccess, setSaveSuccess] = useState(false)
  const [exporting, setExporting] = useState(false)

  useEffect(() => {
    if (!token) {
      router.push('/login')
      return
    }

    if (!params?.id) {
      setLoading(false)
      return
    }

    const fetchDocument = async () => {
      setLoading(true)
      try {
        const response = await documentAPI.getById(params.id)
        setCurrentDocument(response.data)
        setContent(response.data.content)
      } catch (error) {
        console.error('Error cargando documento:', error)
        alert('No se pudo cargar el documento')
      } finally {
        setLoading(false)
      }
    }

    fetchDocument()
  }, [params.id, token, router, setCurrentDocument])

  const handleSave = useCallback(async () => {
    if (!currentDocument) return

    setSaving(true)
    try {
      const response = await documentAPI.update(currentDocument.id, {
        id: currentDocument.id,
        user_id: currentDocument.user_id,
        template_id: currentDocument.template_id,
        title: currentDocument.title,
        type: currentDocument.type,
        content,
        variables: typeof currentDocument.variables === 'string' ? JSON.parse(currentDocument.variables) : currentDocument.variables,
        status: currentDocument.status,
        version: currentDocument.version,
      })
      setCurrentDocument(response.data)
      setContent(response.data.content)
      setSaveSuccess(true)
      setEditing(false)
      setTimeout(() => setSaveSuccess(false), 3000)
    } catch (error) {
      console.error('Error guardando documento:', error)
      alert('No se pudo guardar el documento')
    } finally {
      setSaving(false)
    }
  }, [currentDocument, content, setCurrentDocument])

  const handleDelete = useCallback(async () => {
    if (!currentDocument || !confirm('¿Estás seguro?')) return

    try {
      await documentAPI.delete(currentDocument.id)
      router.push('/dashboard')
    } catch (error) {
      console.error('Error eliminando documento:', error)
      alert('No se pudo eliminar el documento')
    }
  }, [currentDocument, router])

  const handleExportPDF = useCallback(async () => {
    if (!currentDocument || !content) {
      alert('El documento no está completamente cargado')
      return
    }

    setExporting(true)
    try {
      const printWindow = window.open('', '_blank')
      if (!printWindow) {
        alert('Por favor permite las ventanas emergentes para exportar PDF')
        setExporting(false)
        return
      }

      const docTitle = currentDocument.title || 'Documento'
      const docType = currentDocument.type || 'Sin tipo'
      const docStatus = currentDocument.status || 'Sin estado'
      const docVersion = currentDocument.version || '1'
      const docContent = content || 'Sin contenido'

      const htmlContent = `
        <!DOCTYPE html>
        <html>
          <head>
            <meta charset="UTF-8">
            <title>${docTitle}</title>
            <style>
              body { 
                font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
                line-height: 1.6; 
                padding: 40px; 
                color: #333;
                background-color: white;
              }
              h1 { 
                font-size: 28px; 
                margin-bottom: 10px; 
                color: #1a1a1a;
                border-bottom: 2px solid #007bff;
                padding-bottom: 10px;
              }
              .metadata { 
                font-size: 13px; 
                color: #666; 
                margin-bottom: 30px;
                padding: 10px;
                background-color: #f5f5f5;
                border-radius: 4px;
              }
              .content { 
                white-space: pre-wrap; 
                word-wrap: break-word;
                line-height: 1.8;
                font-size: 14px;
              }
            </style>
          </head>
          <body>
            <h1>${docTitle}</h1>
            <div class="metadata">
              <strong>Tipo:</strong> ${docType} | 
              <strong>Estado:</strong> ${docStatus} | 
              <strong>Versión:</strong> v${docVersion}
            </div>
            <div class="content">${docContent}</div>
          </body>
        </html>
      `

      printWindow.document.write(htmlContent)
      printWindow.document.close()

      setTimeout(() => {
        printWindow.print()
      }, 500)
    } catch (error) {
      console.error('Error exportando PDF:', error)
      alert('Error al exportar PDF')
    } finally {
      setExporting(false)
    }
  }, [currentDocument, content])

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p className="text-gray-600">Cargando...</p>
      </div>
    )
  }

  const handleLogout = () => {
    logout()
    router.push('/')
  }

  return (
    <div className="flex min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
      <Sidebar user={user} onLogout={handleLogout} />
      
      <main className="flex-1 lg:ml-0">
        <div className="container mx-auto px-4 py-6 max-w-4xl">
        {saveSuccess && (
          <div className="mb-6 p-4 bg-green-50 border-l-4 border-green-500 rounded-lg animate-fade-in">
            <p className="text-green-700 font-semibold flex items-center gap-2"><IconCheck /> Documento guardado exitosamente</p>
          </div>
        )}

        <div className="flex justify-between items-start mb-8">
          <div>
            <h1 className="heading-2 mb-2">{currentDocument?.title}</h1>
            <div className="flex gap-4 text-sm text-gray-600 flex-wrap">
              <span className="flex items-center gap-2"><IconDocument /> {currentDocument?.type}</span>
              <span className="flex items-center gap-2"><IconCheck /> {currentDocument?.status}</span>
              <span className="flex items-center gap-2"><IconVersion /> v{currentDocument?.version}</span>
              <span className="flex items-center gap-2"><IconDate /> {new Date(currentDocument?.created_at).toLocaleDateString()}</span>
            </div>
          </div>
          <div className="flex gap-2 flex-wrap justify-end">
            {!editing && (
              <>
                <button
                  onClick={() => setEditing(true)}
                  className="btn btn-primary gap-2"
                >
                  <IconEdit /> Editar
                </button>
                <button
                  onClick={handleExportPDF}
                  disabled={exporting}
                  className="btn btn-secondary gap-2 disabled:opacity-50"
                >
                  {exporting ? <IconLoading /> : <IconDownload />} {exporting ? 'Exportando...' : 'Exportar PDF'}
                </button>
              </>
            )}
            {editing && (
              <>
                <button
                  onClick={handleSave}
                  disabled={saving}
                  className="btn btn-primary gap-2 disabled:opacity-50"
                >
                  {saving ? <IconLoading /> : <IconSave />} {saving ? 'Guardando...' : 'Guardar Cambios'}
                </button>
                <button
                  onClick={() => {
                    setEditing(false)
                    setContent(currentDocument?.content || '')
                  }}
                  disabled={saving}
                  className="btn btn-secondary disabled:opacity-50"
                >
                  Cancelar
                </button>
              </>
            )}
          </div>
        </div>

        <div className="card">
          <div className="mb-6">
            {editing ? (
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-semibold text-gray-900 mb-2">Editar Contenido del Documento</label>
                  <textarea
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                    className="w-full h-96 p-4 border-2 border-gray-200 rounded-xl font-mono text-sm focus:border-indigo-600 focus:outline-none"
                  />
                </div>
              </div>
            ) : (
              <div className="bg-white p-8 border-2 border-gray-200 rounded-xl">
                <div className="whitespace-pre-wrap text-sm text-gray-700 leading-relaxed font-serif">{content}</div>
              </div>
            )}
          </div>

          {!editing && (
            <div className="flex gap-3 pt-6 border-t border-gray-200">
              <button
                onClick={() => navigator.clipboard.writeText(content)}
                className="btn btn-secondary gap-2 flex-1"
              >
                <IconDocument />
                Copiar al Portapapeles
              </button>
              <button
                onClick={handleDelete}
                className="btn bg-red-600 text-white hover:bg-red-700 gap-2"
              >
                <IconDelete />
                Eliminar Documento
              </button>
            </div>
          )}
        </div>
        </div>
      </main>
    </div>
  )
}
