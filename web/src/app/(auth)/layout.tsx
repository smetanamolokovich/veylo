export default function AuthLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen grid lg:grid-cols-2">
      {/* Left brand panel — hidden on mobile */}
      <div className="hidden lg:flex flex-col bg-foreground text-background p-10 relative overflow-hidden">
        {/* Background texture */}
        <div className="absolute inset-0 bg-[linear-gradient(135deg,oklch(0.145_0_0)_0%,oklch(0.25_0_0)_100%)]" />
        <div className="absolute inset-0 opacity-5" style={{ backgroundImage: 'radial-gradient(circle at 2px 2px, white 1px, transparent 0)', backgroundSize: '32px 32px' }} />

        {/* Logo */}
        <div className="relative z-10 flex items-center gap-2.5">
          <div className="w-8 h-8 bg-background rounded-lg flex items-center justify-center">
            <svg width="18" height="18" viewBox="0 0 18 18" fill="none" xmlns="http://www.w3.org/2000/svg">
              <rect x="1" y="5" width="16" height="10" rx="2" className="fill-foreground" />
              <rect x="4" y="2" width="10" height="4" rx="1.5" className="fill-foreground" />
              <circle cx="5" cy="14" r="1.5" className="fill-background" />
              <circle cx="13" cy="14" r="1.5" className="fill-background" />
            </svg>
          </div>
          <span className="font-semibold text-lg tracking-tight">Veylo</span>
        </div>

        {/* Tagline */}
        <div className="relative z-10 mt-auto">
          <blockquote className="space-y-2">
            <p className="text-xl font-medium leading-snug">
              &ldquo;The inspection platform built for teams that can&apos;t afford mistakes.&rdquo;
            </p>
            <footer className="text-sm opacity-60">Inspection Management · Fleet · Leasing · Insurance</footer>
          </blockquote>
        </div>
      </div>

      {/* Right form panel */}
      <div className="flex flex-col items-center justify-center p-6 sm:p-10">
        {/* Mobile logo */}
        <div className="lg:hidden mb-8 flex items-center gap-2">
          <div className="w-8 h-8 bg-foreground rounded-lg flex items-center justify-center">
            <svg width="18" height="18" viewBox="0 0 18 18" fill="none" xmlns="http://www.w3.org/2000/svg">
              <rect x="1" y="5" width="16" height="10" rx="2" fill="white" />
              <rect x="4" y="2" width="10" height="4" rx="1.5" fill="white" />
              <circle cx="5" cy="14" r="1.5" fill="oklch(0.145 0 0)" />
              <circle cx="13" cy="14" r="1.5" fill="oklch(0.145 0 0)" />
            </svg>
          </div>
          <span className="font-semibold text-lg tracking-tight">Veylo</span>
        </div>

        <div className="w-full max-w-sm">
          {children}
        </div>
      </div>
    </div>
  )
}
