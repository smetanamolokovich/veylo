import ky, { type KyInstance } from 'ky'
import { getCookieValue, setCookieValue, deleteCookieValue } from '@/lib/auth'

export const BASE_URL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080'

export const publicClient = ky.create({ prefixUrl: BASE_URL })

function getAccessToken(): string | null {
    return getCookieValue('access_token')
}

export function getRefreshToken(): string | null {
    return getCookieValue('refresh_token')
}

export function saveTokens(accessToken: string, refreshToken: string) {
    setCookieValue('access_token', accessToken, 1) // cookie persists 1 day; JWT expiry enforced separately
    setCookieValue('refresh_token', refreshToken, 7)
}

export function clearTokens() {
    deleteCookieValue('access_token')
    deleteCookieValue('refresh_token')
}

let refreshPromise: Promise<string> | null = null

async function refreshTokens(): Promise<string> {
    const refreshToken = getRefreshToken()
    if (!refreshToken) {
        clearTokens()
        window.location.href = '/login'
        throw new Error('No refresh token')
    }

    const data = await ky
        .post(`${BASE_URL}/api/auth/refresh`, {
            json: { refresh_token: refreshToken },
        })
        .json<{ access_token: string; refresh_token: string }>()

    saveTokens(data.access_token, data.refresh_token)
    return data.access_token
}

export const apiClient: KyInstance = ky.create({
    prefixUrl: BASE_URL,
    timeout: 15000,
    retry: 0,
    hooks: {
        beforeRequest: [
            (request) => {
                const token = getAccessToken()
                console.log('[apiClient] beforeRequest', request.url, 'token:', token ? 'present' : 'missing')
                if (token) {
                    request.headers.set('Authorization', `Bearer ${token}`)
                }
            },
        ],
        afterResponse: [
            async (request, _options, response) => {
                if (response.status !== 401) return response

                try {
                    if (!refreshPromise) {
                        refreshPromise = refreshTokens().finally(() => {
                            refreshPromise = null
                        })
                    }
                    const newToken = await refreshPromise

                    request.headers.set('Authorization', `Bearer ${newToken}`)
                    return ky(request)
                } catch {
                    clearTokens()
                    window.location.href = '/login'
                    return response
                }
            },
        ],
    },
})
