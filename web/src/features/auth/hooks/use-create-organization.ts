'use client'

import { useMutation } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'
import { saveTokens, getRefreshToken } from '@/lib/api-client'
import { createOrganization } from '../api'
import type { CreateOrganizationRequest } from '../types'

export function useCreateOrganization(onComplete?: () => void) {
    const router = useRouter()

    return useMutation({
        mutationFn: async (data: CreateOrganizationRequest) => {
            const org = await createOrganization(data)
            // Swap in the new access token that carries org_id before calling onboarding
            saveTokens(org.access_token, getRefreshToken() ?? '')
            return org
        },
        onSuccess() {
            if (onComplete) {
                onComplete()
            } else {
                router.push('/dashboard')
            }
        },
    })
}
