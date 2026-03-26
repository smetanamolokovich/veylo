export interface InviteUserRequest {
  email: string
  role: string
}

export interface InviteUserResponse {
  invitation_id: string
  invite_token: string
  email: string
  role: string
  expires_at: string
}

export interface InvitationInfoResponse {
  email: string
  organization_name: string
  role: string
  expires_at: string
  is_expired: boolean
}

export interface AcceptInvitationRequest {
  full_name: string
  password: string
}

export interface AcceptInvitationResponse {
  user_id: string
  access_token: string
  refresh_token: string
}
