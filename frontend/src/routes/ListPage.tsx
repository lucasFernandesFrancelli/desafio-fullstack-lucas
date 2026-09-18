import { useState, type ReactElement } from "react";

import { useCategories } from "@/api/categories";
import type { SolicitationListFilter } from "@/api/solicitations";
import { useSolicitations } from "@/api/solicitations";
import { Skeleton } from "@/components/ui/skeleton";
import { FiltersBar } from "@/features/solicitations/FiltersBar";
import { SolicitationCard } from "@/features/solicitations/SolicitationCard";

export function ListPage(): ReactElement {
  const [filter, setFilter] = useState<SolicitationListFilter>({});
  const categoriesQuery = useCategories(true);
  const solicitationsQuery = useSolicitations(filter);

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-1">
        <h1 className="text-xl font-semibold tracking-tight">Solicitações</h1>
        <p className="text-sm text-muted-foreground">
          Acompanhe todas as solicitações de melhoria, do rascunho à decisão final.
        </p>
      </div>

      <FiltersBar filter={filter} onChange={setFilter} categories={categoriesQuery.data ?? []} />

      {solicitationsQuery.isPending && (
        <div className="flex flex-col gap-3">
          {Array.from({ length: 4 }, (_, index) => (
            <Skeleton key={`solicitation-skeleton-${index}`} className="h-28 w-full rounded-xl" />
          ))}
        </div>
      )}

      {solicitationsQuery.isError && (
        <p className="text-sm text-destructive">Não foi possível carregar as solicitações.</p>
      )}

      {solicitationsQuery.data && solicitationsQuery.data.length === 0 && (
        <p className="text-sm text-muted-foreground">Nenhuma solicitação encontrada com esses filtros.</p>
      )}

      {solicitationsQuery.data && solicitationsQuery.data.length > 0 && (
        <div className="flex flex-col gap-3">
          {solicitationsQuery.data.map((solicitation) => (
            <SolicitationCard key={solicitation.id} solicitation={solicitation} />
          ))}
        </div>
      )}
    </div>
  );
}
