import type { ReactElement } from "react";

import { SolicitationCard } from "@/features/solicitations/SolicitationCard";
import { STATUS_LABELS } from "@/lib/constants";
import type { SolicitationSummary, Status } from "@/types/api";

interface KanbanColumnProps {
  status: Status;
  items: SolicitationSummary[];
}

export function KanbanColumn({ status, items }: KanbanColumnProps): ReactElement {
  return (
    <div className="flex w-72 shrink-0 flex-col gap-3 rounded-xl bg-muted/40 p-3" data-testid={`kanban-column-${status}`}>
      <div className="flex items-center justify-between px-1">
        <h2 className="text-sm font-semibold">{STATUS_LABELS[status]}</h2>
        <span className="text-xs text-muted-foreground">{items.length}</span>
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
