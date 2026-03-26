"use client"

import { useMutation } from "@tanstack/react-query"
import { inviteUser } from "../api"
import type { InviteUserRequest } from "../types"

export function useInviteUser() {
  return useMutation({
    mutationFn: (data: InviteUserRequest) => inviteUser(data),
  })
}
