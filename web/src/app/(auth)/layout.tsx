import { VeyloLogo, VeyloLogoDark } from '@/components/logo'

export default function AuthLayout({ children }: { children: React.ReactNode }) {
    return (
        <div className="min-h-screen grid lg:grid-cols-2">
            {/* Left brand panel — hidden on mobile */}
            <div className="hidden lg:flex flex-col bg-foreground text-background p-10 relative overflow-hidden">
                {/* Background texture */}
                <div className="absolute inset-0 bg-[linear-gradient(135deg,oklch(0.145_0_0)_0%,oklch(0.25_0_0)_100%)]" />
                <div
                    className="absolute inset-0 opacity-5"
                    style={{
                        backgroundImage:
                            'radial-gradient(circle at 2px 2px, white 1px, transparent 0)',
                        backgroundSize: '32px 32px',
                    }}
                />

                {/* Logo */}
                <div className="relative z-10">
                    <VeyloLogoDark />
                </div>

                {/* Tagline */}
                <div className="relative z-10 mt-auto">
                    <blockquote className="space-y-2">
                        <p className="text-xl font-medium leading-snug">
                            &ldquo;The inspection platform built for teams that can&apos;t afford
                            mistakes.&rdquo;
                        </p>
                        <footer className="text-sm opacity-60">
                            Inspection Management · Fleet · Leasing · Insurance
                        </footer>
                    </blockquote>
                </div>
            </div>

            {/* Right form panel */}
            <div className="flex flex-col items-center justify-center p-6 sm:p-10">
                {/* Mobile logo */}
                <div className="lg:hidden mb-8">
                    <VeyloLogo />
                </div>

                <div className="w-full max-w-sm">{children}</div>
            </div>
        </div>
    )
}
