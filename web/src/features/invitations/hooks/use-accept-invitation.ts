"use client"

import { useMutation } from "@tanstack/react-query"
import { useRouter } from "next/navigation"
import { saveTokens } from "@/lib/api-client"
import { acceptInvitation } from "../api"
import type { AcceptInvitationRequest } from "../types"

export function useAcceptInvitation(token: string) {
  const router = useRouter()

  return useMutation({
    mutationFn: (data: AcceptInvitationRequest) => acceptInvitation(token, data),
    onSuccess(data) {
      saveTokens(data.access_token, data.refresh_token)
      localStorage.setItem("organization_id", "")
      router.push("/dashboard")
    },
  })
}
