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
