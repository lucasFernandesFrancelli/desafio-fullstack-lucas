import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";

import { SolicitationForm } from "@/features/solicitations/SolicitationForm";
import { emptyDraftFormValues } from "@/schemas/solicitationSchemas";
import { categories } from "@/tests/fixtures";

describe("SolicitationForm", () => {
  it("permite salvar rascunho apenas com o título preenchido", async () => {
    const user = userEvent.setup();
    const onSaveDraft = vi.fn().mockResolvedValue(undefined);
    const onSubmitForApproval = vi.fn().mockResolvedValue(undefined);

    render(
      <SolicitationForm
        defaultValues={emptyDraftFormValues}
        categories={categories}
        onSaveDraft={onSaveDraft}
        onSubmitForApproval={onSubmitForApproval}
        isSavingDraft={false}
        isSubmittingForApproval={false}
      />,
    );

    await user.type(screen.getByLabelText("Título"), "Piso escorregadio");
    await user.click(screen.getByRole("button", { name: "Salvar rascunho" }));

    await waitFor(() => expect(onSaveDraft).toHaveBeenCalledTimes(1));
    expect(onSaveDraft.mock.calls[0]?.[0]).toMatchObject({ title: "Piso escorregadio" });
    expect(onSubmitForApproval).not.toHaveBeenCalled();
  });

  it("bloqueia o salvamento de rascunho sem título", async () => {
    const user = userEvent.setup();
    const onSaveDraft = vi.fn().mockResolvedValue(undefined);

    render(
      <SolicitationForm
        defaultValues={emptyDraftFormValues}
        categories={categories}
        onSaveDraft={onSaveDraft}
        onSubmitForApproval={vi.fn()}
        isSavingDraft={false}
        isSubmittingForApproval={false}
      />,
    );

    await user.click(screen.getByRole("button", { name: "Salvar rascunho" }));

    expect(await screen.findByText("Título é obrigatório")).toBeInTheDocument();
    expect(onSaveDraft).not.toHaveBeenCalled();
  });

  it("bloqueia o envio para aprovação quando faltam campos obrigatórios", async () => {
    const user = userEvent.setup();
    const onSubmitForApproval = vi.fn().mockResolvedValue(undefined);

    render(
      <SolicitationForm
        defaultValues={{ ...emptyDraftFormValues, title: "Só o título" }}
        categories={categories}
        onSaveDraft={vi.fn()}
        onSubmitForApproval={onSubmitForApproval}
        isSavingDraft={false}
        isSubmittingForApproval={false}
      />,
    );

    await user.click(screen.getByRole("button", { name: "Enviar para aprovação" }));

    expect(await screen.findByText("Descrição do problema é obrigatória")).toBeInTheDocument();
    expect(await screen.findByText("Melhoria proposta é obrigatória")).toBeInTheDocument();
    expect(await screen.findByText("Categoria é obrigatória")).toBeInTheDocument();
    expect(await screen.findByText("Local é obrigatório")).toBeInTheDocument();
    expect(onSubmitForApproval).not.toHaveBeenCalled();
  });

  it("envia para aprovação quando todos os campos estão preenchidos", async () => {
    const user = userEvent.setup();
    const onSubmitForApproval = vi.fn().mockResolvedValue(undefined);

    render(
      <SolicitationForm
        defaultValues={{
          title: "Piso escorregadio",
          problemDescription: "Risco de queda no refeitório",
          proposedImprovement: "Instalar piso antiderrapante",
          categoryId: categories[0]?.id ?? "",
          location: "Refeitório",
        }}
        categories={categories}
        onSaveDraft={vi.fn()}
        onSubmitForApproval={onSubmitForApproval}
        isSavingDraft={false}
        isSubmittingForApproval={false}
      />,
    );

    await user.click(screen.getByRole("button", { name: "Enviar para aprovação" }));

    await waitFor(() => expect(onSubmitForApproval).toHaveBeenCalledTimes(1));
    expect(onSubmitForApproval.mock.calls[0]?.[0]).toMatchObject({
      title: "Piso escorregadio",
      proposedImprovement: "Instalar piso antiderrapante",
      location: "Refeitório",
    });
  });

  it("desabilita os campos em modo somente leitura e esconde os botões de ação", () => {
    render(
      <SolicitationForm
        defaultValues={emptyDraftFormValues}
        categories={categories}
        readOnly
        onSaveDraft={vi.fn()}
        onSubmitForApproval={vi.fn()}
        isSavingDraft={false}
        isSubmittingForApproval={false}
      />,
    );

    expect(screen.getByLabelText("Título")).toBeDisabled();
    expect(screen.queryByRole("button", { name: "Salvar rascunho" })).not.toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "Enviar para aprovação" })).not.toBeInTheDocument();
  });
});
