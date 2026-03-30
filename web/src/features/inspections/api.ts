import { apiClient } from '@/lib/api-client'
import type { ListInspectionsResponse } from './types'

export async function listInspections(page = 1, pageSize = 20): Promise<ListInspectionsResponse> {
    console.log('[inspections] listInspections called', { page, pageSize })
    return apiClient
        .get('api/v1/inspections', {
            searchParams: { page, page_size: pageSize },
        })
        .json<ListInspectionsResponse>()
}
