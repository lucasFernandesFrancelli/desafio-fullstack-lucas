import type { ReactElement } from "react";

import { Card, CardContent } from "@/components/ui/card";
import { STATUS_LABELS, STATUS_ORDER } from "@/lib/constants";
import type { CountsByStatus } from "@/types/api";

interface StatusCountCardsProps {
  countsByStatus: CountsByStatus;
}

export function StatusCountCards({ countsByStatus }: StatusCountCardsProps): ReactElement {
  return (
    <div className="grid gap-3 sm:grid-cols-3 lg:grid-cols-5">
      {STATUS_ORDER.map((status) => (
        <Card key={status}>
          <CardContent className="flex flex-col gap-1 p-4">
            <span className="text-xs text-muted-foreground">{STATUS_LABELS[status]}</span>
            <span className="text-2xl font-semibold tabular-nums">{countsByStatus[status] ?? 0}</span>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
