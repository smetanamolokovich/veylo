"use client"

import { useFieldArray, useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { useInviteUser } from "../hooks/use-invite-user"
import { inviteTeamSchema, type InviteTeamValues } from "../schemas"

interface Props {
  onDone: () => void
}

const ROLES = [
  { value: "INSPECTOR", label: "Inspector" },
  { value: "EVALUATOR", label: "Evaluator" },
  { value: "MANAGER", label: "Manager" },
]

export function InviteTeamForm({ onDone }: Props) {
  const { mutateAsync, isPending } = useInviteUser()

  const {
    register,
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<InviteTeamValues>({
    resolver: zodResolver(inviteTeamSchema),
    defaultValues: { invites: [{ email: "", role: "INSPECTOR" }] },
  })

  const { fields, append, remove } = useFieldArray({ control, name: "invites" })

  async function onSubmit(data: InviteTeamValues) {
    await Promise.allSettled(
      data.invites.map((inv) => mutateAsync({ email: inv.email, role: inv.role }))
    )
    onDone()
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-5">
      <div className="flex flex-col gap-3">
        {fields.map((field, index) => (
          <div key={field.id} className="flex gap-2 items-start">
            <div className="flex flex-col gap-1 flex-1">
              {index === 0 && <Label>Email</Label>}
              <Input
                type="email"
                placeholder="colleague@company.com"
                {...register(`invites.${index}.email`)}
              />
              {errors.invites?.[index]?.email && (
                <p className="text-destructive text-xs">{errors.invites[index].email?.message}</p>
              )}
            </div>

            <div className="flex flex-col gap-1 w-36">
              {index === 0 && <Label>Role</Label>}
              <select
                {...register(`invites.${index}.role`)}
                className="h-8 rounded-lg border border-input bg-transparent px-2.5 text-sm outline-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50"
              >
                {ROLES.map((r) => (
                  <option key={r.value} value={r.value}>{r.label}</option>
                ))}
              </select>
            </div>

            {fields.length > 1 && (
              <button
                type="button"
                onClick={() => remove(index)}
                className={`text-muted-foreground hover:text-foreground transition-colors ${index === 0 ? "mt-6" : ""}`}
              >
                <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
                  <path d="M4 4l8 8M12 4l-8 8" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round"/>
                </svg>
              </button>
            )}
          </div>
        ))}
      </div>

      <button
        type="button"
        onClick={() => append({ email: "", role: "INSPECTOR" })}
        className="text-sm text-muted-foreground hover:text-foreground transition-colors text-left"
      >
        + Add another
      </button>

      <div className="flex flex-col gap-2">
        <Button type="submit" disabled={isPending} className="w-full">
          {isPending ? "Sending invitations..." : "Send invitations"}
        </Button>
        <button
          type="button"
          onClick={onDone}
          className="text-sm text-muted-foreground hover:text-foreground transition-colors text-center"
        >
          Skip for now
        </button>
      </div>
    </form>
  )
}
