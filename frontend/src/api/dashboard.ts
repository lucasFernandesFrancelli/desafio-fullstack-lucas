import { useQuery, type UseQueryResult } from "@tanstack/react-query";

import { httpClient } from "@/api/httpClient";
import type { DashboardSummary } from "@/types/api";

function fetchDashboard(): Promise<DashboardSummary> {
  return httpClient.get<DashboardSummary>("/dashboard");
}

export function useDashboard(enabled: boolean): UseQueryResult<DashboardSummary, Error> {
  return useQuery({
    queryKey: ["dashboard"],
    queryFn: fetchDashboard,
    enabled,
  });
}
