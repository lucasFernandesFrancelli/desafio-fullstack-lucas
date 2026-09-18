import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router";
import { describe, expect, it } from "vitest";

import { DashboardPage } from "@/routes/DashboardPage";
import { dashboardSummary } from "@/tests/fixtures";

function renderPage() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={client}>
      <MemoryRouter>
        <DashboardPage />
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe("DashboardPage", () => {
  it("mostra as contagens por estado e as tabelas de acompanhamento", async () => {
    renderPage();

    expect(await screen.findByText("Solicitações mais paradas")).toBeInTheDocument();
    expect(screen.getByText(dashboardSummary.oldestPending[0]?.title ?? "")).toBeInTheDocument();
    expect(screen.getByText("Quem tem mais pendências")).toBeInTheDocument();
    expect(screen.getByText(dashboardSummary.awaitingByUser[0]?.userName ?? "")).toBeInTheDocument();
  });
});
