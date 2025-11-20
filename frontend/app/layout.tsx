import type { Metadata } from 'next'
import './globals.css'
import 'remixicon/fonts/remixicon.css'

export const metadata: Metadata = {
  title: 'KiraDoc - Legal Document Generator',
  description: 'Generate professional legal documents with AI',
  icons: {
    // 1. Icono principal y de mayor calidad (para pestañas y marcadores)
    icon: '/Icono-removebg-preview.png', 

    // 2. Icono de Apple Touch (para pantallas de inicio en iOS)
    apple: '/Icono-removebg-preview.png',

    // 3. Opcional: El favicon tradicional (generalmente 16x16 o 32x32)
    shortcut: '/favicon.ico', 
  }
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
      </head>
      <body>
        {children}
      </body>
    </html>
  )
}
