"use client"

import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Button } from "@/components/ui/button"
import { createOrganizationSchema, type CreateOrganizationFormValues } from "../schemas"
import { useCreateOrganization } from "../hooks/use-create-organization"

export function CreateOrganizationForm() {
  const { mutate, isPending, error } = useCreateOrganization()

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateOrganizationFormValues>({
    resolver: zodResolver(createOrganizationSchema),
    defaultValues: { vertical: "vehicle" },
  })

  return (
    <form onSubmit={handleSubmit((data) => mutate(data))} className="flex flex-col gap-5">
      <input type="hidden" {...register("vertical")} />

      <div className="flex flex-col gap-2">
        <Label htmlFor="org_name">Workspace name</Label>
        <Input
          id="org_name"
          type="text"
          placeholder="Acme Fleet Management"
          autoComplete="organization"
          {...register("org_name")}
        />
        {errors.org_name && (
          <p className="text-destructive text-xs">{errors.org_name.message}</p>
        )}
        <p className="text-muted-foreground text-xs">
          This is the name your team will see. You can change it later.
        </p>
      </div>

      {error && (
        <p className="text-destructive text-sm">{error.message}</p>
      )}

      <Button type="submit" disabled={isPending} className="w-full">
        {isPending ? "Creating workspace..." : "Create workspace"}
      </Button>
    </form>
  )
}
