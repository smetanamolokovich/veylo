import ky from "ky"
import { apiClient } from "@/lib/api-client"
import type {
  InviteUserRequest,
  InviteUserResponse,
  InvitationInfoResponse,
  AcceptInvitationRequest,
  AcceptInvitationResponse,
} from "./types"

const BASE_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080"

// Public client — no auth header
const publicClient = ky.create({ prefixUrl: BASE_URL })

export async function inviteUser(data: InviteUserRequest): Promise<InviteUserResponse> {
  return apiClient
    .post("api/v1/organizations/me/invitations", { json: data })
    .json<InviteUserResponse>()
}

export async function getInvitation(token: string): Promise<InvitationInfoResponse> {
  return publicClient
    .get(`api/auth/invite/${token}`)
    .json<InvitationInfoResponse>()
}

export async function acceptInvitation(
  token: string,
  data: AcceptInvitationRequest
): Promise<AcceptInvitationResponse> {
  return publicClient
    .post(`api/auth/invite/${token}/accept`, { json: data })
    .json<AcceptInvitationResponse>()
}
