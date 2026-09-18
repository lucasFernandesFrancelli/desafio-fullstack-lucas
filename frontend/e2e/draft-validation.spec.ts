import { expect, test } from "@playwright/test";

import { loginAs } from "./helpers";

test("bloqueia salvar rascunho sem título e envio sem os demais campos obrigatórios", async ({ page }) => {
  await loginAs(page, "Larissa Mendes");
  await page.getByRole("button", { name: "Nova solicitação" }).click();

  // Salvar rascunho sem nada preenchido.
  await page.getByRole("button", { name: "Salvar rascunho" }).click();
  await expect(page.getByText("Título é obrigatório")).toBeVisible();

  // Preenche só o título e tenta enviar direto para aprovação.
  await page.getByLabel("Título").fill("Sugestão incompleta de propósito");
  await page.getByRole("button", { name: "Enviar para aprovação" }).click();

  await expect(page.getByText("Descrição do problema é obrigatória")).toBeVisible();
  await expect(page.getByText("Melhoria proposta é obrigatória")).toBeVisible();
  await expect(page.getByText("Categoria é obrigatória")).toBeVisible();
  await expect(page.getByText("Local é obrigatório")).toBeVisible();

  // Nenhuma navegação ocorreu — segue na tela de criação.
  await expect(page).toHaveURL(/\/solicitations\/new$/);
});
