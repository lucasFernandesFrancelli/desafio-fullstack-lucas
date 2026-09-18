import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { http, HttpResponse } from "msw";
import { MemoryRouter } from "react-router";
import { describe, expect, it } from "vitest";

import { ListPage } from "@/routes/ListPage";
import { server } from "@/tests/mocks/server";
import { solicitationSummaries } from "@/tests/fixtures";
import type { SolicitationSummary } from "@/types/api";

const BASE = "http://localhost:8080/api/v1";

function findPageIndicator(text: string) {
  return screen.getByText(
    (_content, element) => element?.tagName.toLowerCase() === "span" && element.textContent === text,
  );
}

function manySolicitations(count: number): SolicitationSummary[] {
  return Array.from({ length: count }, (_, index) => ({
    id: `solicitation-${index + 1}`,
    title: `Solicitação número ${index + 1}`,
    categoryName: "Qualidade",
    status: "em_analise",
    currentApprovalStep: null,
    requesterName: "Ana Beatriz Costa",
    priority: null,
    pendingActorName: "",
    createdAt: "2026-01-01T00:00:00Z",
    updatedAt: "2026-01-01T00:00:00Z",
  }));
}

function renderPage() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={client}>
      <MemoryRouter>
        <ListPage />
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe("ListPage", () => {
  it("lista as solicitações retornadas pela API", async () => {
    renderPage();

    for (const item of solicitationSummaries) {
      expect(await screen.findByText(item.title)).toBeInTheDocument();
    }
  });

  it("mostra mensagem de vazio quando não há solicitações", async () => {
    server.use(http.get(`${BASE}/solicitations`, () => HttpResponse.json([])));

    renderPage();

    expect(await screen.findByText("Nenhuma solicitação encontrada com esses filtros.")).toBeInTheDocument();
  });

  it("pagina a lista em blocos de 10 itens", async () => {
    const items = manySolicitations(12);
    server.use(http.get(`${BASE}/solicitations`, () => HttpResponse.json(items)));

    const user = userEvent.setup();
    renderPage();

    expect(await screen.findByText("Solicitação número 1")).toBeInTheDocument();
    expect(screen.getByText("Solicitação número 10")).toBeInTheDocument();
    expect(screen.queryByText("Solicitação número 11")).not.toBeInTheDocument();
    expect(findPageIndicator("Página 1 de 2")).toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: /próxima/i }));

    expect(await screen.findByText("Solicitação número 11")).toBeInTheDocument();
    expect(screen.getByText("Solicitação número 12")).toBeInTheDocument();
    expect(screen.queryByText("Solicitação número 1")).not.toBeInTheDocument();
    expect(screen.getByRole("button", { name: /próxima/i })).toBeDisabled();
  });

  it("não mostra paginação quando cabem todos os itens numa página", async () => {
    renderPage();

    await screen.findByText(solicitationSummaries[0]!.title);
    expect(screen.queryByRole("button", { name: /próxima/i })).not.toBeInTheDocument();
    expect(screen.queryByRole("button", { name: /anterior/i })).not.toBeInTheDocument();
  });

  it("mostra mensagem de erro quando a busca falha", async () => {
    server.use(
      http.get(`${BASE}/solicitations`, () =>
        HttpResponse.json({ error: { code: "internal_error", message: "erro" } }, { status: 500 }),
      ),
    );

    renderPage();

    expect(await screen.findByText("Não foi possível carregar as solicitações.")).toBeInTheDocument();
  });
});
