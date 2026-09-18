import { expect, test } from "@playwright/test";

import { loginAs, logout, openSolicitationByTitle, searchByTitle } from "./helpers";

test.describe.configure({ mode: "serial" });

const TITLE = `Piso escorregadio E2E ${Date.now()}`;

test("caminho completo: rascunho -> aprovação dupla -> análise -> finalizada", async ({ page }) => {
  // 1. Solicitante cria e envia
  await loginAs(page, "Ana Beatriz Costa");
  await page.getByRole("button", { name: "Nova solicitação" }).click();

  await page.getByLabel("Título").fill(TITLE);
  await page.getByRole("button", { name: "Salvar rascunho" }).click();
  await expect(page.getByText("Rascunho", { exact: true })).toBeVisible();

  await page.getByLabel("Descrição do problema").fill("Risco de queda no refeitório após o rodízio.");
  await page.getByLabel("Melhoria proposta").fill("Instalar piso antiderrapante e sinalização.");
  await page.getByLabel("Categoria").click();
  await page.getByRole("option", { name: "Segurança do Trabalho" }).click();
  await page.getByLabel("Local").fill("Refeitório Central");

  await page.getByRole("button", { name: "Enviar para aprovação" }).click();
  await expect(page.getByText("Em aprovação", { exact: true })).toBeVisible();
  await expect(page.getByText("Ricardo Nogueira")).toBeVisible();

  await logout(page, "Ana Beatriz Costa");

  // 2. Primeiro aprovador da categoria Segurança do Trabalho
  await loginAs(page, "Ricardo Nogueira");
  await openSolicitationByTitle(page, TITLE);
  await page.getByRole("button", { name: "Aprovar" }).click();
  await expect(page.getByText("Fernanda Albuquerque")).toBeVisible();

  await logout(page, "Ricardo Nogueira");

  // 3. Segundo aprovador -> libera a análise
  await loginAs(page, "Fernanda Albuquerque");
  await openSolicitationByTitle(page, TITLE);
  await page.getByRole("button", { name: "Aprovar" }).click();
  await expect(page.getByText("Em análise", { exact: true })).toBeVisible();
  await expect(page.getByText("Qualquer analista")).toBeVisible();

  await logout(page, "Fernanda Albuquerque");

  // 4. Analista registra parecer e finaliza
  await loginAs(page, "Renata Souza");
  await openSolicitationByTitle(page, TITLE);

  await page.getByLabel("Gravidade").click();
  await page.getByRole("option", { name: "4", exact: true }).click();
  await page.getByLabel("Urgência").click();
  await page.getByRole("option", { name: "3", exact: true }).click();
  await page.getByLabel("Tendência").click();
  await page.getByRole("option", { name: "5", exact: true }).click();
  await page.getByLabel("Parecer").fill("Risco confirmado; piso liso é causa raiz.");

  await page.getByRole("button", { name: "Finalizar" }).click();
  await expect(page.getByText("Finalizada", { exact: true })).toBeVisible();
  await expect(page.getByText("Prioridade 60")).toBeVisible();

  // 5. Paridade Lista/Kanban: o mesmo item aparece com o mesmo estado nas duas visões
  await page.goto("/solicitations");
  await searchByTitle(page, TITLE);
  const listCard = page.locator('[role="button"]').filter({ hasText: TITLE });
  await expect(listCard).toBeVisible();
  await expect(listCard.getByText("Finalizada", { exact: true })).toBeVisible();

  await page.getByRole("link", { name: "Kanban" }).click();
  await searchByTitle(page, TITLE);
  const finalizedColumn = page.getByTestId("kanban-column-finalizada");
  await expect(finalizedColumn.getByText(TITLE)).toBeVisible();
});
