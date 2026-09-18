import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router";
import { describe, expect, it } from "vitest";

import { KanbanPage } from "@/routes/KanbanPage";
import { solicitationSummaries } from "@/tests/fixtures";

function renderPage() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={client}>
      <MemoryRouter>
        <KanbanPage />
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe("KanbanPage", () => {
  it("renderiza as mesmas solicitações da Lista organizadas em colunas", async () => {
    renderPage();

    for (const item of solicitationSummaries) {
      expect(await screen.findByText(item.title)).toBeInTheDocument();
    }

    expect(screen.getByText("Rascunho")).toBeInTheDocument();
    expect(screen.getByText("Em aprovação")).toBeInTheDocument();
    expect(screen.getByText("Em análise")).toBeInTheDocument();
  });
});
