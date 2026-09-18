import type { ReactElement } from "react";

import { SolicitationCard } from "@/features/solicitations/SolicitationCard";
import { STATUS_ACCENT_TOP_BORDER_CLASSES, STATUS_LABELS, STATUS_TEXT_CLASSES } from "@/lib/constants";
import { cn } from "@/lib/utils";
import type { SolicitationSummary, Status } from "@/types/api";

interface KanbanColumnProps {
  status: Status;
  items: SolicitationSummary[];
}

export function KanbanColumn({ status, items }: KanbanColumnProps): ReactElement {
  return (
    <div
      className={cn(
        "flex w-72 shrink-0 flex-col gap-3 rounded-xl border-t-4 bg-muted/40 p-3",
        STATUS_ACCENT_TOP_BORDER_CLASSES[status],
      )}
      data-testid={`kanban-column-${status}`}
    >
      <div className="flex items-center justify-between px-1">
        <h2 className={cn("text-sm font-semibold", STATUS_TEXT_CLASSES[status])}>{STATUS_LABELS[status]}</h2>
        <span className="flex size-5 items-center justify-center rounded-full bg-background text-xs font-medium text-muted-foreground ring-1 ring-border">
          {items.length}
        </span>
      </div>

      <div className="flex flex-col gap-2">
        {items.length === 0 && (
          <p className="px-1 text-xs text-muted-foreground">Nenhuma solicitação aqui.</p>
        )}
        {items.map((item) => (
          <SolicitationCard key={item.id} solicitation={item} hideStatus />
        ))}
      </div>
    </div>
  );
}
