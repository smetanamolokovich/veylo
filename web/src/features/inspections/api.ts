import { apiClient } from "@/lib/api-client"
import type { ListInspectionsResponse } from "./types"

export async function listInspections(
  page = 1,
  pageSize = 5
): Promise<ListInspectionsResponse> {
  return apiClient
    .get("api/v1/inspections", {
      searchParams: { page, page_size: pageSize },
    })
    .json<ListInspectionsResponse>()
}
