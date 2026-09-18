import { expect, test } from "@playwright/test";

import { loginAs, logout, openSolicitationByTitle } from "./helpers";

test.describe.configure({ mode: "serial" });

const TITLE = `EPI incorreto E2E ${Date.now()}`;
const API_BASE = "http://localhost:8081/api/v1";

test("esconde ações de quem não tem permissão e a API recusa a ação direta", async ({ page, request }) => {
  await loginAs(page, "Diego Farias");
  await page.getByRole("button", { name: "Nova solicitação" }).click();

  await page.getByLabel("Título").fill(TITLE);
  await page.getByLabel("Descrição do problema").fill("Luvas atuais não protegem contra o solvente.");
  await page.getByLabel("Melhoria proposta").fill("Trocar por luvas nitrílicas de cano longo.");
  await page.getByLabel("Categoria").click();
  await page.getByRole("option", { name: "Segurança do Trabalho" }).click();
  await page.getByLabel("Local").fill("Setor de Limpeza Industrial");
  await page.getByRole("button", { name: "Enviar para aprovação" }).click();
  await expect(page.getByText("Em aprovação", { exact: true })).toBeVisible();

  await logout(page, "Diego Farias");

  // Fernanda é a 2ª aprovadora — ainda não é a vez dela (etapa 1 = Ricardo).
  await loginAs(page, "Fernanda Albuquerque");
  await openSolicitationByTitle(page, TITLE);
  await expect(page.getByText("Ricardo Nogueira")).toBeVisible();
  await expect(page.getByRole("button", { name: "Aprovar" })).not.toBeVisible();
  await expect(page.getByRole("button", { name: "Recusar" })).not.toBeVisible();

  // Captura o token dela para testar a API diretamente, sem passar pela UI.
  const fernandaToken = await page.evaluate(() => window.localStorage.getItem("ekaizen.token"));
  const solicitationId = page.url().split("/solicitations/")[1];

  const directAttempt = await request.post(`${API_BASE}/solicitations/${solicitationId}/approve`, {
    headers: { Authorization: `Bearer ${fernandaToken}` },
  });
  expect(directAttempt.status()).toBe(403);

  await logout(page, "Fernanda Albuquerque");

  // Analista não deve ver o painel de análise enquanto a solicitação está em aprovação.
  await loginAs(page, "Tiago Martins");
  await openSolicitationByTitle(page, TITLE);
  await expect(page.getByText("Registre o parecer e as notas para finalizar a análise.")).not.toBeVisible();
});
