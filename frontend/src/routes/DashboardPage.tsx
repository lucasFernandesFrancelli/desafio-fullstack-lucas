import type { ReactElement } from "react";

import { useDashboard } from "@/api/dashboard";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { AwaitingByUserTable } from "@/features/dashboard/AwaitingByUserTable";
import { OldestPendingTable } from "@/features/dashboard/OldestPendingTable";
import { StatusCountCards } from "@/features/dashboard/StatusCountCards";

export function DashboardPage(): ReactElement {
  const dashboardQuery = useDashboard(true);

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-1">
        <h1 className="text-xl font-semibold tracking-tight">Dashboard</h1>
        <p className="text-sm text-muted-foreground">
          Onde o processo está parado, e quem precisa agir — sem perguntar pessoa por pessoa.
        </p>
      </div>

      {dashboardQuery.isPending && (
        <div className="flex flex-col gap-4">
          <Skeleton className="h-24 w-full" />
          <Skeleton className="h-64 w-full" />
        </div>
      )}

      {dashboardQuery.isError && (
        <p className="text-sm text-destructive">Não foi possível carregar o dashboard.</p>
      )}

      {dashboardQuery.data && (
        <>
          <StatusCountCards countsByStatus={dashboardQuery.data.countsByStatus} />

          <div className="grid gap-6 lg:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle className="text-sm">Solicitações mais paradas</CardTitle>
              </CardHeader>
              <CardContent>
                <OldestPendingTable items={dashboardQuery.data.oldestPending} />
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="text-sm">Quem tem mais pendências</CardTitle>
              </CardHeader>
              <CardContent>
                <AwaitingByUserTable items={dashboardQuery.data.awaitingByUser} />
              </CardContent>
            </Card>
          </div>
        </>
      )}
    </div>
  );
}
