import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";

import { AnalysisPanel } from "@/features/solicitations/AnalysisPanel";

const emptyDefaults = { severity: null, urgency: null, trend: null, analysisNotes: "" };

describe("AnalysisPanel", () => {
  it("bloqueia a finalização quando as notas e o parecer estão vazios", async () => {
    const user = userEvent.setup();
    const onFinalize = vi.fn().mockResolvedValue(undefined);

    render(
      <AnalysisPanel
        defaultValues={emptyDefaults}
        onSavePartial={vi.fn()}
        onFinalize={onFinalize}
        isSavingPartial={false}
        isFinalizing={false}
      />,
    );

    await user.click(screen.getByRole("button", { name: "Finalizar" }));

    expect(await screen.findByText("Parecer é obrigatório")).toBeInTheDocument();
    expect(onFinalize).not.toHaveBeenCalled();
  });

  it("salva o parecer parcial mesmo com notas incompletas", async () => {
    const user = userEvent.setup();
    const onSavePartial = vi.fn().mockResolvedValue(undefined);

    render(
      <AnalysisPanel
        defaultValues={emptyDefaults}
        onSavePartial={onSavePartial}
        onFinalize={vi.fn()}
        isSavingPartial={false}
        isFinalizing={false}
      />,
    );

    await user.type(screen.getByLabelText("Parecer"), "Ainda em avaliação");
    await user.click(screen.getByRole("button", { name: "Salvar parecer" }));

    await waitFor(() => expect(onSavePartial).toHaveBeenCalledWith({ analysisNotes: "Ainda em avaliação" }));
  });

  it("finaliza com as três notas e o parecer preenchidos", async () => {
    const user = userEvent.setup();
    const onFinalize = vi.fn().mockResolvedValue(undefined);

    render(
      <AnalysisPanel
        defaultValues={{ severity: 4, urgency: 3, trend: 5, analysisNotes: "" }}
        onSavePartial={vi.fn()}
        onFinalize={onFinalize}
        isSavingPartial={false}
        isFinalizing={false}
      />,
    );

    await user.type(screen.getByLabelText("Parecer"), "Risco confirmado, priorizar tratativa");
    await user.click(screen.getByRole("button", { name: "Finalizar" }));

    await waitFor(() =>
      expect(onFinalize).toHaveBeenCalledWith({
        severity: 4,
        urgency: 3,
        trend: 5,
        analysisNotes: "Risco confirmado, priorizar tratativa",
      }),
    );
  });

  it("bloqueia a finalização quando falta apenas uma das notas", async () => {
    const user = userEvent.setup();
    const onFinalize = vi.fn().mockResolvedValue(undefined);

    render(
      <AnalysisPanel
        defaultValues={{ severity: 4, urgency: null, trend: 5, analysisNotes: "Parecer pronto" }}
        onSavePartial={vi.fn()}
        onFinalize={onFinalize}
        isSavingPartial={false}
        isFinalizing={false}
      />,
    );

    await user.click(screen.getByRole("button", { name: "Finalizar" }));

    expect(await screen.findByText("Obrigatório (1 a 5)")).toBeInTheDocument();
    expect(onFinalize).not.toHaveBeenCalled();
  });
});
