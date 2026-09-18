import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { UsersTable } from "@/features/admin/UsersTable";
import { allUsers } from "@/tests/fixtures";

describe("UsersTable", () => {
  it("mostra mensagem quando não há usuários", () => {
    render(<UsersTable users={[]} />);
    expect(screen.getByText("Nenhum usuário cadastrado.")).toBeInTheDocument();
  });

  it("lista nome, e-mail e papel de cada usuário", () => {
    render(<UsersTable users={allUsers} />);
    const first = allUsers[0];
    expect(first).toBeDefined();
    expect(screen.getByText(first?.name ?? "")).toBeInTheDocument();
    expect(screen.getByText(first?.email ?? "")).toBeInTheDocument();
  });
});
