import { Badge } from '@/components/ui/badge'
import { cn } from '@/lib/utils'

export function StatusBadge({ status }: { status: string }) {
    return (
        <Badge
            variant="outline"
            className={cn('bg-muted text-muted-foreground border-transparent capitalize')}
        >
            {status.replace(/_/g, ' ')}
        </Badge>
    )
}
