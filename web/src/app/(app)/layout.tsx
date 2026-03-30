'use client'

import { usePathname, useRouter } from 'next/navigation'
import Link from 'next/link'
import { LayoutDashboard, ClipboardList, Car, GitBranch, LogOut } from 'lucide-react'
import { clearTokens } from '@/lib/api-client'
import { cn } from '@/lib/utils'
import { VeyloLogo } from '@/components/logo'

const navItems = [
    { href: '/dashboard', label: 'Dashboard', icon: LayoutDashboard },
    { href: '/inspections', label: 'Inspections', icon: ClipboardList },
    { href: '/assets/vehicles', label: 'Vehicles', icon: Car },
    { href: '/workflow', label: 'Workflow', icon: GitBranch },
]

export default function AppLayout({ children }: { children: React.ReactNode }) {
    const router = useRouter()
    const pathname = usePathname()

    return (
        <div className="min-h-screen flex">
            {/* Sidebar */}
            <aside className="w-56 border-r bg-sidebar flex flex-col shrink-0">
                {/* Logo */}
                <div className="h-14 flex items-center px-4 border-b">
                    <VeyloLogo size="sm" />
                </div>

                {/* Navigation */}
                <nav className="flex-1 px-2 py-3 flex flex-col gap-0.5">
                    {navItems.map((item) => {
                        const isActive =
                            pathname === item.href || pathname.startsWith(item.href + '/')
                        return (
                            <Link
                                key={item.href}
                                href={item.href}
                                className={cn(
                                    'flex items-center gap-2.5 px-2.5 py-2 rounded-md text-sm transition-colors',
                                    isActive
                                        ? 'bg-sidebar-primary text-sidebar-primary-foreground font-medium'
                                        : 'text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground',
                                )}
                            >
                                <item.icon size={16} className="shrink-0" />
                                {item.label}
                            </Link>
                        )
                    })}
                </nav>

                {/* Bottom */}
                <div className="p-3 border-t">
                    <button
                        onClick={() => {
                            clearTokens()
                            router.push('/login')
                        }}
                        className="w-full flex items-center gap-2.5 px-2.5 py-2 rounded-md text-sm text-muted-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground transition-colors"
                    >
                        <LogOut size={16} />
                        Sign out
                    </button>
                </div>
            </aside>

            {/* Main content */}
            <main className="flex-1 overflow-auto">{children}</main>
        </div>
    )
}
