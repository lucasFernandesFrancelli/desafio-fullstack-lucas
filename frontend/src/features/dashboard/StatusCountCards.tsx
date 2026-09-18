import type { ReactElement } from "react";

import { Card, CardContent } from "@/components/ui/card";
import { STATUS_ACCENT_TOP_BORDER_CLASSES, STATUS_LABELS, STATUS_ORDER, STATUS_TEXT_CLASSES } from "@/lib/constants";
import { cn } from "@/lib/utils";
import type { CountsByStatus } from "@/types/api";

interface StatusCountCardsProps {
  countsByStatus: CountsByStatus;
}

export function StatusCountCards({ countsByStatus }: StatusCountCardsProps): ReactElement {
  return (
    <div className="grid gap-3 sm:grid-cols-3 lg:grid-cols-5">
      {STATUS_ORDER.map((status) => (
        <Card key={status} className={cn("border-t-4", STATUS_ACCENT_TOP_BORDER_CLASSES[status])}>
          <CardContent className="flex flex-col gap-1 p-4">
            <span className="text-xs text-muted-foreground">{STATUS_LABELS[status]}</span>
            <span className={cn("text-2xl font-semibold tabular-nums", STATUS_TEXT_CLASSES[status])}>
              {countsByStatus[status] ?? 0}
            </span>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
