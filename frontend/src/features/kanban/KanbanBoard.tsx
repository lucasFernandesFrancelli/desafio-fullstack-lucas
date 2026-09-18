import type { ReactElement } from "react";

import { KanbanColumn } from "@/features/kanban/KanbanColumn";
import { STATUS_ORDER } from "@/lib/constants";
import type { SolicitationSummary } from "@/types/api";

interface KanbanBoardProps {
  items: SolicitationSummary[];
}

export function KanbanBoard({ items }: KanbanBoardProps): ReactElement {
  return (
    <div className="flex gap-4 overflow-x-auto pb-2">
      {STATUS_ORDER.map((status) => (
        <KanbanColumn key={status} status={status} items={items.filter((item) => item.status === status)} />
      ))}
    </div>
  );
}
