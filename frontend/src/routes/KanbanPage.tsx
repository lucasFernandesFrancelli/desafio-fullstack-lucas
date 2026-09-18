import { useState, type ReactElement } from "react";

import { useCategories } from "@/api/categories";
import type { SolicitationListFilter } from "@/api/solicitations";
import { useSolicitations } from "@/api/solicitations";
import { Skeleton } from "@/components/ui/skeleton";
import { FiltersBar } from "@/features/solicitations/FiltersBar";
import { KanbanBoard } from "@/features/kanban/KanbanBoard";

export function KanbanPage(): ReactElement {
  const [filter, setFilter] = useState<SolicitationListFilter>({});
  const categoriesQuery = useCategories(true);
  const solicitationsQuery = useSolicitations(filter);

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-1">
        <h1 className="text-xl font-semibold tracking-tight">Kanban</h1>
        <p className="text-sm text-muted-foreground">
          Mesmas solicitações da Lista, organizadas por estado. Clique num card para ver o detalhe e agir.
        </p>
      </div>

      <FiltersBar filter={filter} onChange={setFilter} categories={categoriesQuery.data ?? []} />

      {solicitationsQuery.isPending && (
        <div className="flex gap-4 overflow-x-auto">
          {Array.from({ length: 5 }, (_, index) => (
            <Skeleton key={`kanban-skeleton-${index}`} className="h-96 w-72 shrink-0 rounded-xl" />
          ))}
        </div>
      )}

      {solicitationsQuery.isError && (
        <p className="text-sm text-destructive">Não foi possível carregar as solicitações.</p>
      )}

      {solicitationsQuery.data && <KanbanBoard items={solicitationsQuery.data} />}
    </div>
  );
}
