'use client'

import { useMutation } from '@tanstack/react-query'
import { saveTokens } from '@/lib/api-client'
import { register } from '../api'
import type { RegisterRequest } from '../types'

export function useRegister() {
    return useMutation({
        mutationFn: (data: RegisterRequest) => register(data),
        onSuccess(data) {
            saveTokens(data.access_token, data.refresh_token)
            window.location.replace('/onboarding')
        },
    })
}
