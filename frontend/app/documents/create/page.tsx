'use client'

import { useState, useEffect, Suspense } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import Link from 'next/link'
import Sidebar from '@/components/Sidebar'
import { useAuthStore, useTemplateStore } from '@/lib/store'
import { templateAPI, generateAPI, documentAPI } from '@/lib/api'
import { useTranslation } from '@/hooks/useTranslation'
import { IconCheck, IconDocumentFile, IconSave, IconLoading, IconEdit } from '@/components/Icons'

function CreateDocumentContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { token, user, logout } = useAuthStore()
  const { templates, setTemplates } = useTemplateStore()
  const { t, interpolate } = useTranslation()
  const [selectedTemplate, setSelectedTemplate] = useState<any>(null)
  const [title, setTitle] = useState('')
  const [variables, setVariables] = useState<Record<string, string>>({})
  const [generatedContent, setGeneratedContent] = useState('')
  const [loading, setLoading] = useState(false)
  const [showPreview, setShowPreview] = useState(false)
  const [saving, setSaving] = useState(false)
  const [saveSuccess, setSaveSuccess] = useState(false)

  const handleSelectTemplate = async (template: any) => {
    setSelectedTemplate(template)
    const templateVars = template.variables?.reduce(
      (acc: any, v: any) => ({ ...acc, [v.name]: '' }),
      {}
    ) || {}
    setVariables(templateVars)
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
      }
    }

    fetchTemplates()
  }, [token, router, setTemplates])

  useEffect(() => {
    if (templates.length === 0) return
    
    const templateId = searchParams.get('templateId')
    if (templateId && !selectedTemplate) {
      const template = templates.find((t: any) => t.id === templateId)
      if (template) {
        setSelectedTemplate(template)
        const templateVars = template.variables?.reduce(
          (acc: any, v: any) => ({ ...acc, [v.name]: '' }),
          {}
        ) || {}
        setVariables(templateVars)
      }
    }
  }, [templates, searchParams, selectedTemplate])

  const handleVariableChange = (name: string, value: string) => {
    setVariables((prev) => ({ ...prev, [name]: value }))
  }

  const handleGenerateDocument = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!selectedTemplate || !title) return

    setLoading(true)
    try {
      const response = await generateAPI.generateDocument(
        selectedTemplate.id,
        variables,
        title
      )
      setGeneratedContent(response.data.content)
      setShowPreview(true)
    } catch (error) {
      console.error('Failed to generate document:', error)
      alert('Failed to generate document')
    } finally {
      setLoading(false)
    }
  }

  const handleSaveDocument = async () => {
    if (!selectedTemplate || !title || !generatedContent) return

    setSaving(true)
    try {
      const response = await documentAPI.create({
        template_id: selectedTemplate.id,
        title,
        type: selectedTemplate.type,
        content: generatedContent,
        variables,
      })
      setSaveSuccess(true)
      setTimeout(() => {
        router.push(`/documents/${response.data.id}`)
      }, 1500)
    } catch (error) {
      console.error('Failed to save document:', error)
      alert('Failed to save document')
      setSaving(false)
    }
  }

  const handleLogout = () => {
    logout()
    router.push('/')
  }

  const getCompletionPercentage = () => {
    if (!selectedTemplate?.variables) return 0
    const filled = Object.values(variables).filter(v => v && String(v).trim()).length
    return Math.round((filled / selectedTemplate.variables.length) * 100)
  }

  const getRequiredFields = () => {
    return selectedTemplate?.variables?.filter((v: any) => v.required) || []
  }

  const getFilledRequiredFields = () => {
    const required = getRequiredFields()
    return required.filter((v: any) => variables[v.name] && String(variables[v.name]).trim())
  }

  const isFormValid = () => {
    const required = getRequiredFields()
    return title.trim() && required.every((v: any) => variables[v.name] && String(variables[v.name]).trim())
  }

  const getLivePreview = () => {
    let preview = generatedContent || selectedTemplate?.content || ''
    Object.entries(variables).forEach(([key, value]) => {
      const regex = new RegExp(`{{${key}}}`, 'g')
      preview = preview.replace(regex, String(value || `[${key}]`))
    })
    return preview
  }

  return (
    <div className="flex min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
      <Sidebar user={user} onLogout={handleLogout} />
      
      <main className="flex-1 lg:ml-0">
        <div className="container mx-auto px-4 py-6 max-w-7xl">
          {saveSuccess && (
            <div className="mb-6 p-4 bg-green-50 border-l-4 border-green-500 rounded-lg animate-fade-in">
              <p className="text-green-700 font-semibold flex items-center gap-2"><IconCheck /> {t('documents.documentSavedSuccess')}</p>
            </div>
          )}

          <div className="mb-8">
            <h1 className="heading-1 mb-2">{t('documents.createNew')}</h1>
            <p className="text-gray-600">{selectedTemplate ? `Using: ${selectedTemplate.name}` : t('documents.chooseTemplate')}</p>
          </div>

          {!selectedTemplate ? (
            <div className="space-y-4">
              <div className="mb-6">
                <h2 className="heading-2 mb-4">{t('documents.selectTemplate')}</h2>
                <p className="text-gray-600">{t('documents.chooseTemplate')}</p>
              </div>
              {templates.length === 0 ? (
                <div className="card text-center py-12">
                  <p className="text-gray-600">{t('documents.loadingTemplates')}</p>
                </div>
              ) : (
                <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
                  {templates.map((template: any) => (
                    <button
                      key={template.id}
                      onClick={() => handleSelectTemplate(template)}
                      className="card text-left hover:shadow-lg hover:border-indigo-300 transition-all"
                    >
                      <div className="flex flex-col gap-3 h-full">
                        <div className="text-4xl text-indigo-500"><IconDocumentFile /></div>
                        <div>
                          <h3 className="font-semibold text-lg text-gray-900">{template.name}</h3>
                          <p className="text-sm text-gray-600 mt-1 line-clamp-2">{template.description}</p>
                        </div>
                        <div className="flex-1" />
                        <span className="text-indigo-600 text-sm font-semibold flex items-center gap-1">
                          {t('documents.selectTemplate')} →
                        </span>
                      </div>
                    </button>
                  ))}
                </div>
              )}
            </div>
          ) : (
            <div className="grid lg:grid-cols-3 gap-6">
              {/* Left: Form */}
              <div className="lg:col-span-2">
                <button
                  onClick={() => {
                    setSelectedTemplate(null)
                    setVariables({})
                  }}
                  className="mb-6 btn btn-secondary gap-2"
                >
                  <span>←</span> {t('documents.changeTemplate')}
                </button>

                <form onSubmit={handleGenerateDocument} className="space-y-6">
                  {/* Document Title */}
                  <div className="card">
                    <h3 className="heading-3 mb-4 flex items-center gap-2">
                      <i className="ri-file-text-line text-indigo-600" />
                      {t('documents.documentTitle')}
                    </h3>
                    <input
                      type="text"
                      value={title}
                      onChange={(e) => setTitle(e.target.value)}
                      placeholder={t('documents.documentTitlePlaceholder')}
                      className="input-field mb-2"
                      required
                    />
                    <p className="text-xs text-gray-500">{t('documents.giveDocumentName')}</p>
                  </div>

                  {/* Variables */}
                  {selectedTemplate?.variables && selectedTemplate.variables.length > 0 && (
                    <div className="card">
                      <div className="flex items-center justify-between mb-6">
                        <h3 className="heading-3 flex items-center gap-2 m-0">
                          <i className="ri-list-check-line text-indigo-600" />
                          {t('documents.fillDetails')}
                        </h3>
                        <div className="text-xs font-semibold text-gray-600">
                          {getFilledRequiredFields().length} / {getRequiredFields().length} {t('templates.required')}
                        </div>
                      </div>

                      {/* Progress Bar */}
                      <div className="mb-6">
                        <div className="flex items-center justify-between mb-2">
                          <p className="text-xs font-semibold text-gray-600">Completion</p>
                          <p className="text-xs font-bold text-indigo-600">{getCompletionPercentage()}%</p>
                        </div>
                        <div className="w-full bg-gray-200 rounded-full h-2">
                          <div
                            className="bg-indigo-600 h-2 rounded-full transition-all duration-300"
                            style={{ width: `${getCompletionPercentage()}%` }}
                          />
                        </div>
                      </div>

                      <div className="space-y-5">
                        {selectedTemplate.variables.map((variable: any) => (
                          <div key={variable.name} className="border border-gray-200 rounded-lg p-4 hover:border-indigo-300 transition-colors">
                            <div className="flex items-start justify-between mb-2">
                              <label className="block text-sm font-semibold text-gray-900">
                                {variable.name}
                                {variable.required && <span className="text-red-500 ml-1">*</span>}
                              </label>
                              {variables[variable.name] && String(variables[variable.name]).trim() && (
                                <span className="text-green-600 text-xs font-bold">✓</span>
                              )}
                            </div>
                            <p className="text-xs text-gray-600 mb-3">{variable.description}</p>
                            <input
                              type={variable.type === 'number' ? 'number' : variable.type === 'date' ? 'date' : 'text'}
                              value={variables[variable.name] || ''}
                              onChange={(e) => handleVariableChange(variable.name, e.target.value)}
                              placeholder={interpolate(t('documents.enterVariable'), { variable: variable.name.toLowerCase() })}
                              className={`input-field ${variables[variable.name] && String(variables[variable.name]).trim() ? 'border-green-300 bg-green-50' : ''}`}
                              required={variable.required}
                            />
                          </div>
                        ))}
                      </div>
                    </div>
                  )}

                  {/* Generate Button */}
                  <button
                    type="submit"
                    disabled={loading || !isFormValid()}
                    className="w-full btn btn-primary disabled:opacity-50 gap-2 py-3"
                  >
                    <i className={`ri-sparkles-line ${loading ? 'animate-spin' : ''}`} />
                    {loading ? t('documents.generatingDocument') : t('documents.generateDocument')}
                  </button>
                </form>
              </div>

              {/* Right: Template Info + Preview */}
              <div className="space-y-6">
                {/* Template Info Card */}
                <div className="card bg-gradient-to-br from-indigo-50 to-blue-50 sticky top-24">
                  <div className="flex items-start gap-3 mb-4">
                    <div className="text-3xl text-indigo-600"><IconDocumentFile /></div>
                    <div>
                      <h3 className="font-bold text-gray-900">{selectedTemplate.name}</h3>
                      <p className="text-xs text-gray-600">{selectedTemplate.type}</p>
                    </div>
                  </div>
                  {selectedTemplate.description && (
                    <p className="text-sm text-gray-700 mb-4">{selectedTemplate.description}</p>
                  )}
                  <div className="pt-4 border-t border-indigo-200">
                    <p className="text-xs font-semibold text-indigo-900 mb-2">Template Variables: {selectedTemplate.variables?.length || 0}</p>
                  </div>
                </div>

                {/* Live Preview */}
                <div className="card bg-gradient-to-br from-blue-50 to-purple-50 sticky top-96">
                  <div className="flex items-center gap-2 mb-4">
                    <i className="ri-eye-line text-lg text-purple-600" />
                    <h3 className="font-bold text-gray-900">{t('documents.documentPreview')}</h3>
                  </div>

                  {showPreview && generatedContent ? (
                    <div className="space-y-4">
                      <div className="bg-white p-4 rounded-lg border-2 border-gray-200 max-h-64 overflow-y-auto">
                        <p className="text-xs text-gray-700 whitespace-pre-wrap font-mono leading-relaxed">
                          {generatedContent}
                        </p>
                      </div>
                      <div className="flex gap-2">
                        <button 
                          onClick={handleSaveDocument}
                          disabled={saving}
                          className="flex-1 btn btn-primary gap-2 disabled:opacity-50 text-sm py-2"
                        >
                          {saving ? <IconLoading /> : <IconSave />}
                          {saving ? t('documents.saving') : t('documents.saveDocument')}
                        </button>
                        <button
                          onClick={() => setShowPreview(false)}
                          disabled={saving}
                          className="flex-1 btn btn-secondary disabled:opacity-50 text-sm py-2"
                        >
                          {t('common.cancel')}
                        </button>
                      </div>
                    </div>
                  ) : (
                    <div className="text-center py-8">
                      <div className="text-3xl text-gray-300 mb-2">
                        <i className="ri-file-text-line" />
                      </div>
                      <p className="text-sm text-gray-600">{t('documents.fillDetailsPreview')}</p>
                      {isFormValid() && (
                        <p className="text-xs text-indigo-600 font-semibold mt-2">✓ Ready to generate</p>
                      )}
                    </div>
                  )}
                </div>
              </div>
            </div>
          )}
        </div>
      </main>
    </div>
  )
}

export default function CreateDocumentPage() {
  return (
    <Suspense>
      <CreateDocumentContent />
    </Suspense>
  )
}
