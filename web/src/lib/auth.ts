interface JWTPayload {
  sub: string
  org_id: string
  role: string
  exp: number
}

export function decodeJWT(token: string): JWTPayload | null {
  try {
    const parts = token.split(".")
    if (parts.length !== 3) return null
    const payload = parts[1]
    const decoded = atob(payload.replace(/-/g, "+").replace(/_/g, "/"))
    return JSON.parse(decoded) as JWTPayload
  } catch {
    return null
  }
}

export function isAuthenticated(): boolean {
  if (typeof window === "undefined") return false
  const token = localStorage.getItem("access_token")
  if (!token) return false
  const payload = decodeJWT(token)
  if (!payload) return false
  return payload.exp * 1000 > Date.now()
}

export function getOrgID(): string {
  if (typeof window === "undefined") return ""
  const token = localStorage.getItem("access_token")
  if (!token) return ""
  const payload = decodeJWT(token)
  return payload?.org_id ?? ""
}

export function getUserID(): string {
  if (typeof window === "undefined") return ""
  const token = localStorage.getItem("access_token")
  if (!token) return ""
  const payload = decodeJWT(token)
  return payload?.sub ?? ""
}
