'use client'

import { useRef, useEffect } from 'react'

interface SignatureCanvasProps {
  onSave: (signatureData: string) => void
  onCancel: () => void
  signerName: string
  signerEmail: string
}

export default function SignatureCanvas({
  onSave,
  onCancel,
  signerName,
  signerEmail,
}: SignatureCanvasProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const isDrawing = useRef(false)
  const lastX = useRef(0)
  const lastY = useRef(0)

  useEffect(() => {
    const canvas = canvasRef.current
    if (!canvas) return

    const rect = canvas.getBoundingClientRect()
    canvas.width = rect.width * window.devicePixelRatio
    canvas.height = rect.height * window.devicePixelRatio

    const ctx = canvas.getContext('2d')
    if (ctx) {
      ctx.scale(window.devicePixelRatio, window.devicePixelRatio)
      ctx.lineCap = 'round'
      ctx.lineJoin = 'round'
      ctx.lineWidth = 2
      ctx.strokeStyle = '#1e293b'
      ctx.fillStyle = 'white'
      ctx.fillRect(0, 0, rect.width, rect.height)
    }
  }, [])

  const startDrawing = (e: React.MouseEvent<HTMLCanvasElement>) => {
    const canvas = canvasRef.current
    if (!canvas) return

    const rect = canvas.getBoundingClientRect()
    lastX.current = e.clientX - rect.left
    lastY.current = e.clientY - rect.top
    isDrawing.current = true
  }

  const draw = (e: React.MouseEvent<HTMLCanvasElement>) => {
    if (!isDrawing.current) return

    const canvas = canvasRef.current
    if (!canvas) return

    const ctx = canvas.getContext('2d')
    if (!ctx) return

    const rect = canvas.getBoundingClientRect()
    const x = e.clientX - rect.left
    const y = e.clientY - rect.top

    ctx.beginPath()
    ctx.moveTo(lastX.current, lastY.current)
    ctx.lineTo(x, y)
    ctx.stroke()

    lastX.current = x
    lastY.current = y
  }

  const endDrawing = () => {
    isDrawing.current = false
  }

  const clearCanvas = () => {
    const canvas = canvasRef.current
    if (!canvas) return

    const ctx = canvas.getContext('2d')
    const rect = canvas.getBoundingClientRect()
    if (ctx) {
      ctx.fillStyle = 'white'
      ctx.fillRect(0, 0, rect.width, rect.height)
    }
  }

  const handleSave = () => {
    const canvas = canvasRef.current
    if (!canvas) return

    const signatureData = canvas.toDataURL('image/png')
    onSave(signatureData)
  }

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-white rounded-xl shadow-2xl max-w-2xl w-full mx-4">
        <div className="p-6 border-b border-gray-200">
          <h2 className="heading-2 mb-2">Sign Document</h2>
          <p className="text-sm text-gray-600">
            <strong>{signerName}</strong> ({signerEmail})
          </p>
        </div>

        <div className="p-6">
          <div className="mb-4">
            <label className="block text-sm font-semibold text-gray-900 mb-2">
              Your Signature *
            </label>
            <div className="border-2 border-dashed border-gray-300 rounded-lg overflow-hidden bg-white">
              <canvas
                ref={canvasRef}
                onMouseDown={startDrawing}
                onMouseMove={draw}
                onMouseUp={endDrawing}
                onMouseLeave={endDrawing}
                className="w-full h-48 cursor-crosshair bg-white"
                style={{ display: 'block' }}
              />
            </div>
            <p className="text-xs text-gray-500 mt-2">
              Draw your signature above. Be consistent for better verification.
            </p>
          </div>

          <div className="flex gap-3">
            <button
              type="button"
              onClick={clearCanvas}
              className="flex-1 px-4 py-2 text-sm font-semibold text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition"
            >
              <i className="ri-refresh-line mr-2" />
              Clear
            </button>
            <button
              type="button"
              onClick={onCancel}
              className="flex-1 px-4 py-2 text-sm font-semibold text-gray-700 border border-gray-300 rounded-lg hover:bg-gray-50 transition"
            >
              Cancel
            </button>
            <button
              type="button"
              onClick={handleSave}
              className="flex-1 px-4 py-2 text-sm font-semibold text-white bg-indigo-600 rounded-lg hover:bg-indigo-700 transition"
            >
              <i className="ri-check-line mr-2" />
              Sign Document
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
