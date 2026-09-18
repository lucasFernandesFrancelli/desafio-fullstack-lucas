import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import { http, HttpResponse } from "msw";
import { MemoryRouter, Route, Routes } from "react-router";
import { describe, expect, it } from "vitest";

import { AuthProvider } from "@/context/AuthContext";
import { RequireAuth } from "@/routes/RequireAuth";
import { server } from "@/tests/mocks/server";
import { collaboratorUser } from "@/tests/fixtures";

const BASE = "http://localhost:8080/api/v1";

function renderProtected(initialEntry: string) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={client}>
      <MemoryRouter initialEntries={[initialEntry]}>
        <AuthProvider>
          <Routes>
            <Route path="/" element={<p>seletor de perfil</p>} />
            <Route
              path="/solicitations"
              element={
                <RequireAuth>
                  <p>área protegida</p>
                </RequireAuth>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe("RequireAuth", () => {
  it("redireciona para \"/\" quando não há usuário autenticado", async () => {
    renderProtected("/solicitations");

    expect(await screen.findByText("seletor de perfil")).toBeInTheDocument();
    expect(screen.queryByText("área protegida")).not.toBeInTheDocument();
  });

  it("libera o acesso quando a sessão é reidratada com sucesso", async () => {
    window.localStorage.setItem("ekaizen.token", "token-de-teste");
    server.use(http.get(`${BASE}/auth/me`, () => HttpResponse.json(collaboratorUser)));

    renderProtected("/solicitations");

    expect(await screen.findByText("área protegida")).toBeInTheDocument();
  });
});
