import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import { http, HttpResponse } from "msw";
import { MemoryRouter } from "react-router";
import { describe, expect, it } from "vitest";

import { ListPage } from "@/routes/ListPage";
import { server } from "@/tests/mocks/server";
import { solicitationSummaries } from "@/tests/fixtures";

const BASE = "http://localhost:8080/api/v1";

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
