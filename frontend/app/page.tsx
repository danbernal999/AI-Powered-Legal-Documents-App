'use client'

/* LANDING PAGE */

import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useEffect } from 'react'
import { useAuthStore } from '@/lib/store'
import { IconBolt, IconRobot, IconTemplates } from '@/components/Icons'

export default function Home() {
  const router = useRouter()
  const { token } = useAuthStore()

  useEffect(() => {
    if (token) {
      router.push('/dashboard')
    }
  }, [token, router])

  return (
    <main className="min-h-screen flex flex-col">
      <nav className="bg-white/80 backdrop-blur-md sticky top-0 z-50 border-b border-gray-200/50">
        <div className="container mx-auto px-4 py-4 flex justify-between items-center">
          <div className="flex items-center gap-2">
            <div className="w-10 h-10 bg-gradient-to-r from-indigo-600 to-purple-600 rounded-xl flex items-center justify-center">
              <span className="text-white font-bold text-lg">K</span>
            </div>
            <h1 className="text-2xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">KiraDoc</h1>
          </div>
          <div className="space-x-3">
            <Link href="/login" className="btn btn-secondary">
              Login
            </Link>
            <Link href="/register" className="btn btn-primary">
              Get Started
            </Link>
          </div>
        </div>
      </nav>

      <section className="flex-1 container mx-auto px-4 py-24">
        <div className="text-center max-w-3xl mx-auto">
          <h2 className="heading-1 mb-6 text-6xl">Generate Professional Legal Documents</h2>
          <p className="text-xl text-gray-600 mb-12 leading-relaxed">
            KiraDoc uses artificial intelligence to help you create customized legal contracts in minutes. 
            From NDAs to employment agreements, we've got you covered.
          </p>

          <div className="grid md:grid-cols-3 gap-6 mb-16">
            <div className="card text-center">
              <div className="text-4xl mb-4"><IconBolt /></div>
              <h3 className="heading-3 mb-3">Lightning Fast</h3>
              <p className="text-gray-600">Generate professional legal documents in seconds, not hours.</p>
            </div>
            <div className="card text-center">
              <div className="text-4xl mb-4"><IconRobot /></div>
              <h3 className="heading-3 mb-3">AI-Powered</h3>
              <p className="text-gray-600">Leverages advanced AI models to ensure accuracy and completeness.</p>
            </div>
            <div className="card text-center">
              <div className="text-4xl mb-4"><IconTemplates /></div>
              <h3 className="heading-3 mb-3">Fully Customizable</h3>
              <p className="text-gray-600">Fully editable templates tailored to your specific needs.</p>
            </div>
          </div>

          <Link href="/register" className="btn btn-primary text-lg px-8 py-4">
            Get Started For Free
          </Link>
        </div>
      </section>
      <section className="py-20 px-8 bg-gray-50">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 max-w-6xl mx-auto items-center">
          <div className="col-span-1 w-1/2 mx-auto lg:w-10/12">
            <h2 className="heading-2 mb-4 text-3xl">Mobile Convenience</h2>
            <p className="text-xl text-gray-600 mb-8">
              Access your documents and templates on the go.
            </p>
            <div className="grid gap-4 grid-cols-2">
              <div className="card">
                <h4 className="font-semibold mb-2">Feature 1</h4>
                <p className="text-gray-600 text-sm">Description</p>
              </div>
              <div className="card">
                <h4 className="font-semibold mb-2">Feature 2</h4>
                <p className="text-gray-600 text-sm">Description</p>
              </div>
            </div>
          </div>
          <div className="col-span-1 mx-auto">
            <img src="/IphoneBackgroundNoFondo.png" alt="iphone" className="w-full" />
          </div>
        </div>
      </section>

      <footer className="bg-gradient-to-r from-gray-900 to-gray-800 text-white py-12 border-t border-gray-700">
        <div className="container mx-auto px-4">
          <div className="grid md:grid-cols-3 gap-8 mb-8">
            <div>
              <h4 className="font-semibold text-lg mb-4">KiraDoc</h4>
              <p className="text-gray-400">Professional legal document generation with AI.</p>
            </div>
            <div>
              <h4 className="font-semibold text-lg mb-4">Features</h4>
              <ul className="space-y-2 text-gray-400">
                <li><a href="#" className="hover:text-white transition">Templates</a></li>
                <li><a href="#" className="hover:text-white transition">AI Generation</a></li>
                <li><a href="#" className="hover:text-white transition">Export</a></li>
              </ul>
            </div>
            <div>
              <h4 className="font-semibold text-lg mb-4">Legal</h4>
              <ul className="space-y-2 text-gray-400">
                <li><a href="#" className="hover:text-white transition">Privacy</a></li>
                <li><a href="#" className="hover:text-white transition">Terms</a></li>
                <li><a href="#" className="hover:text-white transition">Contact</a></li>
              </ul>
            </div>
          </div>
          <div className="border-t border-gray-700 pt-8 text-center text-gray-400">
            <p>&copy; 2024 KiraDoc. All rights reserved.</p>
          </div>
        </div>
      </footer>
    </main>
  )
}
