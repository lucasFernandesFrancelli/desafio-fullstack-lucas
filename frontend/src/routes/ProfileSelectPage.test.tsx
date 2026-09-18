import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { http, HttpResponse } from "msw";
import { MemoryRouter, Route, Routes } from "react-router";
import { describe, expect, it } from "vitest";

import { AuthProvider } from "@/context/AuthContext";
import { ProfileSelectPage } from "@/routes/ProfileSelectPage";
import { server } from "@/tests/mocks/server";
import { profiles } from "@/tests/fixtures";

const BASE = "http://localhost:8080/api/v1";

function renderPage() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={client}>
      <MemoryRouter initialEntries={["/"]}>
        <AuthProvider>
          <Routes>
            <Route path="/" element={<ProfileSelectPage />} />
            <Route path="/solicitations" element={<p>tela de solicitações</p>} />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe("ProfileSelectPage", () => {
  it("lista os perfis disponíveis vindos da API", async () => {
    renderPage();

    for (const profile of profiles) {
      expect(await screen.findByText(profile.name)).toBeInTheDocument();
    }
  });

  it("mostra mensagem de erro quando a API de perfis falha", async () => {
    server.use(
      http.get(`${BASE}/auth/profiles`, () =>
        HttpResponse.json({ error: { code: "internal_error", message: "erro" } }, { status: 500 }),
      ),
    );

    renderPage();

    expect(
      await screen.findByText(/Não foi possível carregar os perfis/i),
    ).toBeInTheDocument();
  });

  it("loga e navega para a tela de solicitações ao escolher um perfil", async () => {
    const user = userEvent.setup();
    renderPage();

    const firstProfile = profiles[0];
    if (!firstProfile) throw new Error("fixture profiles está vazia");

    await user.click(await screen.findByText(firstProfile.name));

    expect(await screen.findByText("tela de solicitações")).toBeInTheDocument();
  });
});
