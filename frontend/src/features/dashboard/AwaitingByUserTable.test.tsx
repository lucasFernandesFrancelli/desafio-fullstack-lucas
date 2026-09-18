import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { AwaitingByUserTable } from "@/features/dashboard/AwaitingByUserTable";
import { dashboardSummary } from "@/tests/fixtures";

describe("AwaitingByUserTable", () => {
  it("mostra mensagem quando não há pendências", () => {
    render(<AwaitingByUserTable items={[]} />);
    expect(screen.getByText("Ninguém com aprovações pendentes no momento.")).toBeInTheDocument();
  });

  it("lista as pessoas e suas contagens", () => {
    render(<AwaitingByUserTable items={dashboardSummary.awaitingByUser} />);
    const first = dashboardSummary.awaitingByUser[0];
    expect(first).toBeDefined();
    expect(screen.getByText(first?.userName ?? "")).toBeInTheDocument();
  });
});
