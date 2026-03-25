import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { CreateOrganizationForm } from "@/features/auth/components/create-organization-form"

export default function OnboardingPage() {
  return (
    <div className="flex flex-col gap-6">
      {/* Step indicator */}
      <div className="flex items-center gap-2 justify-center">
        <div className="flex items-center gap-1.5">
          <div className="w-6 h-6 rounded-full bg-muted text-muted-foreground text-xs font-semibold flex items-center justify-center">
            <svg width="12" height="12" viewBox="0 0 12 12" fill="none"><path d="M2 6l3 3 5-5" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/></svg>
          </div>
          <span className="text-xs text-muted-foreground">Your account</span>
        </div>
        <div className="h-px w-8 bg-foreground/30" />
        <div className="flex items-center gap-1.5">
          <div className="w-6 h-6 rounded-full bg-foreground text-background text-xs font-semibold flex items-center justify-center">2</div>
          <span className="text-xs font-medium">Workspace</span>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Set up your workspace</CardTitle>
          <CardDescription>Give your organization a name to get started.</CardDescription>
        </CardHeader>
        <CardContent>
          <CreateOrganizationForm />
        </CardContent>
      </Card>
    </div>
  )
}
