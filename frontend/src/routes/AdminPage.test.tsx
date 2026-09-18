import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter } from "react-router";
import { describe, expect, it } from "vitest";

import { AdminPage } from "@/routes/AdminPage";
import { allUsers, categories, securityCategory } from "@/tests/fixtures";

function renderPage() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={client}>
      <MemoryRouter>
        <AdminPage />
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe("AdminPage", () => {
  it("mostra as categorias com seus aprovadores por padrão", async () => {
    renderPage();
    expect(await screen.findByText(securityCategory.name)).toBeInTheDocument();
    expect(screen.getByText(/1º aprovador: Ricardo Nogueira/)).toBeInTheDocument();
  });

  it("troca para a aba de usuários e lista as pessoas cadastradas", async () => {
    const user = userEvent.setup();
    renderPage();

    await screen.findByText(securityCategory.name);
    await user.click(screen.getByRole("tab", { name: "Usuários" }));

    const first = allUsers[0];
    expect(first).toBeDefined();
    expect(await screen.findByText(first?.email ?? "")).toBeInTheDocument();
  });

  it("renderiza um card por categoria retornada pela API", async () => {
    renderPage();
    for (const category of categories) {
      expect(await screen.findByText(category.name)).toBeInTheDocument();
    }
  });
});
