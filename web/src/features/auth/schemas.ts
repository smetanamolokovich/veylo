import { z } from "zod"

export const registerSchema = z.object({
  full_name: z.string().min(2, "Full name must be at least 2 characters"),
  email: z.string().email("Invalid email address"),
  password: z.string().min(8, "Password must be at least 8 characters"),
})

export const createOrganizationSchema = z.object({
  org_name: z.string().min(2, "Workspace name must be at least 2 characters"),
  vertical: z.enum(["vehicle"]),
})

export const loginSchema = z.object({
  email: z.string().email("Invalid email address"),
  password: z.string().min(1, "Password is required"),
})

export type RegisterFormValues = z.infer<typeof registerSchema>
export type CreateOrganizationFormValues = z.infer<typeof createOrganizationSchema>
export type LoginFormValues = z.infer<typeof loginSchema>
