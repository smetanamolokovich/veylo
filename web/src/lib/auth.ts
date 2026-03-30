interface JWTPayload {
    sub: string
    org_id: string
    role: string
    exp: number
}

export function decodeJWT(token: string): JWTPayload | null {
    try {
        const parts = token.split('.')
        if (parts.length !== 3) return null
        const payload = parts[1]
        const decoded = atob(payload.replace(/-/g, '+').replace(/_/g, '/'))
        return JSON.parse(decoded) as JWTPayload
    } catch {
        return null
    }
}

export function isAuthenticated(): boolean {
    if (typeof document === 'undefined') return false
    const token = getCookieValue('access_token')
    if (!token) return false
    const payload = decodeJWT(token)
    if (!payload) return false
    return payload.exp * 1000 > Date.now()
}

export function getOrgID(): string {
    if (typeof document === 'undefined') return ''
    const token = getCookieValue('access_token')
    if (!token) return ''
    const payload = decodeJWT(token)
    return payload?.org_id ?? ''
}

export function getUserID(): string {
    if (typeof document === 'undefined') return ''
    const token = getCookieValue('access_token')
    if (!token) return ''
    const payload = decodeJWT(token)
    return payload?.sub ?? ''
}

// Cookie helpers — used internally and by api-client
export function getCookieValue(name: string): string | null {
    if (typeof document === 'undefined') return null
    const match = document.cookie.match(new RegExp('(?:^|; )' + name + '=([^;]*)'))
    return match ? decodeURIComponent(match[1]) : null
}

export function setCookieValue(name: string, value: string, maxAgeDays: number) {
    document.cookie = `${name}=${encodeURIComponent(value)}; max-age=${maxAgeDays * 86400}; path=/; SameSite=Lax`
}

export function deleteCookieValue(name: string) {
    document.cookie = `${name}=; max-age=0; path=/`
}
