import { apiClient, publicClient } from '@/lib/api-client'
import type {
    InviteUserRequest,
    InviteUserResponse,
    InvitationInfoResponse,
    AcceptInvitationRequest,
    AcceptInvitationResponse,
} from './types'

export async function inviteUser(data: InviteUserRequest): Promise<InviteUserResponse> {
    return apiClient
        .post('api/v1/organizations/me/invitations', { json: data })
        .json<InviteUserResponse>()
}

export async function getInvitation(token: string): Promise<InvitationInfoResponse> {
    return publicClient.get(`api/auth/invite/${token}`).json<InvitationInfoResponse>()
}

export async function acceptInvitation(
    token: string,
    data: AcceptInvitationRequest,
): Promise<AcceptInvitationResponse> {
    return publicClient
        .post(`api/auth/invite/${token}/accept`, { json: data })
        .json<AcceptInvitationResponse>()
}
