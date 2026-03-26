"use client"

import { useInspections } from "@/features/inspections/hooks/use-inspections"
import { DashboardEmptyState } from "@/features/inspections/components/dashboard-empty-state"

export default function DashboardPage() {
  const { data, isLoading } = useInspections(1, 5)

  if (isLoading) return null

  if (data?.total === 0) return <DashboardEmptyState />

  return (
    <div className="p-6 max-w-5xl">
      <div className="mb-8">
        <h1 className="text-2xl font-semibold tracking-tight">Dashboard</h1>
        <p className="text-muted-foreground text-sm mt-1">Your inspections overview.</p>
      </div>

      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        {[
          { label: "Total inspections", value: data?.total ?? "\u2014" },
          { label: "In progress", value: "\u2014" },
          { label: "Completed this month", value: "\u2014" },
          { label: "Pending review", value: "\u2014" },
        ].map((stat) => (
          <div key={stat.label} className="rounded-xl border bg-card p-4">
            <p className="text-xs text-muted-foreground font-medium">{stat.label}</p>
            <p className="text-2xl font-semibold mt-1">{stat.value}</p>
          </div>
        ))}
      </div>

      <div className="rounded-xl border bg-card">
        <div className="p-4 border-b">
          <h2 className="font-medium text-sm">Recent inspections</h2>
        </div>
        <div className="divide-y">
          {data?.items.map((inspection) => (
            <div key={inspection.id} className="p-4 flex items-center justify-between">
              <span className="text-sm font-medium">{inspection.contract_number}</span>
              <span className="text-xs text-muted-foreground">{inspection.status}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
