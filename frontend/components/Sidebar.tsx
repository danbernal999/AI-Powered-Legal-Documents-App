'use client'

import { useState } from 'react'
import { usePathname } from 'next/navigation'
import Link from 'next/link'
import { useTranslation } from '@/hooks/useTranslation'
import { useLanguage, Language } from '@/lib/i18n'
import { IconDashboard, IconDocument, IconTemplates, IconSettings } from '@/components/Icons'

interface SidebarProps {
  user?: {
    name?: string
    email?: string
  } | null
  onLogout: () => void
}

export default function Sidebar({ user, onLogout }: SidebarProps) {
  const [isOpen, setIsOpen] = useState(false)
  const [showLanguageMenu, setShowLanguageMenu] = useState(false)
  const pathname = usePathname()
  const { t } = useTranslation()
  const { language, setLanguage } = useLanguage()

  const closeSidebar = () => setIsOpen(false)

  const isActive = (path: string) => pathname === path

  const languages: { code: Language; name: string; flag: string }[] = [
    { code: 'en', name: 'English', flag: '🇺🇸' },
    { code: 'es', name: 'Español', flag: '🇪🇸' },
    { code: 'pt', name: 'Português', flag: '🇵🇹' },
  ]

  const navItems = [
    { href: '/dashboard', label: t('sidebar.dashboard'), icon: IconDashboard },
    { href: '/documents', label: t('sidebar.documents'), icon: IconDocument },
    { href: '/templates', label: t('sidebar.templates'), icon: IconTemplates },
    { href: '/settings', label: t('sidebar.settings'), icon: IconSettings },
  ]

  return (
    <>
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="fixed lg:hidden top-4 left-4 z-40 p-2 rounded-lg bg-slate-600 text-white hover:bg-slate-700 transition"
        aria-label="Toggle menu"
      >
        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d={isOpen ? "M6 18L18 6M6 6l12 12" : "M4 6h16M4 12h16M4 18h16"} />
        </svg>
      </button>

      <div className={`fixed inset-y-0 left-0 w-64 bg-white border-r border-gray-200 z-30 transform transition-transform duration-300 ${isOpen ? 'translate-x-0' : '-translate-x-full'} lg:translate-x-0 lg:sticky lg:top-0`}>
        <div className="flex flex-col h-screen">
          <div className="p-6 border-b border-gray-200">
            <div className="flex items-center gap-3">
              <div className="w-12 h-12 bg-slate-300 rounded-full flex items-center justify-center">
                <img src="/LogoDashboard.png" alt="Legal app icon" className="w-11 h-11" />
              </div>
              <div>
                <h1 className="text-base font-bold text-slate-900">Legal</h1>
                <p className="text-xs text-gray-600">Contract Generator</p>
              </div>
            </div>
          </div>

          <nav className="flex-1 overflow-y-auto px-3 py-4 space-y-1">
            {navItems.map((item) => {
              const IconComponent = item.icon
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  onClick={closeSidebar}
                  className={`flex items-center gap-3 px-3 py-3 rounded-lg transition font-medium text-sm ${
                    isActive(item.href)
                      ? 'bg-gray-100 text-slate-700'
                      : 'text-gray-700 hover:bg-gray-50'
                  }`}
                >
                  <span className="text-lg"><IconComponent /></span>
                  {item.label}
                </Link>
              )
            })}
          </nav>

          <div className="border-t border-gray-200 p-4 space-y-3">
            <div className="flex items-center gap-2 px-2">
              <div className="w-8 h-8 bg-orange-200 rounded-full flex items-center justify-center text-orange-700 font-semibold flex-shrink-0 text-sm">
                {user?.name?.charAt(0).toUpperCase()}
              </div>
              <div className="min-w-0 flex-1">
                <p className="text-xs font-semibold text-gray-900 truncate">{user?.name}</p>
                <p className="text-xs text-gray-500 truncate">{user?.email}</p>
              </div>
            </div>

            <div className="relative">
              <button
                onClick={() => setShowLanguageMenu(!showLanguageMenu)}
                className="w-full flex items-center justify-between px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-100 rounded-lg transition"
              >
                <span>
                  {languages.find((l) => l.code === language)?.flag}{' '}
                  {t('sidebar.language')}
                </span>
                <svg
                  className={`w-4 h-4 transition ${
                    showLanguageMenu ? 'rotate-180' : ''
                  }`}
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M19 14l-7 7m0 0l-7-7m7 7V3"
                  />
                </svg>
              </button>

              {showLanguageMenu && (
                <div className="absolute bottom-12 left-0 right-0 bg-white border border-gray-200 rounded-lg shadow-lg overflow-hidden z-50">
                  {languages.map((lang) => (
                    <button
                      key={lang.code}
                      onClick={() => {
                        setLanguage(lang.code)
                        setShowLanguageMenu(false)
                      }}
                      className={`w-full flex items-center gap-2 px-3 py-2 text-sm transition ${
                        language === lang.code
                          ? 'bg-slate-100 text-slate-700 font-semibold'
                          : 'text-gray-700 hover:bg-gray-50'
                      }`}
                    >
                      <span>{lang.flag}</span>
                      {lang.name}
                    </button>
                  ))}
                </div>
              )}
            </div>

            <button
              onClick={() => {
                onLogout()
                closeSidebar()
              }}
              className="w-full px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-100 rounded-lg transition"
            >
              {t('sidebar.logout')}
            </button>
          </div>
        </div>
      </div>

      {isOpen && (
        <div
          className="fixed inset-0 bg-black/30 z-20 lg:hidden"
          onClick={() => setIsOpen(false)}
        />
      )}
    </>
  )
}
