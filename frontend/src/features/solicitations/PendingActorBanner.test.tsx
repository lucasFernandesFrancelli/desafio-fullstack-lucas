import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { PendingActorBanner } from "@/features/solicitations/PendingActorBanner";

describe("PendingActorBanner", () => {
  it("não renderiza nada quando pendingActor.kind é 'none'", () => {
    const { container } = render(
      <PendingActorBanner status="finalizada" pendingActor={{ kind: "none", name: "" }} />,
    );
    expect(container).toBeEmptyDOMElement();
  });

  it("usa 'Aguardando conclusão de' para rascunho", () => {
    render(<PendingActorBanner status="rascunho" pendingActor={{ kind: "user", name: "Ana Beatriz Costa" }} />);
    expect(screen.getByText(/Aguardando conclusão de/)).toBeInTheDocument();
    expect(screen.getByText("Ana Beatriz Costa")).toBeInTheDocument();
  });

  it("usa 'Aguardando ação de' para os demais estados", () => {
    render(<PendingActorBanner status="em_aprovacao" pendingActor={{ kind: "user", name: "Ricardo Nogueira" }} />);
    expect(screen.getByText(/Aguardando ação de/)).toBeInTheDocument();
  });
});
