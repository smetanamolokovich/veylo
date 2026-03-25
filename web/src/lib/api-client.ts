import ky, { type KyInstance } from "ky"

const BASE_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080"

function getAccessToken(): string | null {
  if (typeof window === "undefined") return null
  return localStorage.getItem("access_token")
}

function getRefreshToken(): string | null {
  if (typeof window === "undefined") return null
  return localStorage.getItem("refresh_token")
}

function saveTokens(accessToken: string, refreshToken: string) {
  localStorage.setItem("access_token", accessToken)
  localStorage.setItem("refresh_token", refreshToken)
}

function clearTokens() {
  localStorage.removeItem("access_token")
  localStorage.removeItem("refresh_token")
  localStorage.removeItem("organization_id")
}

let isRefreshing = false

export const apiClient: KyInstance = ky.create({
  prefixUrl: BASE_URL,
  hooks: {
    beforeRequest: [
      (request) => {
        const token = getAccessToken()
        if (token) {
          request.headers.set("Authorization", `Bearer ${token}`)
        }
      },
    ],
    afterResponse: [
      async (request, options, response) => {
        if (response.status !== 401) return response
        if (isRefreshing) return response

        const refreshToken = getRefreshToken()
        if (!refreshToken) {
          clearTokens()
          window.location.href = "/login"
          return response
        }

        isRefreshing = true
        try {
          const data = await ky
            .post(`${BASE_URL}/api/auth/refresh`, {
              json: { refresh_token: refreshToken },
            })
            .json<{ access_token: string; refresh_token: string }>()

          saveTokens(data.access_token, data.refresh_token)

          request.headers.set("Authorization", `Bearer ${data.access_token}`)
          return ky(request)
        } catch {
          clearTokens()
          window.location.href = "/login"
          return response
        } finally {
          isRefreshing = false
        }
      },
    ],
  },
})

export { saveTokens, clearTokens, getAccessToken }
