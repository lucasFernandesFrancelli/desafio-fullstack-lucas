import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { PriorityBadge } from "@/features/solicitations/PriorityBadge";

describe("PriorityBadge", () => {
  it("mostra 'Sem prioridade' quando priority é null", () => {
    render(<PriorityBadge priority={null} />);
    expect(screen.getByText("Sem prioridade")).toBeInTheDocument();
  });

  it("mostra o valor calculado para prioridade baixa", () => {
    render(<PriorityBadge priority={10} />);
    expect(screen.getByText("Prioridade 10")).toBeInTheDocument();
  });

  it("mostra o valor calculado para prioridade média", () => {
    render(<PriorityBadge priority={30} />);
    expect(screen.getByText("Prioridade 30")).toBeInTheDocument();
  });

  it("mostra o valor calculado para prioridade alta", () => {
    render(<PriorityBadge priority={100} />);
    expect(screen.getByText("Prioridade 100")).toBeInTheDocument();
  });
});
