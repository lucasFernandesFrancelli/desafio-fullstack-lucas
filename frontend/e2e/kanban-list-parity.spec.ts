import { expect, test } from "@playwright/test";

import { loginAs, searchByTitle } from "./helpers";

const TITLE = `Vazamento de óleo E2E ${Date.now()}`;

test("Lista e Kanban mostram o mesmo conjunto e abrem o mesmo detalhe", async ({ page }) => {
  await loginAs(page, "Patrícia Lemos");
  await page.getByRole("button", { name: "Nova solicitação" }).click();
  await page.getByLabel("Título").fill(TITLE);
  await page.getByRole("button", { name: "Salvar rascunho" }).click();
  await expect(page.getByText("Rascunho", { exact: true })).toBeVisible();

  // Lista: aparece com o filtro de busca.
  await page.goto("/solicitations");
  await searchByTitle(page, TITLE);
  const listCard = page.locator('[role="button"]').filter({ hasText: TITLE });
  await expect(listCard).toBeVisible();
  await expect(listCard.getByText("Rascunho", { exact: true })).toBeVisible();

  // Kanban: mesmo filtro, mesma solicitação, agora na coluna Rascunho.
  await page.getByRole("link", { name: "Kanban" }).click();
  await searchByTitle(page, TITLE);
  const rascunhoColumn = page.getByTestId("kanban-column-rascunho");
  const kanbanCard = rascunhoColumn.locator('[role="button"]').filter({ hasText: TITLE });
  await expect(kanbanCard).toBeVisible();

  // O card do Kanban não é arrastável: é um botão estático que abre o detalhe ao clicar.
  await expect(kanbanCard).toHaveAttribute("role", "button");
  await expect(kanbanCard).not.toHaveAttribute("draggable", "true");

  await kanbanCard.click();
  await expect(page).toHaveURL(/\/solicitations\/[^/]+$/);
  await expect(page.getByRole("heading", { name: TITLE })).toBeVisible();
});
