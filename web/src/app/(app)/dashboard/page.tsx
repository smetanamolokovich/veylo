export default function DashboardPage() {
  return (
    <div className="p-6 max-w-5xl">
      {/* Page header */}
      <div className="mb-8">
        <h1 className="text-2xl font-semibold tracking-tight">Dashboard</h1>
        <p className="text-muted-foreground text-sm mt-1">Welcome to Veylo. Your inspections overview.</p>
      </div>

      {/* Stats row */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        {[
          { label: "Total inspections", value: "\u2014" },
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

      {/* Recent inspections placeholder */}
      <div className="rounded-xl border bg-card">
        <div className="p-4 border-b">
          <h2 className="font-medium text-sm">Recent inspections</h2>
        </div>
        <div className="p-10 flex flex-col items-center justify-center text-center">
          <div className="w-12 h-12 rounded-full bg-muted flex items-center justify-center mb-3">
            <svg width="20" height="20" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M3 2h10a1 1 0 0 1 1 1v10a1 1 0 0 1-1 1H3a1 1 0 0 1-1-1V3a1 1 0 0 1 1-1z" stroke="currentColor" strokeWidth="1.5"/>
              <path d="M5 5h6M5 8h6M5 11h4" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round"/>
            </svg>
          </div>
          <p className="text-sm font-medium">No inspections yet</p>
          <p className="text-xs text-muted-foreground mt-1">Create your first inspection to get started</p>
        </div>
      </div>
    </div>
  )
}
