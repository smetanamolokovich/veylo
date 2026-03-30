'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { CreateOrganizationForm } from '@/features/auth/components/create-organization-form'
import { InviteTeamForm } from '@/features/invitations/components/invite-team-form'
import { completeOnboarding } from '@/features/auth/api'

type Step = 2 | 3

function StepIndicator({ current }: { current: Step }) {
    return (
        <div className="flex items-center gap-2 justify-center">
            <div className="flex items-center gap-1.5">
                <div className="w-6 h-6 rounded-full bg-muted text-muted-foreground text-xs font-semibold flex items-center justify-center">
                    <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
                        <path
                            d="M2 6l3 3 5-5"
                            stroke="currentColor"
                            strokeWidth="1.5"
                            strokeLinecap="round"
                            strokeLinejoin="round"
                        />
                    </svg>
                </div>
                <span className="text-xs text-muted-foreground">Your account</span>
            </div>
            <div className="h-px w-8 bg-foreground/30" />
            <div className="flex items-center gap-1.5">
                <div
                    className={`w-6 h-6 rounded-full text-xs font-semibold flex items-center justify-center ${current > 2 ? 'bg-muted text-muted-foreground' : 'bg-foreground text-background'}`}
                >
                    {current > 2 ? (
                        <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
                            <path
                                d="M2 6l3 3 5-5"
                                stroke="currentColor"
                                strokeWidth="1.5"
                                strokeLinecap="round"
                                strokeLinejoin="round"
                            />
                        </svg>
                    ) : (
                        '2'
                    )}
                </div>
                <span
                    className={`text-xs ${current === 2 ? 'font-medium' : 'text-muted-foreground'}`}
                >
                    Workspace
                </span>
            </div>
            <div className="h-px w-8 bg-foreground/30" />
            <div className="flex items-center gap-1.5">
                <div
                    className={`w-6 h-6 rounded-full text-xs font-semibold flex items-center justify-center ${current === 3 ? 'bg-foreground text-background' : 'bg-muted text-muted-foreground'}`}
                >
                    3
                </div>
                <span
                    className={`text-xs ${current === 3 ? 'font-medium' : 'text-muted-foreground'}`}
                >
                    Invite team
                </span>
            </div>
        </div>
    )
}

export default function OnboardingPage() {
    const [step, setStep] = useState<Step>(2)
    const [finishError, setFinishError] = useState<string | null>(null)
    const router = useRouter()

    async function finishOnboarding() {
        try {
            setFinishError(null)
            await completeOnboarding()
            router.push('/dashboard')
        } catch {
            setFinishError('Something went wrong. Please try again.')
        }
    }

    return (
        <div className="flex flex-col gap-6">
            <StepIndicator current={step} />

            {step === 2 && (
                <Card>
                    <CardHeader>
                        <CardTitle>Set up your workspace</CardTitle>
                        <CardDescription>
                            Give your organization a name to get started.
                        </CardDescription>
                    </CardHeader>
                    <CardContent>
                        <CreateOrganizationForm onComplete={() => setStep(3)} />
                    </CardContent>
                </Card>
            )}

            {step === 3 && (
                <Card>
                    <CardHeader>
                        <CardTitle>Invite your team</CardTitle>
                        <CardDescription>
                            Add inspectors, evaluators, and managers. You can also do this later.
                        </CardDescription>
                    </CardHeader>
                    <CardContent>
                        <InviteTeamForm onDone={finishOnboarding} />
                        {finishError && (
                            <p className="mt-3 text-sm text-destructive">{finishError}</p>
                        )}
                    </CardContent>
                </Card>
            )}
        </div>
    )
}
