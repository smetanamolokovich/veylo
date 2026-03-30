export function InspectionSkeletonRow() {
    return (
        <tr className="border-b">
            {[130, 100, 80, 80, 90].map((w, i) => (
                <td key={i} className="px-4 py-3">
                    <div className="h-4 rounded bg-muted animate-pulse" style={{ width: w }} />
                </td>
            ))}
        </tr>
    )
}
