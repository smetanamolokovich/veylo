'use client'

import { useState } from 'react'
import Link from 'next/link'
import { buttonVariants } from '@/components/ui/button'
import { useInspections } from '@/features/inspections/hooks/use-inspections'
import { InspectionsTable } from '@/features/inspections/components/inspections-table'

const PAGE_SIZE = 20

export default function InspectionsPage() {
    const [page, setPage] = useState(1)
    const { data, isLoading, isError } = useInspections(page, PAGE_SIZE)

    return (
        <div className="p-6 w-full">
            {/* Header */}
            <div className="mb-6 flex items-center justify-between">
                <div>
                    <h1 className="text-2xl font-semibold tracking-tight">Inspections</h1>
                    <p className="text-muted-foreground text-sm mt-1">
                        {isLoading ? 'Loading…' : data ? `${data.total} total` : ''}
                    </p>
                </div>
                <Link href="/inspections/new" className={buttonVariants({ variant: 'default' })}>
                    New inspection
                </Link>
            </div>

            {/* Error state */}
            {isError && (
                <div className="rounded-xl border border-destructive/30 bg-destructive/5 px-4 py-3 text-sm text-destructive mb-4">
                    Failed to load inspections. Please try again.
                </div>
            )}

            {/* Table */}
            <InspectionsTable
                items={data?.items ?? []}
                total={data?.total ?? 0}
                page={page}
                pageSize={PAGE_SIZE}
                isLoading={isLoading}
                onPageChange={setPage}
            />
        </div>
    )
}
