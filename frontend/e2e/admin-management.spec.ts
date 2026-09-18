import { expect, test } from "@playwright/test";

import { loginAs } from "./helpers";

test.describe.configure({ mode: "serial" });

const NEW_USER_NAME = `Pessoa E2E ${Date.now()}`;
const NEW_USER_EMAIL = `pessoa.e2e.${Date.now()}@ekaizen.example`;
const NEW_CATEGORY_NAME = `Ergonomia E2E ${Date.now()}`;

test("gestor gerencia usuários e categorias, e colaborador não acessa a área", async ({ page }) => {
  await loginAs(page, "Sérgio Andrade");
  await page.getByRole("link", { name: "Gestão" }).click();
  await expect(page.getByRole("heading", { name: "Gestão de cadastros" })).toBeVisible();

  // Cria um novo usuário.
  await page.getByRole("tab", { name: "Usuários" }).click();
  await page.getByRole("button", { name: "Novo usuário" }).click();
  await page.getByLabel("Nome").fill(NEW_USER_NAME);
  await page.getByLabel("E-mail").fill(NEW_USER_EMAIL);
  await page.getByRole("button", { name: "Criar usuário" }).click();
  await expect(page.getByText("Usuário criado.")).toBeVisible();
  await expect(page.getByRole("cell", { name: NEW_USER_NAME })).toBeVisible();

  // Cria uma categoria nova usando o usuário recém-criado como 1º aprovador.
  await page.getByRole("tab", { name: "Categorias e aprovadores" }).click();
  await page.getByRole("button", { name: "Nova categoria" }).click();
  await page.getByLabel("Nome").fill(NEW_CATEGORY_NAME);
  await page.getByLabel("1º aprovador").click();
  await page.getByRole("option", { name: NEW_USER_NAME }).click();
  await page.getByLabel("2º aprovador").click();
  await page.getByRole("option", { name: "Ana Beatriz Costa" }).click();
  await page.getByRole("button", { name: "Criar categoria" }).click();

  await expect(page.getByText("Categoria criada.")).toBeVisible();
  const categoryCard = page.locator('[data-slot="card"]').filter({ hasText: NEW_CATEGORY_NAME });
  await expect(categoryCard.getByText(`1º aprovador: ${NEW_USER_NAME}`)).toBeVisible();
  await expect(categoryCard.getByText("2º aprovador: Ana Beatriz Costa")).toBeVisible();

  // Reatribui os aprovadores da categoria recém-criada.
  await categoryCard.getByRole("button", { name: "Editar aprovadores" }).click();
  await page.getByLabel("2º aprovador").click();
  await page.getByRole("option", { name: "João Pedro Rocha" }).click();
  await page.getByRole("button", { name: "Salvar" }).click();
  await expect(page.getByText("Aprovadores atualizados.")).toBeVisible();
  await expect(categoryCard.getByText("2º aprovador: João Pedro Rocha")).toBeVisible();

  // A nova categoria já aparece disponível ao criar uma solicitação.
  await page.getByRole("button", { name: "Nova solicitação" }).click();
  await page.getByLabel("Categoria").click();
  await expect(page.getByRole("option", { name: NEW_CATEGORY_NAME })).toBeVisible();
  await page.keyboard.press("Escape");

  await page.getByRole("button", { name: "Sérgio Andrade" }).click();
  await page.getByRole("menuitem", { name: "Trocar de perfil" }).click();

  // Colaborador não tem acesso à área de gestão nem vê o link no menu.
  await loginAs(page, "Ana Beatriz Costa");
  await expect(page.getByRole("link", { name: "Gestão" })).not.toBeVisible();
  await page.goto("/admin");
  await expect(page).toHaveURL(/\/solicitations$/);
});
