'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import Sidebar from '@/components/Sidebar'
import { useAuthStore } from '@/lib/store'
import { useTranslation } from '@/hooks/useTranslation'
import { IconSave, IconLoading } from '@/components/Icons'

export default function SettingsPage() {
  const router = useRouter()
  const { token, user, logout } = useAuthStore()
  const { t } = useTranslation()
  const [loading, setLoading] = useState(true)
  const [email, setEmail] = useState('')
  const [name, setName] = useState('')
  const [saving, setSaving] = useState(false)
  const [saveSuccess, setSaveSuccess] = useState(false)

  useEffect(() => {
    if (!token) {
      router.push('/login')
      return
    }

    if (user) {
      setEmail(user.email || '')
      setName(user.name || '')
    }
    setLoading(false)
  }, [token, router, user])

  const handleLogout = () => {
    logout()
    router.push('/')
  }

  const handleSaveSettings = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)
    try {
      setSaveSuccess(true)
      setTimeout(() => setSaveSuccess(false), 3000)
    } catch (error) {
      console.error('Error saving settings:', error)
    } finally {
      setSaving(false)
    }
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
        <div className="container mx-auto px-4 py-6 max-w-2xl">
          <div className="mb-8">
            <h2 className="heading-2 mb-2 mt-6">{t('settings.title')}</h2>
            <p className="text-gray-600">{t('settings.subtitle')}</p>
          </div>

          {saveSuccess && (
            <div className="mb-6 p-4 bg-green-50 border-l-4 border-green-500 rounded-lg">
              <p className="text-green-700 font-semibold flex items-center gap-2"><IconSave /> {t('common.success')}</p>
            </div>
          )}

          <div className="space-y-6">
            <div className="card">
              <h3 className="heading-3 mb-6">{t('settings.profile')}</h3>
              <form onSubmit={handleSaveSettings} className="space-y-6">
                <div>
                  <label className="block text-sm font-semibold text-gray-900 mb-2">
                    {t('settings.fullName')}
                  </label>
                  <input
                    type="text"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    className="input-field"
                    placeholder={t('settings.fullNamePlaceholder')}
                  />
                </div>

                <div>
                  <label className="block text-sm font-semibold text-gray-900 mb-2">
                    {t('settings.emailAddress')}
                  </label>
                  <input
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    className="input-field"
                    placeholder={t('settings.emailPlaceholder')}
                    disabled
                  />
                  <p className="text-xs text-gray-500 mt-2">{t('settings.emailCannotChange')}</p>
                </div>

                <button
                  type="submit"
                  disabled={saving}
                  className="btn btn-primary disabled:opacity-50 gap-2"
                >
                  {saving ? <IconLoading /> : <IconSave />}
                  {saving ? t('settings.saving') : t('settings.saveChanges')}
                </button>
              </form>
            </div>

            <div className="card">
              <h3 className="heading-3 mb-4">{t('settings.account')}</h3>
              <div className="space-y-4 text-sm">
                <div className="flex justify-between items-center py-2 border-b border-gray-200">
                  <span className="text-gray-600">{t('settings.userId')}</span>
                  <span className="text-gray-900 font-mono">{user?.id}</span>
                </div>
                <div className="flex justify-between items-center py-2 border-b border-gray-200">
                  <span className="text-gray-600">{t('documents.created')}</span>
                  <span className="text-gray-900">{user?.email}</span>
                </div>
                <div className="flex justify-between items-center py-2">
                  <span className="text-gray-600">{t('settings.memberSince')}</span>
                  <span className="text-gray-900">{user?.created_at ? new Date(user.created_at).toLocaleDateString() : t('settings.notAvailable')}</span>
                </div>
              </div>
            </div>

            <div className="card bg-red-50 border-red-200">
              <h3 className="heading-3 mb-4 text-red-900">{t('settings.dangerZone')}</h3>
              <p className="text-sm text-red-700 mb-4">
                {t('settings.logoutConfirm')}
              </p>
              <button
                onClick={handleLogout}
                className="btn bg-red-600 text-white hover:bg-red-700 gap-2"
              >
                <i className="ri-logout-box-line" />
                {t('settings.logout')}
              </button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
