'use client'

import { useTranslation } from '@/hooks/useTranslation'

interface Signature {
  id: string
  signer_name: string
  signer_email: string
  signed_at: string
  timestamp: string
}

interface SignaturesListProps {
  signatures: Signature[]
  onDelete?: (signatureId: string) => void
  isLoading?: boolean
}

export default function SignaturesList({
  signatures,
  onDelete,
  isLoading = false,
}: SignaturesListProps) {
  const { t } = useTranslation()

  if (isLoading) {
    return (
      <div className="space-y-2">
        {[1, 2].map((i) => (
          <div key={i} className="bg-gray-200 rounded h-16 animate-pulse" />
        ))}
      </div>
    )
  }

  if (!signatures || signatures.length === 0) {
    return (
      <div className="text-center py-8">
        <div className="text-3xl text-gray-300 mb-2">
          <i className="ri-edit-2-line" />
        </div>
        <p className="text-sm text-gray-600">No signatures yet. Be the first to sign!</p>
      </div>
    )
  }

  return (
    <div className="space-y-3">
      {signatures.map((signature) => (
        <div
          key={signature.id}
          className="bg-gradient-to-r from-green-50 to-emerald-50 border border-green-200 rounded-lg p-4 hover:shadow-md transition"
        >
          <div className="flex items-start justify-between mb-2">
            <div className="flex-1">
              <div className="flex items-center gap-2 mb-1">
                <i className="ri-check-double-fill text-green-600" />
                <p className="font-semibold text-gray-900">{signature.signer_name}</p>
              </div>
              <p className="text-xs text-gray-600">{signature.signer_email}</p>
            </div>
            {onDelete && (
              <button
                onClick={() => onDelete(signature.id)}
                className="p-2 text-red-600 hover:bg-red-50 rounded-lg transition"
                title="Remove signature"
              >
                <i className="ri-delete-bin-line" />
              </button>
            )}
          </div>
          <div className="flex items-center gap-4 text-xs text-gray-600 pt-2 border-t border-green-200">
            <span className="flex items-center gap-1">
              <i className="ri-calendar-line" />
              {new Date(signature.signed_at).toLocaleDateString()}
            </span>
            <span className="flex items-center gap-1">
              <i className="ri-time-line" />
              {new Date(signature.signed_at).toLocaleTimeString()}
            </span>
          </div>
        </div>
      ))}
    </div>
  )
}
