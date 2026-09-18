import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";

import { ProfileCard } from "@/features/auth/ProfileCard";
import { approverStep1User } from "@/tests/fixtures";
import type { UserProfile } from "@/types/api";

const profile: UserProfile = {
  ...approverStep1User,
  approverFor: [{ categoryId: "cat-1", categoryName: "Segurança do Trabalho", order: 1 }],
};

describe("ProfileCard", () => {
  it("mostra nome, papel, email e categorias em que é aprovador", () => {
    render(<ProfileCard profile={profile} onSelect={vi.fn()} />);

    expect(screen.getByText(profile.name)).toBeInTheDocument();
    expect(screen.getByText(profile.email)).toBeInTheDocument();
    expect(screen.getByText("Colaborador")).toBeInTheDocument();
    expect(screen.getByText("1º aprovador · Segurança do Trabalho")).toBeInTheDocument();
  });

  it("chama onSelect com o id do perfil ao clicar", async () => {
    const user = userEvent.setup();
    const onSelect = vi.fn();

    render(<ProfileCard profile={profile} onSelect={onSelect} />);
    await user.click(screen.getByText(profile.name));

    expect(onSelect).toHaveBeenCalledWith(profile.id);
  });

  it("gera as iniciais a partir de um nome com uma só palavra", () => {
    const singleNameProfile: UserProfile = { ...profile, name: "Ana" };
    render(<ProfileCard profile={singleNameProfile} onSelect={vi.fn()} />);
    expect(screen.getByText("A")).toBeInTheDocument();
  });

  it("não chama onSelect quando disabled", async () => {
    const user = userEvent.setup();
    const onSelect = vi.fn();

    render(<ProfileCard profile={profile} onSelect={onSelect} disabled />);
    await user.click(screen.getByText(profile.name));

    expect(onSelect).not.toHaveBeenCalled();
  });
});
