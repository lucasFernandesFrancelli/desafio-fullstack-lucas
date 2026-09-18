import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { http, HttpResponse } from "msw";
import { describe, expect, it } from "vitest";

import { getStoredToken } from "@/api/httpClient";
import { AuthProvider, useAuth } from "@/context/AuthContext";
import { server } from "@/tests/mocks/server";
import { collaboratorUser } from "@/tests/fixtures";

const BASE = "http://localhost:8080/api/v1";

function TestConsumer() {
  const { user, isLoading, login, logout } = useAuth();

  if (isLoading) return <p>carregando</p>;

  return (
    <div>
      <p>{user ? `logado como ${user.name}` : "não autenticado"}</p>
      <button onClick={() => void login(collaboratorUser.id)}>entrar</button>
      <button onClick={logout}>sair</button>
    </div>
  );
}

function renderAuth() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={client}>
      <AuthProvider>
        <TestConsumer />
      </AuthProvider>
    </QueryClientProvider>,
  );
}

describe("AuthContext", () => {
  it("começa deslogado quando não há token salvo", async () => {
    renderAuth();
    expect(await screen.findByText("não autenticado")).toBeInTheDocument();
  });

  it("loga e persiste o token no localStorage", async () => {
    const user = userEvent.setup();
    renderAuth();

    await screen.findByText("não autenticado");
    await user.click(screen.getByRole("button", { name: "entrar" }));

    expect(await screen.findByText(`logado como ${collaboratorUser.name}`)).toBeInTheDocument();
    expect(getStoredToken()).toBe("fake-jwt-token");
  });

  it("desloga e limpa o token", async () => {
    const user = userEvent.setup();
    renderAuth();

    await screen.findByText("não autenticado");
    await user.click(screen.getByRole("button", { name: "entrar" }));
    await screen.findByText(`logado como ${collaboratorUser.name}`);

    await user.click(screen.getByRole("button", { name: "sair" }));

    expect(await screen.findByText("não autenticado")).toBeInTheDocument();
    expect(getStoredToken()).toBeNull();
  });

  it("reidrata a sessão a partir de um token salvo, buscando /auth/me", async () => {
    window.localStorage.setItem("ekaizen.token", "token-existente");

    renderAuth();

    expect(await screen.findByText(`logado como ${collaboratorUser.name}`)).toBeInTheDocument();
  });

  it("limpa o token quando /auth/me falha (token inválido)", async () => {
    server.use(
      http.get(`${BASE}/auth/me`, () =>
        HttpResponse.json({ error: { code: "unauthorized", message: "token inválido" } }, { status: 401 }),
      ),
    );
    window.localStorage.setItem("ekaizen.token", "token-invalido");

    renderAuth();

    await waitFor(() => expect(getStoredToken()).toBeNull());
    expect(await screen.findByText("não autenticado")).toBeInTheDocument();
  });
});
