import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { HistoryTimeline } from "@/features/solicitations/HistoryTimeline";
import type { HistoryEntry } from "@/types/api";

describe("HistoryTimeline", () => {
  it("mostra mensagem quando não há histórico", () => {
    render(<HistoryTimeline history={[]} />);
    expect(screen.getByText("Ainda não há histórico.")).toBeInTheDocument();
  });

  it("mostra o motivo da recusa quando a ação é 'recusada'", () => {
    const history: HistoryEntry[] = [
      {
        id: "h1",
        solicitationId: "s1",
        fromStatus: "em_aprovacao",
        toStatus: "recusada",
        action: "recusada",
        actorId: "u1",
        actorName: "Fernanda Albuquerque",
        comment: "Fora do escopo",
        createdAt: "2026-01-01T00:00:00Z",
      },
    ];

    render(<HistoryTimeline history={history} />);

    expect(screen.getByText("Recusou a solicitação")).toBeInTheDocument();
    expect(screen.getByText("Motivo: Fora do escopo")).toBeInTheDocument();
  });

  it("mostra o parecer e as notas quando a ação é 'finalizada'", () => {
    const history: HistoryEntry[] = [
      {
        id: "h2",
        solicitationId: "s1",
        fromStatus: "em_analise",
        toStatus: "finalizada",
        action: "finalizada",
        actorId: "u2",
        actorName: "Renata Souza",
        comment: "Risco confirmado",
        severity: 4,
        urgency: 3,
        trend: 5,
        createdAt: "2026-01-01T00:00:00Z",
      },
    ];

    render(<HistoryTimeline history={history} />);

    expect(screen.getByText("Parecer: Risco confirmado")).toBeInTheDocument();
    expect(screen.getByText("Gravidade 4 · Urgência 3 · Tendência 5")).toBeInTheDocument();
  });
});
