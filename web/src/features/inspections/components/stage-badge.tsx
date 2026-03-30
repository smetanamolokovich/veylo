import { Badge } from '@/components/ui/badge'
import { cn } from '@/lib/utils'
import type { SystemStage } from '../types'

const stageColors: Record<SystemStage, string> = {
    ENTRY: 'bg-blue-100 text-blue-700 border-transparent',
    EVALUATION: 'bg-yellow-100 text-yellow-700 border-transparent',
    REVIEW: 'bg-orange-100 text-orange-700 border-transparent',
    FINAL: 'bg-green-100 text-green-700 border-transparent',
}

const stageLabels: Record<SystemStage, string> = {
    ENTRY: 'Entry',
    EVALUATION: 'Evaluation',
    REVIEW: 'Review',
    FINAL: 'Final',
}

export function StageBadge({ stage }: { stage: SystemStage }) {
    return (
        <Badge variant="outline" className={cn(stageColors[stage])}>
            {stageLabels[stage]}
        </Badge>
    )
}
