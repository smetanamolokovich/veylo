'use client'

import { useQuery } from '@tanstack/react-query'
import { listInspections } from '../api'

export function useInspections(page = 1, pageSize = 20) {
    return useQuery({
        queryKey: ['inspections', { page, pageSize }],
        queryFn: () => listInspections(page, pageSize),
    })
}
