import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it } from "vitest";

import { NewCategoryDialog } from "@/features/admin/NewCategoryDialog";
import { allUsers } from "@/tests/fixtures";

function renderDialog() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={client}>
      <NewCategoryDialog users={allUsers} />
    </QueryClientProvider>,
  );
}

describe("NewCategoryDialog", () => {
  it("bloqueia a criação sem nome e sem os dois aprovadores selecionados", async () => {
    const user = userEvent.setup();
    renderDialog();

    await user.click(screen.getByRole("button", { name: "Nova categoria" }));
    await user.click(screen.getByRole("button", { name: "Criar categoria" }));

    expect(await screen.findByText("Nome é obrigatório")).toBeInTheDocument();
    expect(await screen.findByText("Selecione o 1º aprovador")).toBeInTheDocument();
    expect(await screen.findByText("Selecione o 2º aprovador")).toBeInTheDocument();
  });
});
