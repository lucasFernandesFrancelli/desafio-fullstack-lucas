import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { CategoriesList } from "@/features/admin/CategoriesList";
import { allUsers, categories, securityCategory } from "@/tests/fixtures";
import type { Category } from "@/types/api";

function renderWithClient(categories: Category[], users = allUsers) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={client}>
      <CategoriesList categories={categories} users={users} />
    </QueryClientProvider>,
  );
}

describe("CategoriesList", () => {
  it("mostra mensagem quando não há categorias", () => {
    renderWithClient([]);
    expect(screen.getByText("Nenhuma categoria cadastrada.")).toBeInTheDocument();
  });

  it("mostra nome e os dois aprovadores de cada categoria", () => {
    renderWithClient(categories);
    expect(screen.getByText(securityCategory.name)).toBeInTheDocument();
    expect(screen.getByText(/1º aprovador: Ricardo Nogueira/)).toBeInTheDocument();
    expect(screen.getByText(/2º aprovador: Fernanda Albuquerque/)).toBeInTheDocument();
  });

  it("avisa quando um aprovador não está definido", () => {
    const incomplete: Category = { id: "cat-incompleta", name: "Sem aprovadores", description: "", approvers: [] };
    renderWithClient([incomplete]);
    expect(screen.getAllByText("não definido")).toHaveLength(2);
  });
});
