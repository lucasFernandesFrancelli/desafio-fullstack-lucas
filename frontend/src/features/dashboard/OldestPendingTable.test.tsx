import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router";
import { describe, expect, it } from "vitest";

import { OldestPendingTable } from "@/features/dashboard/OldestPendingTable";
import { dashboardSummary } from "@/tests/fixtures";

describe("OldestPendingTable", () => {
  it("mostra mensagem quando não há itens", () => {
    render(
      <MemoryRouter>
        <OldestPendingTable items={[]} />
      </MemoryRouter>,
    );
    expect(screen.getByText("Nenhuma solicitação parada no momento.")).toBeInTheDocument();
  });

  it("lista os itens recebidos", () => {
    render(
      <MemoryRouter>
        <OldestPendingTable items={dashboardSummary.oldestPending} />
      </MemoryRouter>,
    );
    const first = dashboardSummary.oldestPending[0];
    expect(first).toBeDefined();
    expect(screen.getByText(first?.title ?? "")).toBeInTheDocument();
  });
});
