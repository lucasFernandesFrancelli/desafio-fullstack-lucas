import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { http, HttpResponse } from "msw";
import { describe, expect, it, vi } from "vitest";

import { Toaster } from "@/components/ui/sonner";
import { EditApproversDialog } from "@/features/admin/EditApproversDialog";
import { server } from "@/tests/mocks/server";
import { allUsers, securityCategory } from "@/tests/fixtures";
import type { Category } from "@/types/api";

const BASE = "http://localhost:8080/api/v1";

function renderDialog(category: Category) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={client}>
      <EditApproversDialog category={category} users={allUsers} />
      <Toaster />
    </QueryClientProvider>,
  );
}

describe("EditApproversDialog", () => {
  it("mostra os aprovadores atuais já selecionados", async () => {
    const user = userEvent.setup();
    renderDialog(securityCategory);

    await user.click(screen.getByRole("button", { name: "Editar aprovadores" }));
    await screen.findByRole("dialog");

    const firstCombos = screen.getAllByRole("combobox", { name: "1º aprovador" });
    expect(firstCombos.some((el) => el.textContent?.includes("Ricardo Nogueira"))).toBe(true);

    const secondCombos = screen.getAllByRole("combobox", { name: "2º aprovador" });
    expect(secondCombos.some((el) => el.textContent?.includes("Fernanda Albuquerque"))).toBe(true);
  });

  it("bloqueia salvar quando as duas ordens apontam para a mesma pessoa", async () => {
    const user = userEvent.setup();
    const brokenCategory: Category = {
      id: "cat-quebrada",
      name: "Categoria Quebrada",
      description: "",
      approvers: [
        { userId: "user-approver-1", userName: "Ricardo Nogueira", order: 1 },
        { userId: "user-approver-1", userName: "Ricardo Nogueira", order: 2 },
      ],
    };
    renderDialog(brokenCategory);

    await user.click(screen.getByRole("button", { name: "Editar aprovadores" }));
    await user.click(screen.getByRole("button", { name: "Salvar" }));

    expect(await screen.findByText("Os dois aprovadores devem ser pessoas distintas")).toBeInTheDocument();
  });

  it("salva quando o servidor confirma a mudança", async () => {
    const setApproversHandler = vi.fn(() => HttpResponse.json(securityCategory));
    server.use(http.patch(`${BASE}/categories/:id/approvers`, setApproversHandler));

    const user = userEvent.setup();
    renderDialog(securityCategory);

    await user.click(screen.getByRole("button", { name: "Editar aprovadores" }));
    await user.click(screen.getByRole("button", { name: "Salvar" }));

    await waitFor(() => expect(setApproversHandler).toHaveBeenCalledTimes(1));
    expect(await screen.findByText("Aprovadores atualizados.")).toBeInTheDocument();
  });
});
