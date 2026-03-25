"use client"

import { useMutation } from "@tanstack/react-query"
import { useRouter } from "next/navigation"
import { saveTokens } from "@/lib/api-client"
import { login } from "../api"
import type { LoginRequest } from "../types"

export function useLogin() {
  const router = useRouter()

  return useMutation({
    mutationFn: (data: LoginRequest) => login(data),
    onSuccess(data) {
      saveTokens(data.access_token, data.refresh_token)
      localStorage.setItem("user_id", data.user_id)
      router.push("/dashboard")
    },
  })
}
