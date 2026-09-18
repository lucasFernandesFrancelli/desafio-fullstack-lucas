import { useEffect, useMemo, useState, type ReactElement } from "react";

import { useCategories } from "@/api/categories";
import type { SolicitationListFilter } from "@/api/solicitations";
import { useSolicitations } from "@/api/solicitations";
import { Pagination } from "@/components/Pagination";
import { Skeleton } from "@/components/ui/skeleton";
import { FiltersBar } from "@/features/solicitations/FiltersBar";
import { SolicitationCard } from "@/features/solicitations/SolicitationCard";

const PAGE_SIZE = 10;

export function ListPage(): ReactElement {
  const [filter, setFilter] = useState<SolicitationListFilter>({});
  const [page, setPage] = useState(1);
  const categoriesQuery = useCategories(true);
  const solicitationsQuery = useSolicitations(filter);

  const items = solicitationsQuery.data;
  const totalPages = items ? Math.max(1, Math.ceil(items.length / PAGE_SIZE)) : 1;

  const pageItems = useMemo(() => {
    if (!items) return [];
    const start = (page - 1) * PAGE_SIZE;
    return items.slice(start, start + PAGE_SIZE);
  }, [items, page]);

  // Volta para a primeira página sempre que o filtro muda, para não ficar
  // numa página vazia depois de restringir a busca.
  useEffect(() => {
    setPage(1);
  }, [filter]);

  // Se a lista encolher (novo filtro, item resolvido) e a página atual deixar
  // de existir, recua para a última página válida.
  useEffect(() => {
    if (page > totalPages) {
      setPage(totalPages);
    }
  }, [page, totalPages]);

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

      {items && items.length === 0 && (
        <p className="text-sm text-muted-foreground">Nenhuma solicitação encontrada com esses filtros.</p>
      )}

      {items && items.length > 0 && (
        <>
          <p className="text-xs text-muted-foreground">
            {items.length} {items.length === 1 ? "solicitação encontrada" : "solicitações encontradas"}
          </p>

          <div className="flex flex-col gap-3">
            {pageItems.map((solicitation) => (
              <SolicitationCard key={solicitation.id} solicitation={solicitation} />
            ))}
          </div>

          <Pagination page={page} totalPages={totalPages} onPageChange={setPage} />
        </>
      )}
    </div>
  );
}
