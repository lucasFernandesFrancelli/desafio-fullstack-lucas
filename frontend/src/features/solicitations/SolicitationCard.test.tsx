import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter, Route, Routes } from "react-router";
import { describe, expect, it } from "vitest";

import { SolicitationCard } from "@/features/solicitations/SolicitationCard";
import type { SolicitationSummary } from "@/types/api";

const summary: SolicitationSummary = {
  id: "sol-1",
  title: "Piso escorregadio",
  categoryName: "Segurança do Trabalho",
  status: "em_aprovacao",
  currentApprovalStep: 1,
  requesterName: "Ana Beatriz Costa",
  priority: null,
  pendingActorName: "Ricardo Nogueira",
  createdAt: "2026-01-01T00:00:00Z",
  updatedAt: "2026-01-01T00:00:00Z",
};

function renderCard(hideStatus = false) {
  return render(
    <MemoryRouter initialEntries={["/solicitations"]}>
      <Routes>
        <Route path="/solicitations" element={<SolicitationCard solicitation={summary} hideStatus={hideStatus} />} />
        <Route path="/solicitations/:id" element={<p>detalhe da solicitação</p>} />
      </Routes>
    </MemoryRouter>,
  );
}

describe("SolicitationCard", () => {
  it("mostra os dados principais da solicitação", () => {
    renderCard();
    expect(screen.getByText("Piso escorregadio")).toBeInTheDocument();
    expect(screen.getByText("Segurança do Trabalho")).toBeInTheDocument();
    expect(screen.getByText("Por Ana Beatriz Costa")).toBeInTheDocument();
    expect(screen.getByText(/Ricardo Nogueira/)).toBeInTheDocument();
  });

  it("esconde o status quando hideStatus é true (usado no Kanban)", () => {
    renderCard(true);
    expect(screen.queryByText("Em aprovação")).not.toBeInTheDocument();
  });

  it("navega para o detalhe ao clicar no card", async () => {
    const user = userEvent.setup();
    renderCard();

    await user.click(screen.getByText("Piso escorregadio"));

    expect(await screen.findByText("detalhe da solicitação")).toBeInTheDocument();
  });

  it("navega para o detalhe ao pressionar Enter com o card focado", async () => {
    const user = userEvent.setup();
    renderCard();

    screen.getByRole("button").focus();
    await user.keyboard("{Enter}");

    expect(await screen.findByText("detalhe da solicitação")).toBeInTheDocument();
  });
});
