'use client'

import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { cn } from '@/lib/utils'
import { buttonVariants } from '@/components/ui/button'
import { type InspectionItem, DEFAULT_STATUS_STAGE_MAP } from '../types'
import { StageBadge } from './stage-badge'
import { StatusBadge } from './status-badge'
import { InspectionSkeletonRow } from './inspection-skeleton'

// ── Empty state ───────────────────────────────────────────────────────────────

function EmptyState() {
    return (
        <tr>
            <td colSpan={5} className="px-4 py-16 text-center">
                <div className="flex flex-col items-center gap-4">
                    <div className="w-12 h-12 rounded-full bg-muted flex items-center justify-center">
                        <svg
                            width="22"
                            height="22"
                            viewBox="0 0 24 24"
                            fill="none"
                            xmlns="http://www.w3.org/2000/svg"
                            className="text-muted-foreground"
                        >
                            <path
                                d="M9 17H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11a2 2 0 0 1 2 2v3"
                                stroke="currentColor"
                                strokeWidth="1.5"
                                strokeLinecap="round"
                            />
                            <rect
                                x="13"
                                y="13"
                                width="8"
                                height="8"
                                rx="1"
                                stroke="currentColor"
                                strokeWidth="1.5"
                            />
                            <path
                                d="M16 16v3M16 16h3"
                                stroke="currentColor"
                                strokeWidth="1.5"
                                strokeLinecap="round"
                            />
                        </svg>
                    </div>
                    <div className="space-y-1 text-center">
                        <p className="text-sm font-medium">No inspections yet</p>
                        <p className="text-xs text-muted-foreground">
                            Create your first inspection to get started.
                        </p>
                    </div>
                    <Link
                        href="/inspections/new"
                        className={buttonVariants({ variant: 'default', size: 'sm' })}
                    >
                        New inspection
                    </Link>
                </div>
            </td>
        </tr>
    )
}

// ── Pagination ────────────────────────────────────────────────────────────────

interface PaginationProps {
    page: number
    pageSize: number
    total: number
    onPageChange: (page: number) => void
}

function Pagination({ page, pageSize, total, onPageChange }: PaginationProps) {
    const totalPages = Math.max(1, Math.ceil(total / pageSize))
    const from = total === 0 ? 0 : (page - 1) * pageSize + 1
    const to = Math.min(page * pageSize, total)

    return (
        <div className="flex items-center justify-between px-4 py-3 border-t text-sm text-muted-foreground">
            <span>{total === 0 ? 'No results' : `${from}–${to} of ${total}`}</span>
            <div className="flex items-center gap-1">
                <button
                    onClick={() => onPageChange(page - 1)}
                    disabled={page <= 1}
                    className={cn(
                        'inline-flex items-center justify-center rounded-md h-7 w-7 text-sm transition-colors',
                        page <= 1
                            ? 'opacity-40 pointer-events-none'
                            : 'hover:bg-muted hover:text-foreground',
                    )}
                    aria-label="Previous page"
                >
                    <svg width="14" height="14" viewBox="0 0 16 16" fill="none">
                        <path
                            d="M10 12L6 8l4-4"
                            stroke="currentColor"
                            strokeWidth="1.5"
                            strokeLinecap="round"
                            strokeLinejoin="round"
                        />
                    </svg>
                </button>
                <span className="px-2 tabular-nums">
                    {page} / {totalPages}
                </span>
                <button
                    onClick={() => onPageChange(page + 1)}
                    disabled={page >= totalPages}
                    className={cn(
                        'inline-flex items-center justify-center rounded-md h-7 w-7 text-sm transition-colors',
                        page >= totalPages
                            ? 'opacity-40 pointer-events-none'
                            : 'hover:bg-muted hover:text-foreground',
                    )}
                    aria-label="Next page"
                >
                    <svg width="14" height="14" viewBox="0 0 16 16" fill="none">
                        <path
                            d="M6 4l4 4-4 4"
                            stroke="currentColor"
                            strokeWidth="1.5"
                            strokeLinecap="round"
                            strokeLinejoin="round"
                        />
                    </svg>
                </button>
            </div>
        </div>
    )
}

// ── Main table component ──────────────────────────────────────────────────────

interface InspectionsTableProps {
    items: InspectionItem[]
    total: number
    page: number
    pageSize: number
    isLoading: boolean
    onPageChange: (page: number) => void
}

export function InspectionsTable({
    items,
    total,
    page,
    pageSize,
    isLoading,
    onPageChange,
}: InspectionsTableProps) {
    const router = useRouter()

    return (
        <div className="rounded-xl border bg-card overflow-hidden">
            <table className="w-full text-sm">
                <thead>
                    <tr className="border-b bg-muted/40">
                        <th className="px-4 py-2.5 text-left font-medium text-muted-foreground">
                            Contract
                        </th>
                        <th className="px-4 py-2.5 text-left font-medium text-muted-foreground">
                            Status
                        </th>
                        <th className="px-4 py-2.5 text-left font-medium text-muted-foreground">
                            Stage
                        </th>
                        <th className="px-4 py-2.5 text-left font-medium text-muted-foreground">
                            Created
                        </th>
                        <th className="px-4 py-2.5 text-right font-medium text-muted-foreground">
                            Action
                        </th>
                    </tr>
                </thead>
                <tbody className="divide-y">
                    {isLoading ? (
                        Array.from({ length: 5 }).map((_, i) => <InspectionSkeletonRow key={i} />)
                    ) : items.length === 0 ? (
                        <EmptyState />
                    ) : (
                        items.map((inspection) => {
                            const stage = DEFAULT_STATUS_STAGE_MAP[inspection.status] ?? 'ENTRY'

                            const formattedDate = inspection.created_at
                                ? new Intl.DateTimeFormat('en-US', {
                                      month: 'short',
                                      day: 'numeric',
                                      year: 'numeric',
                                  }).format(new Date(inspection.created_at))
                                : '—'

                            return (
                                <tr
                                    key={inspection.id}
                                    className="hover:bg-muted/30 transition-colors cursor-pointer"
                                    onClick={() => {
                                        router.push(`/inspections/${inspection.id}`)
                                    }}
                                >
                                    <td className="px-4 py-3 font-medium">
                                        {inspection.contract_number}
                                    </td>
                                    <td className="px-4 py-3">
                                        <StatusBadge status={inspection.status} />
                                    </td>
                                    <td className="px-4 py-3">
                                        <StageBadge stage={stage} />
                                    </td>
                                    <td className="px-4 py-3 text-muted-foreground">
                                        {formattedDate}
                                    </td>
                                    <td className="px-4 py-3 text-right">
                                        <Link
                                            href={`/inspections/${inspection.id}`}
                                            onClick={(e) => e.stopPropagation()}
                                            className="text-xs text-muted-foreground hover:text-foreground transition-colors"
                                        >
                                            View
                                        </Link>
                                    </td>
                                </tr>
                            )
                        })
                    )}
                </tbody>
            </table>
            {!isLoading && items.length > 0 && (
                <Pagination
                    page={page}
                    pageSize={pageSize}
                    total={total}
                    onPageChange={onPageChange}
                />
            )}
        </div>
    )
}
