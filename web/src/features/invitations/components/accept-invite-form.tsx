"use client"

import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { useAcceptInvitation } from "../hooks/use-accept-invitation"
import { acceptInvitationSchema, type AcceptInvitationValues } from "../schemas"

interface Props {
  token: string
  email: string
  organizationName: string
  role: string
}

export function AcceptInviteForm({ token, email, organizationName, role }: Props) {
  const { mutate, isPending, error } = useAcceptInvitation(token)

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<AcceptInvitationValues>({
    resolver: zodResolver(acceptInvitationSchema),
  })

  return (
    <form onSubmit={handleSubmit((data) => mutate(data))} className="flex flex-col gap-5">
      <div className="flex flex-col gap-1 rounded-lg bg-muted px-3 py-2.5 text-sm">
        <span className="text-muted-foreground text-xs">Joining as</span>
        <span className="font-medium">{email}</span>
        <span className="text-muted-foreground text-xs capitalize">
          {role.toLowerCase()} at {organizationName}
        </span>
      </div>

      <div className="flex flex-col gap-2">
        <Label htmlFor="full_name">Full name</Label>
        <Input
          id="full_name"
          type="text"
          placeholder="John Smith"
          autoComplete="name"
          {...register("full_name")}
        />
        {errors.full_name && (
          <p className="text-destructive text-xs">{errors.full_name.message}</p>
        )}
      </div>

      <div className="flex flex-col gap-2">
        <Label htmlFor="password">Password</Label>
        <Input
          id="password"
          type="password"
          placeholder="Min. 8 characters"
          autoComplete="new-password"
          {...register("password")}
        />
        {errors.password && (
          <p className="text-destructive text-xs">{errors.password.message}</p>
        )}
      </div>

      {error && (
        <p className="text-destructive text-sm">{error.message}</p>
      )}

      <Button type="submit" disabled={isPending} className="w-full">
        {isPending ? "Joining..." : `Join ${organizationName}`}
      </Button>
    </form>
  )
}
