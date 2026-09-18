import { render, screen, within } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { MemoryRouter } from "react-router";

import { KanbanBoard } from "@/features/kanban/KanbanBoard";
import { STATUS_LABELS } from "@/lib/constants";
import type { SolicitationSummary } from "@/types/api";

function summary(overrides: Partial<SolicitationSummary>): SolicitationSummary {
  return {
    id: "id-1",
    title: "Item",
    status: "rascunho",
    currentApprovalStep: null,
    requesterName: "Fulano",
    priority: null,
    pendingActorName: "",
    createdAt: "2026-01-01T00:00:00Z",
    updatedAt: "2026-01-01T00:00:00Z",
    ...overrides,
  };
}

describe("KanbanBoard", () => {
  it("agrupa as solicitações na coluna correta de cada estado", () => {
    const items: SolicitationSummary[] = [
      summary({ id: "a", title: "Rascunho A", status: "rascunho" }),
      summary({ id: "b", title: "Em aprovação B", status: "em_aprovacao" }),
      summary({ id: "c", title: "Em análise C", status: "em_analise" }),
      summary({ id: "d", title: "Finalizada D", status: "finalizada" }),
      summary({ id: "e", title: "Recusada E", status: "recusada" }),
      summary({ id: "f", title: "Rascunho F", status: "rascunho" }),
    ];

    render(
      <MemoryRouter>
        <KanbanBoard items={items} />
      </MemoryRouter>,
    );

    const rascunhoColumn = screen.getByText(STATUS_LABELS.rascunho).closest("div");
    expect(rascunhoColumn).not.toBeNull();
    if (rascunhoColumn?.parentElement) {
      expect(within(rascunhoColumn.parentElement).getByText("Rascunho A")).toBeInTheDocument();
      expect(within(rascunhoColumn.parentElement).getByText("Rascunho F")).toBeInTheDocument();
      expect(within(rascunhoColumn.parentElement).queryByText("Em análise C")).not.toBeInTheDocument();
    }

    expect(screen.getByText("Em aprovação B")).toBeInTheDocument();
    expect(screen.getByText("Em análise C")).toBeInTheDocument();
    expect(screen.getByText("Finalizada D")).toBeInTheDocument();
    expect(screen.getByText("Recusada E")).toBeInTheDocument();
  });

  it("mostra a mensagem de coluna vazia quando não há itens naquele estado", () => {
    render(
      <MemoryRouter>
        <KanbanBoard items={[]} />
      </MemoryRouter>,
    );

    expect(screen.getAllByText("Nenhuma solicitação aqui.")).toHaveLength(5);
  });
});
