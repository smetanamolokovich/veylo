import Link from 'next/link'
import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from '@/components/ui/card'
import { LoginForm } from '@/features/auth/components/login-form'

export default function LoginPage() {
    return (
        <Card>
            <CardHeader>
                <CardTitle>Welcome back</CardTitle>
                <CardDescription>Sign in to your workspace</CardDescription>
            </CardHeader>
            <CardContent>
                <LoginForm />
            </CardContent>
            <CardFooter className="justify-center">
                <p className="text-sm text-muted-foreground">
                    Don&apos;t have an account?{' '}
                    <Link
                        href="/signup"
                        className="text-foreground font-medium underline underline-offset-4"
                    >
                        Create workspace
                    </Link>
                </p>
            </CardFooter>
        </Card>
    )
}
