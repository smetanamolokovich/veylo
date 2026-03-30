export interface RegisterRequest {
    email: string
    password: string
    full_name: string
}

export interface RegisterResponse {
    user_id: string
    access_token: string
    refresh_token: string
}

export interface CreateOrganizationRequest {
    org_name: string
    vertical: string
}

export interface CreateOrganizationResponse {
    organization_id: string
    access_token: string
}

export interface LoginRequest {
    email: string
    password: string
}

export interface LoginResponse {
    access_token: string
    refresh_token: string
    user_id: string
    role: string
}
