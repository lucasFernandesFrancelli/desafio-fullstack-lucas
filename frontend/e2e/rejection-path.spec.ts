import { expect, test } from "@playwright/test";

import { loginAs, logout, openSolicitationByTitle } from "./helpers";

test.describe.configure({ mode: "serial" });

const TITLE = `Retrabalho na embalagem E2E ${Date.now()}`;
const REASON = "Fora do escopo de melhoria de processo — sugerido abrir chamado de facilities.";

test("recusa encerra o fluxo definitivamente e registra o motivo", async ({ page }) => {
  await loginAs(page, "João Pedro Rocha");
  await page.getByRole("button", { name: "Nova solicitação" }).click();

  await page.getByLabel("Título").fill(TITLE);
  await page.getByLabel("Descrição do problema").fill("Itens saem com etiqueta trocada.");
  await page.getByLabel("Melhoria proposta").fill("Checklist de dupla checagem.");
  await page.getByLabel("Categoria").click();
  await page.getByRole("option", { name: "Qualidade" }).click();
  await page.getByLabel("Local").fill("Linha de Embalagem 3");

  await page.getByRole("button", { name: "Enviar para aprovação" }).click();
  await expect(page.getByText("Em aprovação", { exact: true })).toBeVisible();

  await logout(page, "João Pedro Rocha");

  await loginAs(page, "Camila Duarte");
  await openSolicitationByTitle(page, TITLE);

  await page.getByRole("button", { name: "Recusar" }).click();
  await page.getByLabel("Motivo").fill(REASON);
  await page.getByRole("button", { name: "Confirmar recusa" }).click();

  await expect(page.getByText("Recusada", { exact: true })).toBeVisible();
  await expect(page.getByText(`Motivo da recusa: ${REASON}`)).toBeVisible();

  // Estado terminal: nenhuma ação de aprovação/recusa deve continuar disponível.
  await expect(page.getByRole("button", { name: "Aprovar" })).not.toBeVisible();
  await expect(page.getByRole("button", { name: "Recusar" })).not.toBeVisible();

  await logout(page, "Camila Duarte");

  // A recusa é visível para o gestor na consulta, mesmo sendo estado terminal.
  await loginAs(page, "Sérgio Andrade");
  await openSolicitationByTitle(page, TITLE);
  await expect(page.getByText("Recusada", { exact: true })).toBeVisible();
});
