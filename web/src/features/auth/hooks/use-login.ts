'use client'

import { useMutation } from '@tanstack/react-query'
import { saveTokens } from '@/lib/api-client'
import { login } from '../api'
import type { LoginRequest } from '../types'

export function useLogin() {
    return useMutation({
        mutationFn: (data: LoginRequest) => login(data),
        onSuccess(data) {
            saveTokens(data.access_token, data.refresh_token)
            window.location.replace('/dashboard')
        },
    })
}
