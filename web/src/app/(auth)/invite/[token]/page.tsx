'use client'

import { use } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { useInvitation } from '@/features/invitations/hooks/use-invitation'
import { AcceptInviteForm } from '@/features/invitations/components/accept-invite-form'

interface Props {
    params: Promise<{ token: string }>
}

export default function InvitePage({ params }: Props) {
    const { token } = use(params)
    const { data, isLoading, isError } = useInvitation(token)

    if (isLoading) {
        return (
            <div className="flex flex-col gap-6">
                <Card>
                    <CardContent className="pt-6 pb-6 text-center text-sm text-muted-foreground">
                        Loading invitation...
                    </CardContent>
                </Card>
            </div>
        )
    }

    if (isError || !data) {
        return (
            <div className="flex flex-col gap-6">
                <Card>
                    <CardContent className="pt-6 pb-6 text-center">
                        <p className="text-sm font-medium">Invalid invitation</p>
                        <p className="text-xs text-muted-foreground mt-1">
                            This invitation link is invalid or has already been used.
                        </p>
                    </CardContent>
                </Card>
            </div>
        )
    }

    if (data.is_expired) {
        return (
            <div className="flex flex-col gap-6">
                <Card>
                    <CardContent className="pt-6 pb-6 text-center">
                        <p className="text-sm font-medium">Invitation expired</p>
                        <p className="text-xs text-muted-foreground mt-1">
                            Ask your team admin to send a new invitation.
                        </p>
                    </CardContent>
                </Card>
            </div>
        )
    }

    return (
        <div className="flex flex-col gap-6">
            <Card>
                <CardHeader>
                    <CardTitle>Join {data.organization_name}</CardTitle>
                    <CardDescription>
                        Create your account to accept this invitation.
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <AcceptInviteForm
                        token={token}
                        email={data.email}
                        organizationName={data.organization_name}
                        role={data.role}
                    />
                </CardContent>
            </Card>
        </div>
    )
}
