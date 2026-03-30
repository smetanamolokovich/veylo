export type SystemStage = 'ENTRY' | 'EVALUATION' | 'REVIEW' | 'FINAL'

export interface InspectionItem {
    id: string
    contract_number: string
    status: string
    created_at: string
}

export interface ListInspectionsResponse {
    items: InspectionItem[]
    total: number
    page: number
    page_size: number
}

// Default vehicle workflow: maps status name -> system stage
// Used for display purposes when backend does not return stage directly
export const DEFAULT_STATUS_STAGE_MAP: Record<string, SystemStage> = {
    new: 'ENTRY',
    damage_entered: 'ENTRY',
    damage_evaluated: 'EVALUATION',
    inspected: 'REVIEW',
    completed: 'FINAL',
}
