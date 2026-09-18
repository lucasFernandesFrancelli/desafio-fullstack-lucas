import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";

import { ApproveRejectPanel } from "@/features/solicitations/ApproveRejectPanel";

describe("ApproveRejectPanel", () => {
  it("chama onApprove ao clicar em Aprovar", async () => {
    const user = userEvent.setup();
    const onApprove = vi.fn().mockResolvedValue(undefined);

    render(
      <ApproveRejectPanel onApprove={onApprove} onReject={vi.fn()} isApproving={false} isRejecting={false} />,
    );

    await user.click(screen.getByRole("button", { name: "Aprovar" }));

    await waitFor(() => expect(onApprove).toHaveBeenCalledTimes(1));
  });

  it("exige motivo antes de confirmar a recusa", async () => {
    const user = userEvent.setup();
    const onReject = vi.fn().mockResolvedValue(undefined);

    render(<ApproveRejectPanel onApprove={vi.fn()} onReject={onReject} isApproving={false} isRejecting={false} />);

    await user.click(screen.getByRole("button", { name: "Recusar" }));
    await user.click(screen.getByRole("button", { name: "Confirmar recusa" }));

    expect(await screen.findByText("Motivo da recusa é obrigatório")).toBeInTheDocument();
    expect(onReject).not.toHaveBeenCalled();
  });

  it("chama onReject com o motivo preenchido", async () => {
    const user = userEvent.setup();
    const onReject = vi.fn().mockResolvedValue(undefined);

    render(<ApproveRejectPanel onApprove={vi.fn()} onReject={onReject} isApproving={false} isRejecting={false} />);

    await user.click(screen.getByRole("button", { name: "Recusar" }));
    await user.type(screen.getByLabelText("Motivo"), "Fora do escopo de melhoria de processo");
    await user.click(screen.getByRole("button", { name: "Confirmar recusa" }));

    await waitFor(() => expect(onReject).toHaveBeenCalledWith("Fora do escopo de melhoria de processo"));
  });
});
