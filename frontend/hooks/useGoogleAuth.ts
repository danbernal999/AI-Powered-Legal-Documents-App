import { useState } from 'react'
import { signInWithPopup, auth, googleProvider } from '@/lib/firebase'
import { useAuthStore } from '@/lib/store'
import { authAPI } from '@/lib/api'
import { useRouter } from 'next/navigation'

export const useGoogleAuth = () => {
    const [isLoading, setIsLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)
    const { login } = useAuthStore()
    const router = useRouter()

    const signInWithGoogle = async () => {
        setIsLoading(true)
        setError(null)

        try {
            // Sign in with Google popup
            const result = await signInWithPopup(auth, googleProvider)
            const user = result.user

            // Get the ID token from Firebase
            const idToken = await user.getIdToken()

            // Send the ID token to your backend for verification
            const response = await authAPI.googleLogin(idToken)
            const { token, user: userData } = response.data

            // Store the token and user data
            login(token, userData)

            // Redirect to dashboard
            router.push('/dashboard')
        } catch (err: any) {
            console.error('Google sign-in error:', err)

            // Handle specific Firebase errors
            if (err.code === 'auth/popup-closed-by-user') {
                setError('Sign-in cancelled')
            } else if (err.code === 'auth/popup-blocked') {
                setError('Popup blocked. Please allow popups for this site.')
            } else if (err.response?.data?.message) {
                setError(err.response.data.message)
            } else {
                setError('Failed to sign in with Google. Please try again.')
            }
        } finally {
            setIsLoading(false)
        }
    }

    return {
        signInWithGoogle,
        isLoading,
        error,
    }
}
