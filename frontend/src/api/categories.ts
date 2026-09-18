import { useQuery, type UseQueryResult } from "@tanstack/react-query";

import { httpClient } from "@/api/httpClient";
import type { Category } from "@/types/api";

function fetchCategories(): Promise<Category[]> {
  return httpClient.get<Category[]>("/categories");
}

export function useCategories(enabled: boolean): UseQueryResult<Category[], Error> {
  return useQuery({
    queryKey: ["categories"],
    queryFn: fetchCategories,
    enabled,
  });
}
