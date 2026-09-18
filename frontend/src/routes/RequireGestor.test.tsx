import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import { http, HttpResponse } from "msw";
import { MemoryRouter, Route, Routes } from "react-router";
import { describe, expect, it } from "vitest";

import { AuthProvider } from "@/context/AuthContext";
import { RequireGestor } from "@/routes/RequireGestor";
import { server } from "@/tests/mocks/server";
import { collaboratorUser, managerUser } from "@/tests/fixtures";
import type { User } from "@/types/api";

const BASE = "http://localhost:8080/api/v1";

function renderAsUser(user: User) {
  window.localStorage.setItem("ekaizen.token", "token-de-teste");
  server.use(http.get(`${BASE}/auth/me`, () => HttpResponse.json(user)));

  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={client}>
      <MemoryRouter initialEntries={["/dashboard"]}>
        <AuthProvider>
          <Routes>
            <Route path="/solicitations" element={<p>lista</p>} />
            <Route
              path="/dashboard"
              element={
                <RequireGestor>
                  <p>dashboard do gestor</p>
                </RequireGestor>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe("RequireGestor", () => {
  it("permite acesso quando o usuário é gestor", async () => {
    renderAsUser(managerUser);
    expect(await screen.findByText("dashboard do gestor")).toBeInTheDocument();
  });

  it("redireciona colaboradores para a lista de solicitações", async () => {
    renderAsUser(collaboratorUser);
    expect(await screen.findByText("lista")).toBeInTheDocument();
    expect(screen.queryByText("dashboard do gestor")).not.toBeInTheDocument();
  });
});
