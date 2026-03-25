"use client"

import { useMutation } from "@tanstack/react-query"
import { useRouter } from "next/navigation"
import { saveTokens } from "@/lib/api-client"
import { createOrganization, completeOnboarding } from "../api"
import type { CreateOrganizationRequest } from "../types"

export function useCreateOrganization() {
  const router = useRouter()

  return useMutation({
    mutationFn: async (data: CreateOrganizationRequest) => {
      const org = await createOrganization(data)
      // Swap in the new access token that carries org_id before calling onboarding
      saveTokens(org.access_token, localStorage.getItem("refresh_token") ?? "")
      localStorage.setItem("organization_id", org.organization_id)
      await completeOnboarding()
      return org
    },
    onSuccess() {
      router.push("/dashboard")
    },
  })
}
