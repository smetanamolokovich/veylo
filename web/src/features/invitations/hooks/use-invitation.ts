"use client"

import { useQuery } from "@tanstack/react-query"
import { getInvitation } from "../api"

export function useInvitation(token: string) {
  return useQuery({
    queryKey: ["invitation", token],
    queryFn: () => getInvitation(token),
    enabled: !!token,
    retry: false,
  })
}
