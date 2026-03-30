import { z } from 'zod'

export const inviteRowSchema = z.object({
    email: z.string().email('Invalid email'),
    role: z.enum(['INSPECTOR', 'EVALUATOR', 'MANAGER'], { message: 'Select a role' }),
})

export const inviteTeamSchema = z.object({
    invites: z.array(inviteRowSchema).min(1),
})

export const acceptInvitationSchema = z.object({
    full_name: z.string().min(2, 'Full name must be at least 2 characters'),
    password: z.string().min(8, 'Password must be at least 8 characters'),
})

export type InviteRowValues = z.infer<typeof inviteRowSchema>
export type InviteTeamValues = z.infer<typeof inviteTeamSchema>
export type AcceptInvitationValues = z.infer<typeof acceptInvitationSchema>
