---
name: fe-feature
description: Scaffold a new frontend feature in web/features/<name>/ with components, hooks, api, schemas, types
argument-hint: <feature-name>
allowed-tools: Read, Glob, Grep, Write, Edit, Bash
---

Scaffold a new frontend feature for the Veylo web app.

Feature name: $ARGUMENTS

## Steps

1. Read `web/CLAUDE.md` for architecture rules
2. Check `web/features/` for existing features to follow the pattern
3. Check `web/lib/api-client.ts` to understand how to make API calls

## What to create

### `web/features/$ARGUMENTS/types.ts`
TypeScript types mirroring backend DTOs:
```ts
export interface $Entity {
  id: string
  organizationId: string
  // ... fields
  createdAt: string
}

export interface Create$EntityRequest {
  // ... request fields
}
```

### `web/features/$ARGUMENTS/schemas.ts`
Zod schemas for form validation:
```ts
import { z } from "zod"

export const create$EntitySchema = z.object({
  // ... fields with validation
})

export type Create$EntityFormData = z.infer<typeof create$EntitySchema>
```

### `web/features/$ARGUMENTS/api.ts`
API call functions using the ky client:
```ts
import { apiClient } from "@/lib/api-client"
import type { $Entity, Create$EntityRequest } from "./types"

export async function create$Entity(data: Create$EntityRequest): Promise<$Entity> {
  return apiClient.post("$arguments", { json: data }).json()
}

export async function get$Entity(id: string): Promise<$Entity> {
  return apiClient.get(`$arguments/${id}`).json()
}

export async function list$Entities(): Promise<{ items: $Entity[]; total: number }> {
  return apiClient.get("$arguments").json()
}
```

### `web/features/$ARGUMENTS/hooks/use-$arguments.ts`
TanStack Query hooks:
```ts
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { create$Entity, get$Entity, list$Entities } from "../api"
import type { Create$EntityRequest } from "../types"

const QUERY_KEYS = {
  all: ["$arguments"] as const,
  detail: (id: string) => ["$arguments", id] as const,
}

export function use$Entities() {
  return useQuery({
    queryKey: QUERY_KEYS.all,
    queryFn: list$Entities,
  })
}

export function use$Entity(id: string) {
  return useQuery({
    queryKey: QUERY_KEYS.detail(id),
    queryFn: () => get$Entity(id),
  })
}

export function useCreate$Entity() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: Create$EntityRequest) => create$Entity(data),
    onSuccess: () => qc.invalidateQueries({ queryKey: QUERY_KEYS.all }),
  })
}
```

### `web/features/$ARGUMENTS/components/$Entity-list.tsx`
Basic list component:
```tsx
"use client"

import { use$Entities } from "../hooks/use-$arguments"

export function $EntityList() {
  const { data, isLoading } = use$Entities()

  if (isLoading) return <div>Loading...</div>

  return (
    <ul>
      {data?.items.map((item) => (
        <li key={item.id}>{item.id}</li>
      ))}
    </ul>
  )
}
```

### `web/features/$ARGUMENTS/components/create-$arguments-form.tsx`
Form component:
```tsx
"use client"

import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { create$EntitySchema, type Create$EntityFormData } from "../schemas"
import { useCreate$Entity } from "../hooks/use-$arguments"

export function Create$EntityForm() {
  const { register, handleSubmit, formState: { errors } } = useForm<Create$EntityFormData>({
    resolver: zodResolver(create$EntitySchema),
  })
  const mutation = useCreate$Entity()

  const onSubmit = (data: Create$EntityFormData) => mutation.mutate(data)

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      {/* fields */}
      <button type="submit" disabled={mutation.isPending}>Create</button>
    </form>
  )
}
```

### `web/features/$ARGUMENTS/index.ts`
Barrel export:
```ts
export * from "./types"
export * from "./hooks/use-$arguments"
export * from "./components/$entity-list"
export * from "./components/create-$arguments-form"
```

## Rules
- Use `"use client"` only on components that use hooks or browser APIs
- Never import from other features — only from `@/lib/` and `@/components/ui/`
- All API types must match backend DTOs exactly (check `internal/application/` for response structs)
- After creating, verify TypeScript with: `cd web && npx tsc --noEmit`
