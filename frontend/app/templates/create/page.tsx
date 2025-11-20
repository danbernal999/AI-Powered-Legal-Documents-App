'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import Sidebar from '@/components/Sidebar'
import { useAuthStore } from '@/lib/store'
import { templateAPI } from '@/lib/api'
import { useTranslation } from '@/hooks/useTranslation'
import { IconAdd, IconDelete, IconSave, IconLoading } from '@/components/Icons'

export default function CreateTemplatePage() {
  const router = useRouter()
  const { token, user, logout } = useAuthStore()
  const { t } = useTranslation()
  const [loading, setLoading] = useState(false)
  const [saving, setSaving] = useState(false)
  const [templateName, setTemplateName] = useState('')
  const [templateDescription, setTemplateDescription] = useState('')
  const [templateType, setTemplateType] = useState('custom')
  const [templateContent, setTemplateContent] = useState('')
  const [variables, setVariables] = useState<Array<{ name: string; description: string; type: string; required: boolean }>>([])
  const [newVarName, setNewVarName] = useState('')
  const [newVarDescription, setNewVarDescription] = useState('')
  const [newVarType, setNewVarType] = useState('string')
  const [newVarRequired, setNewVarRequired] = useState(true)
  const [saveSuccess, setSaveSuccess] = useState(false)
  const [errors, setErrors] = useState<Record<string, string>>({})

  useEffect(() => {
    if (!token) {
      router.push('/login')
      return
    }
  }, [token, router])

  const handleLogout = () => {
    logout()
    router.push('/')
  }

  const addVariable = () => {
    const newErrors = { ...errors }
    delete newErrors.newVarName
    
    if (!newVarName.trim()) {
      setErrors({ ...newErrors, newVarName: t('templates.variableNameRequired') })
      return
    }

    if (variables.some(v => v.name.toLowerCase() === newVarName.toLowerCase())) {
      setErrors({ ...newErrors, newVarName: t('templates.variableAlreadyExists') })
      return
    }
    
    setVariables([
      ...variables,
      {
        name: newVarName,
        description: newVarDescription,
        type: newVarType,
        required: newVarRequired,
      },
    ])
    
    setNewVarName('')
    setNewVarDescription('')
    setNewVarType('string')
    setNewVarRequired(true)
    setErrors(newErrors)
  }

  const removeVariable = (index: number) => {
    setVariables(variables.filter((_, i) => i !== index))
  }

  const insertVariableInContent = (varName: string) => {
    const textarea = document.getElementById('content-textarea') as HTMLTextAreaElement
    if (!textarea) return
    
    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    const placeholder = `{{${varName}}}`
    const newContent = templateContent.substring(0, start) + placeholder + templateContent.substring(end)
    setTemplateContent(newContent)
    
    setTimeout(() => {
      textarea.focus()
      textarea.setSelectionRange(start + placeholder.length, start + placeholder.length)
    }, 0)
  }

  const getPreviewContent = () => {
    let preview = templateContent
    variables.forEach(variable => {
      const regex = new RegExp(`{{${variable.name}}}`, 'g')
      preview = preview.replace(regex, `[${variable.name}]`)
    })
    return preview
  }

  const getCompletionPercentage = () => {
    let completed = 0
    if (templateName.trim()) completed++
    if (templateContent.trim()) completed++
    if (variables.length > 0) completed++
    return Math.round((completed / 3) * 100)
  }

  const isFormValid = () => {
    return templateName.trim() && templateContent.trim()
  }

  const getVariableUsageCount = () => {
    let count = 0
    variables.forEach(variable => {
      const regex = new RegExp(`{{${variable.name}}}`, 'g')
      const matches = templateContent.match(regex)
      if (matches) count += matches.length
    })
    return count
  }

  const handleSaveTemplate = async (e: React.FormEvent) => {
    e.preventDefault()
    const newErrors: Record<string, string> = {}
    
    if (!templateName.trim()) {
      newErrors.name = t('templates.templateNameRequired')
    }
    if (!templateContent.trim()) {
      newErrors.content = t('templates.contentRequired')
    }

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors)
      return
    }

    setSaving(true)
    try {
      await templateAPI.create({
        name: templateName,
        description: templateDescription,
        type: templateType,
        content: templateContent,
        variables,
      })
      
      setSaveSuccess(true)
      setTimeout(() => {
        router.push('/templates')
      }, 1500)
    } catch (error) {
      console.error('Failed to save template:', error)
      setErrors({ submit: t('templates.failedToSaveTemplate') })
      setSaving(false)
    }
  }

  return (
    <div className="flex min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
      <Sidebar user={user} onLogout={handleLogout} />
      
      <main className="flex-1 lg:ml-0">
        <div className="container mx-auto px-4 py-6 max-w-7xl">
          {saveSuccess && (
            <div className="mb-6 p-4 bg-green-50 border-l-4 border-green-500 rounded-lg animate-fade-in">
              <p className="text-green-700 font-semibold flex items-center gap-2">
                <i className="ri-check-line" /> {t('templates.templateCreatedSuccess')}
              </p>
            </div>
          )}

          {errors.submit && (
            <div className="mb-6 p-4 bg-red-50 border-l-4 border-red-500 rounded-lg">
              <p className="text-red-700 font-semibold">{errors.submit}</p>
            </div>
          )}

          <div className="mb-8">
            <h1 className="heading-1 mb-2">{t('templates.createCustomTemplate')}</h1>
            <p className="text-gray-600">{t('templates.createCustomTemplateDesc')}</p>
          </div>

          {/* Progress Bar */}
          <div className="mb-6 bg-white rounded-lg p-4 border border-gray-200">
            <div className="flex items-center justify-between mb-2">
              <p className="text-xs font-semibold text-gray-600">{t('templates.basicInfo')} Setup</p>
              <p className="text-xs font-bold text-indigo-600">{getCompletionPercentage()}%</p>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-2">
              <div
                className="bg-indigo-600 h-2 rounded-full transition-all duration-300"
                style={{ width: `${getCompletionPercentage()}%` }}
              />
            </div>
          </div>

          <form onSubmit={handleSaveTemplate} className="grid lg:grid-cols-3 gap-6">
            {/* Left Column: Basic Info & Content */}
            <div className="lg:col-span-2 space-y-6">
              {/* Basic Information */}
              <div className="card">
                <h2 className="heading-2 mb-6">{t('templates.basicInfo')}</h2>
                
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-semibold text-gray-900 mb-2">
                      {t('templates.templateName')} *
                    </label>
                    <input
                      type="text"
                      value={templateName}
                      onChange={(e) => setTemplateName(e.target.value)}
                      placeholder={t('templates.templateNamePlaceholder')}
                      className={`input-field ${errors.name ? 'border-red-500 focus:border-red-500' : ''}`}
                    />
                    {errors.name && <p className="text-xs text-red-600 mt-1">{errors.name}</p>}
                    <p className="text-xs text-gray-500 mt-1">{t('templates.giveTemplateName')}</p>
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-gray-900 mb-2">
                      {t('templates.templateDescription')}
                    </label>
                    <textarea
                      value={templateDescription}
                      onChange={(e) => setTemplateDescription(e.target.value)}
                      placeholder={t('templates.templateDescriptionPlaceholder')}
                      className="input-field resize-none"
                      rows={3}
                    />
                    <p className="text-xs text-gray-500 mt-1">{t('templates.optionalTemplateDescription')}</p>
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-gray-900 mb-2">
                      {t('templates.templateType')}
                    </label>
                    <select
                      value={templateType}
                      onChange={(e) => setTemplateType(e.target.value)}
                      className="input-field"
                    >
                      <option value="custom">{t('templates.typeCustom')}</option>
                      <option value="NDA">{t('templates.typeNDA')}</option>
                      <option value="Employment">{t('templates.typeEmployment')}</option>
                      <option value="Rental">{t('templates.typeRental')}</option>
                      <option value="Freelance">{t('templates.typeFreelance')}</option>
                    </select>
                  </div>
                </div>
              </div>

              {/* Template Content */}
              <div className="card">
                <h2 className="heading-2 mb-6">{t('templates.content')}</h2>
                
                <div>
                  <div className="flex items-center justify-between mb-2">
                    <label className="block text-sm font-semibold text-gray-900">
                      {t('templates.content')} *
                    </label>
                    {variables.length > 0 && (
                      <p className="text-xs text-gray-500">
                        {t('templates.quickInsert')}:
                      </p>
                    )}
                  </div>
                  
                  {variables.length > 0 && (
                    <div className="flex gap-1 flex-wrap mb-2 pb-2 border-b border-gray-200">
                      {variables.map(variable => (
                        <button
                          key={variable.name}
                          type="button"
                          onClick={() => insertVariableInContent(variable.name)}
                          className="px-2 py-1 text-xs bg-indigo-100 text-indigo-700 rounded hover:bg-indigo-200 transition"
                        >
                          {`{{${variable.name}}}`}
                        </button>
                      ))}
                    </div>
                  )}
                  
                  <textarea
                    id="content-textarea"
                    value={templateContent}
                    onChange={(e) => setTemplateContent(e.target.value)}
                    placeholder={t('templates.contentPlaceholder')}
                    className={`w-full h-64 p-4 border-2 border-gray-200 rounded-xl font-mono text-sm focus:border-indigo-600 focus:outline-none ${errors.content ? 'border-red-500 focus:border-red-500' : ''}`}
                  />
                  {errors.content && <p className="text-xs text-red-600 mt-1">{errors.content}</p>}
                  <p className="text-xs text-gray-500 mt-2">
                    {t('templates.contentHint')}
                  </p>
                </div>
              </div>

              {/* Variables */}
              <div className="card">
                <div className="flex items-center justify-between mb-6">
                  <h2 className="heading-2 flex items-center gap-2 m-0">
                    <i className="ri-variable-line text-indigo-600" />
                    {t('templates.variables')}
                  </h2>
                  {variables.length > 0 && (
                    <span className="px-3 py-1 bg-indigo-100 text-indigo-700 rounded-full text-xs font-bold">
                      {variables.length} {variables.length === 1 ? 'variable' : 'variables'}
                    </span>
                  )}
                </div>
                
                {variables.length > 0 && (
                  <div className="mb-6 space-y-3">
                    <h3 className="text-sm font-semibold text-gray-900 flex items-center gap-2">
                      <i className="ri-check-double-line text-green-600" />
                      {t('templates.addedVariables')}
                    </h3>
                    <div className="grid gap-3">
                      {variables.map((variable, index) => (
                        <div key={index} className="border border-gray-200 rounded-lg p-4 hover:border-indigo-300 hover:shadow-sm transition-all bg-white">
                          <div className="flex items-start justify-between mb-2">
                            <div>
                              <p className="font-semibold text-gray-900 flex items-center gap-2">
                                {variable.name}
                                {variable.required && (
                                  <span className="px-2 py-0.5 bg-red-100 text-red-700 text-xs rounded font-semibold">
                                    Required
                                  </span>
                                )}
                              </p>
                              <p className="text-xs text-gray-600 mt-1">{variable.description || t('templates.noDescription')}</p>
                            </div>
                            <button
                              type="button"
                              onClick={() => removeVariable(index)}
                              className="p-2 text-red-600 hover:bg-red-50 rounded-lg transition"
                            >
                              <IconDelete />
                            </button>
                          </div>
                          <div className="flex gap-2 mt-3">
                            <span className="px-2 py-1 bg-blue-100 text-blue-700 text-xs rounded font-semibold">
                              {variable.type === 'string' ? t('templates.typeString') : variable.type === 'number' ? t('templates.typeNumber') : variable.type === 'date' ? t('templates.typeDate') : variable.type}
                            </span>
                            {templateContent.includes(`{{${variable.name}}}`) && (
                              <span className="px-2 py-1 bg-green-100 text-green-700 text-xs rounded font-semibold flex items-center gap-1">
                                <i className="ri-check-line" /> Used
                              </span>
                            )}
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                <div className="bg-indigo-50 border-2 border-indigo-200 rounded-xl p-4 space-y-4">
                  <h3 className="font-semibold text-indigo-900">{t('templates.addVariable')}</h3>
                  
                  <div className="grid grid-cols-2 gap-3">
                    <div>
                      <label className="block text-xs font-semibold text-gray-900 mb-1">
                        {t('templates.variableName')}
                      </label>
                      <input
                        type="text"
                        value={newVarName}
                        onChange={(e) => setNewVarName(e.target.value)}
                        placeholder={t('templates.variablePlaceholder')}
                        className={`input-field text-sm ${errors.newVarName ? 'border-red-500 focus:border-red-500' : ''}`}
                      />
                      {errors.newVarName && <p className="text-xs text-red-600 mt-1">{errors.newVarName}</p>}
                    </div>
                    <div>
                      <label className="block text-xs font-semibold text-gray-900 mb-1">
                        {t('templates.variableType')}
                      </label>
                      <select
                        value={newVarType}
                        onChange={(e) => setNewVarType(e.target.value)}
                        className="input-field text-sm"
                      >
                        <option value="string">{t('templates.typeString')}</option>
                        <option value="number">{t('templates.typeNumber')}</option>
                        <option value="date">{t('templates.typeDate')}</option>
                      </select>
                    </div>
                  </div>

                  <div>
                    <label className="block text-xs font-semibold text-gray-900 mb-1">
                      {t('templates.variableDescription')}
                    </label>
                    <input
                      type="text"
                      value={newVarDescription}
                      onChange={(e) => setNewVarDescription(e.target.value)}
                      placeholder={t('templates.variableDescriptionPlaceholder')}
                      className="input-field text-sm"
                    />
                  </div>

                  <div className="flex items-center gap-2">
                    <input
                      type="checkbox"
                      id="required"
                      checked={newVarRequired}
                      onChange={(e) => setNewVarRequired(e.target.checked)}
                      className="w-4 h-4 text-indigo-600 rounded"
                    />
                    <label htmlFor="required" className="text-sm text-gray-700">
                      {t('templates.requiredField')}
                    </label>
                  </div>

                  <button
                    type="button"
                    onClick={addVariable}
                    className="w-full btn btn-secondary gap-2"
                  >
                    <IconAdd /> {t('templates.addVariable')}
                  </button>
                </div>
              </div>
            </div>

            {/* Right Column: Preview & Stats */}
            <div className="lg:col-span-1 space-y-6">
              {/* Template Stats */}
              <div className="card bg-gradient-to-br from-indigo-50 to-blue-50 sticky top-24">
                <h3 className="heading-3 mb-4 flex items-center gap-2">
                  <i className="ri-bar-chart-box-line text-indigo-600" />
                  Template Stats
                </h3>
                
                <div className="space-y-3">
                  <div className="bg-white rounded-lg p-3 border border-indigo-200">
                    <p className="text-xs font-semibold text-gray-600 mb-1">Name</p>
                    <p className="font-semibold text-gray-900">{templateName || '—'}</p>
                  </div>

                  <div className="bg-white rounded-lg p-3 border border-indigo-200">
                    <p className="text-xs font-semibold text-gray-600 mb-1">Type</p>
                    <p className="font-semibold text-gray-900 capitalize">{templateType}</p>
                  </div>

                  <div className="bg-white rounded-lg p-3 border border-indigo-200">
                    <p className="text-xs font-semibold text-gray-600 mb-1">Content Length</p>
                    <p className="font-semibold text-gray-900">{templateContent.length} characters</p>
                  </div>

                  <div className="bg-white rounded-lg p-3 border border-indigo-200">
                    <p className="text-xs font-semibold text-gray-600 mb-1">Variables</p>
                    <div className="flex items-center justify-between">
                      <p className="font-semibold text-gray-900">{variables.length} defined</p>
                      {getVariableUsageCount() > 0 && (
                        <span className="text-xs bg-green-100 text-green-700 px-2 py-1 rounded font-semibold">
                          {getVariableUsageCount()} used
                        </span>
                      )}
                    </div>
                  </div>
                </div>
              </div>

              {/* Live Preview */}
              <div className="card bg-gradient-to-br from-blue-50 to-purple-50 sticky top-80">
                <h3 className="heading-3 mb-4 flex items-center gap-2">
                  <i className="ri-eye-line text-purple-600" />
                  {t('templates.preview')}
                </h3>
                
                {templateContent ? (
                  <div className="space-y-4">
                    <div>
                      <p className="text-xs font-semibold text-gray-600 mb-2">{t('templates.content')} Preview</p>
                      <div className="bg-white p-4 rounded-lg border-2 border-gray-200 max-h-64 overflow-y-auto">
                        <p className="text-xs text-gray-700 whitespace-pre-wrap font-mono leading-relaxed">
                          {getPreviewContent()}
                        </p>
                      </div>
                    </div>

                    {variables.length > 0 && (
                      <div>
                        <p className="text-xs font-semibold text-gray-600 mb-2">{t('templates.variables')} ({variables.length})</p>
                        <div className="space-y-2 max-h-48 overflow-y-auto">
                          {variables.map(variable => (
                            <div key={variable.name} className={`bg-white p-3 rounded border-l-4 ${templateContent.includes(`{{${variable.name}}}`) ? 'border-green-500' : 'border-gray-300'}`}>
                              <p className="text-xs font-semibold text-gray-900 flex items-center gap-2">
                                {variable.name}
                                {templateContent.includes(`{{${variable.name}}}`) && (
                                  <i className="ri-check-line text-green-600 text-xs" />
                                )}
                              </p>
                              {variable.description && (
                                <p className="text-xs text-gray-600 mt-1">{variable.description}</p>
                              )}
                              <p className="text-xs text-gray-500 mt-1">
                                Type: <span className="font-semibold">{variable.type === 'string' ? t('templates.typeString') : variable.type === 'number' ? t('templates.typeNumber') : variable.type === 'date' ? t('templates.typeDate') : variable.type}</span>
                              </p>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}
                  </div>
                ) : (
                  <div className="text-center py-8">
                    <div className="text-3xl text-gray-300 mb-2">
                      <i className="ri-file-text-line" />
                    </div>
                    <p className="text-sm text-gray-600">{t('templates.addContentToPreview')}</p>
                  </div>
                )}
              </div>

              {/* Form Validation */}
              <div className="card bg-gradient-to-br from-amber-50 to-orange-50">
                <h3 className="heading-3 mb-3 flex items-center gap-2">
                  <i className="ri-checkbox-multiple-mark-line text-amber-600" />
                  Requirements
                </h3>
                
                <div className="space-y-2">
                  <div className="flex items-center gap-2 text-sm">
                    {templateName.trim() ? (
                      <i className="ri-check-circle-fill text-green-600" />
                    ) : (
                      <i className="ri-circle-line text-gray-400" />
                    )}
                    <span className={templateName.trim() ? 'text-gray-900 font-semibold' : 'text-gray-600'}>
                      Template name
                    </span>
                  </div>

                  <div className="flex items-center gap-2 text-sm">
                    {templateContent.trim() ? (
                      <i className="ri-check-circle-fill text-green-600" />
                    ) : (
                      <i className="ri-circle-line text-gray-400" />
                    )}
                    <span className={templateContent.trim() ? 'text-gray-900 font-semibold' : 'text-gray-600'}>
                      Content added
                    </span>
                  </div>

                  <div className="flex items-center gap-2 text-sm">
                    {variables.length > 0 ? (
                      <i className="ri-check-circle-fill text-green-600" />
                    ) : (
                      <i className="ri-circle-line text-gray-400" />
                    )}
                    <span className={variables.length > 0 ? 'text-gray-900 font-semibold' : 'text-gray-600'}>
                      At least 1 variable (optional)
                    </span>
                  </div>
                </div>

                <div className="flex gap-2 mt-4">
                  <button
                    type="submit"
                    disabled={saving || !isFormValid()}
                    className="flex-1 btn btn-primary gap-2 disabled:opacity-50"
                  >
                    {saving ? <IconLoading /> : <IconSave />}
                    {saving ? t('templates.creatingTemplate') : t('templates.createTemplate')}
                  </button>
                  <button
                    type="button"
                    onClick={() => router.push('/templates')}
                    disabled={saving}
                    className="flex-1 btn btn-secondary disabled:opacity-50"
                  >
                    {t('common.cancel')}
                  </button>
                </div>
              </div>
            </div>
          </form>
        </div>
      </main>
    </div>
  )
}
