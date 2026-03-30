import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

const PUBLIC_PATHS = ['/login', '/signup', '/onboarding', '/invite']
const AUTH_REDIRECT_PATHS = ['/login', '/signup']

export function proxy(request: NextRequest) {
    const { pathname } = request.nextUrl
    const token = request.cookies.get('access_token')?.value

    const isPublic = PUBLIC_PATHS.some((p) => pathname.startsWith(p))

    // Not authenticated → redirect to login
    if (!isPublic && !token) {
        return NextResponse.redirect(new URL('/login', request.url))
    }

    // Authenticated → redirect away from login/signup to dashboard
    if (AUTH_REDIRECT_PATHS.some((p) => pathname.startsWith(p)) && token) {
        return NextResponse.redirect(new URL('/dashboard', request.url))
    }

    return NextResponse.next()
}

export const config = {
    matcher: ['/((?!_next/static|_next/image|favicon.ico|public).*)'],
}
