import { cn } from '@/lib/utils'

interface VeyloLogoProps {
    /** Size of the icon container in pixels. Default: 32 (w-8 h-8). */
    size?: 'sm' | 'default'
    className?: string
}

/**
 * Veylo logotype — icon + wordmark.
 *
 * "default" variant: dark background container (used on auth panel, mobile header)
 * "sm" variant: slightly smaller container (used in sidebar)
 */
export function VeyloLogo({ size = 'default', className }: VeyloLogoProps) {
    const containerSize = size === 'sm' ? 'w-7 h-7' : 'w-8 h-8'
    const rounded = size === 'sm' ? 'rounded-md' : 'rounded-lg'
    const iconSize = size === 'sm' ? 15 : 18
    const textSize = size === 'sm' ? 'text-sm' : 'text-lg'

    return (
        <div className={cn('flex items-center gap-2', className)}>
            <div
                className={cn(
                    'bg-foreground flex items-center justify-center shrink-0',
                    containerSize,
                    rounded,
                )}
            >
                <svg
                    width={iconSize}
                    height={iconSize}
                    viewBox="0 0 18 18"
                    fill="none"
                    xmlns="http://www.w3.org/2000/svg"
                >
                    <rect x="1" y="5" width="16" height="10" rx="2" fill="white" />
                    <rect x="4" y="2" width="10" height="4" rx="1.5" fill="white" />
                    <circle cx="5" cy="14" r="1.5" fill="oklch(0.145 0 0)" />
                    <circle cx="13" cy="14" r="1.5" fill="oklch(0.145 0 0)" />
                </svg>
            </div>
            <span className={cn('font-semibold tracking-tight', textSize)}>Veylo</span>
        </div>
    )
}

/**
 * Inverted variant for use on dark backgrounds (auth panel left side).
 * Icon has a light container; text is inherited (typically white).
 */
export function VeyloLogoDark({ className }: { className?: string }) {
    return (
        <div className={cn('flex items-center gap-2.5', className)}>
            <div className="w-8 h-8 bg-background rounded-lg flex items-center justify-center">
                <svg
                    width="18"
                    height="18"
                    viewBox="0 0 18 18"
                    fill="none"
                    xmlns="http://www.w3.org/2000/svg"
                >
                    <rect x="1" y="5" width="16" height="10" rx="2" className="fill-foreground" />
                    <rect x="4" y="2" width="10" height="4" rx="1.5" className="fill-foreground" />
                    <circle cx="5" cy="14" r="1.5" className="fill-background" />
                    <circle cx="13" cy="14" r="1.5" className="fill-background" />
                </svg>
            </div>
            <span className="font-semibold text-lg tracking-tight">Veylo</span>
        </div>
    )
}
