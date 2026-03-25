"use client"

import { useMutation } from "@tanstack/react-query"
import { useRouter } from "next/navigation"
import { saveTokens } from "@/lib/api-client"
import { register } from "../api"
import type { RegisterRequest } from "../types"

export function useRegister() {
  const router = useRouter()

  return useMutation({
    mutationFn: (data: RegisterRequest) => register(data),
    onSuccess(data) {
      saveTokens(data.access_token, data.refresh_token)
      localStorage.setItem("user_id", data.user_id)
      router.push("/onboarding")
    },
  })
}
