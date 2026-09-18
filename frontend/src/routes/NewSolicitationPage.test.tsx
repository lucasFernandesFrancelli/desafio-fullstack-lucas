import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { http, HttpResponse } from "msw";
import { MemoryRouter, Route, Routes } from "react-router";
import { describe, expect, it } from "vitest";

import { Toaster } from "@/components/ui/sonner";
import { NewSolicitationPage } from "@/routes/NewSolicitationPage";
import { server } from "@/tests/mocks/server";
import { draftSolicitation } from "@/tests/fixtures";

const BASE = "http://localhost:8080/api/v1";

function renderPage() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={client}>
      <MemoryRouter initialEntries={["/solicitations/new"]}>
        <Routes>
          <Route path="/solicitations/new" element={<NewSolicitationPage />} />
          <Route path="/solicitations/:id" element={<p>detalhe: {draftSolicitation.id}</p>} />
        </Routes>
        <Toaster />
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe("NewSolicitationPage", () => {
  it("cria o rascunho e navega para o detalhe ao salvar", async () => {
    const user = userEvent.setup();
    renderPage();

    await user.type(await screen.findByLabelText("Título"), "Piso escorregadio");
    await user.click(screen.getByRole("button", { name: "Salvar rascunho" }));

    expect(await screen.findByText(`detalhe: ${draftSolicitation.id}`, {}, { timeout: 3000 })).toBeInTheDocument();
  });

  it("mostra erro quando a criação falha", async () => {
    server.use(
      http.post(`${BASE}/solicitations`, () =>
        HttpResponse.json({ error: { code: "internal_error", message: "Erro ao criar" } }, { status: 500 }),
      ),
    );
    const user = userEvent.setup();
    renderPage();

    await user.type(await screen.findByLabelText("Título"), "Piso escorregadio");
    await user.click(screen.getByRole("button", { name: "Salvar rascunho" }));

    expect(await screen.findByText("Erro ao criar", {}, { timeout: 3000 })).toBeInTheDocument();
  });
});
