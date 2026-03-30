import { apiClient, publicClient } from '@/lib/api-client'
import type {
    RegisterRequest,
    RegisterResponse,
    CreateOrganizationRequest,
    CreateOrganizationResponse,
    LoginRequest,
    LoginResponse,
} from './types'

export async function register(data: RegisterRequest): Promise<RegisterResponse> {
    return publicClient.post('api/auth/register', { json: data }).json<RegisterResponse>()
}

export async function createOrganization(
    data: CreateOrganizationRequest,
): Promise<CreateOrganizationResponse> {
    return apiClient.post('api/v1/organizations', { json: data }).json<CreateOrganizationResponse>()
}

export async function completeOnboarding(): Promise<void> {
    await apiClient.post('api/v1/organizations/me/onboarding')
}

export async function login(data: LoginRequest): Promise<LoginResponse> {
    return publicClient.post('api/auth/login', { json: data }).json<LoginResponse>()
}
