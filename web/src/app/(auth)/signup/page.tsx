import Link from "next/link"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { RegisterForm } from "@/features/auth/components/register-form"

export default function SignupPage() {
  return (
    <div className="flex flex-col gap-6">
      {/* Step indicator */}
      <div className="flex items-center gap-2 justify-center">
        <div className="flex items-center gap-1.5">
          <div className="w-6 h-6 rounded-full bg-foreground text-background text-xs font-semibold flex items-center justify-center">1</div>
          <span className="text-xs font-medium">Your account</span>
        </div>
        <div className="h-px w-8 bg-border" />
        <div className="flex items-center gap-1.5">
          <div className="w-6 h-6 rounded-full bg-muted text-muted-foreground text-xs font-semibold flex items-center justify-center">2</div>
          <span className="text-xs text-muted-foreground">Workspace</span>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Create your account</CardTitle>
          <CardDescription>Start managing inspections today. Free to get started.</CardDescription>
        </CardHeader>
        <CardContent>
          <RegisterForm />
        </CardContent>
        <CardFooter className="justify-center">
          <p className="text-sm text-muted-foreground">
            Already have an account?{" "}
            <Link href="/login" className="text-foreground font-medium underline underline-offset-4">
              Sign in
            </Link>
          </p>
        </CardFooter>
      </Card>
    </div>
  )
}
