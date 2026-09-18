import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { http, HttpResponse } from "msw";
import { describe, expect, it } from "vitest";

import { Toaster } from "@/components/ui/sonner";
import { NewUserDialog } from "@/features/admin/NewUserDialog";
import { server } from "@/tests/mocks/server";

const BASE = "http://localhost:8080/api/v1";

function renderDialog() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={client}>
      <NewUserDialog />
      <Toaster />
    </QueryClientProvider>,
  );
}

describe("NewUserDialog", () => {
  it("bloqueia a criação sem nome e e-mail válidos", async () => {
    const user = userEvent.setup();
    renderDialog();

    await user.click(screen.getByRole("button", { name: "Novo usuário" }));
    await user.click(screen.getByRole("button", { name: "Criar usuário" }));

    expect(await screen.findByText("Nome é obrigatório")).toBeInTheDocument();
    expect(await screen.findByText("E-mail é obrigatório")).toBeInTheDocument();
  });

  it("cria o usuário com os dados preenchidos", async () => {
    const user = userEvent.setup();
    renderDialog();

    await user.click(screen.getByRole("button", { name: "Novo usuário" }));
    await user.type(screen.getByLabelText("Nome"), "Pessoa Nova");
    await user.type(screen.getByLabelText("E-mail"), "pessoa.nova@ekaizen.example");
    await user.click(screen.getByRole("button", { name: "Criar usuário" }));

    expect(await screen.findByText("Usuário criado.")).toBeInTheDocument();
  });

  it("mostra erro quando o e-mail já existe", async () => {
    server.use(
      http.post(`${BASE}/users`, () =>
        HttpResponse.json({ error: { code: "validation_error", message: "e-mail já cadastrado" } }, { status: 422 }),
      ),
    );
    const user = userEvent.setup();
    renderDialog();

    await user.click(screen.getByRole("button", { name: "Novo usuário" }));
    await user.type(screen.getByLabelText("Nome"), "Pessoa Nova");
    await user.type(screen.getByLabelText("E-mail"), "duplicado@ekaizen.example");
    await user.click(screen.getByRole("button", { name: "Criar usuário" }));

    expect(await screen.findByText("e-mail já cadastrado")).toBeInTheDocument();
  });
});
