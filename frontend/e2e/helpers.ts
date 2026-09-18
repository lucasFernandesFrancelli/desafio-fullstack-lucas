import { expect, type Page } from "@playwright/test";

export async function loginAs(page: Page, name: string): Promise<void> {
  await page.goto("/");
  await page.getByText(name, { exact: true }).click();
  await expect(page).toHaveURL(/\/solicitations$/);
}

export async function logout(page: Page, name: string): Promise<void> {
  await page.getByRole("button", { name }).click();
  await page.getByRole("menuitem", { name: "Trocar de perfil" }).click();
  await expect(page).toHaveURL("http://localhost:5174/");
}

interface DraftFormData {
  title: string;
  problemDescription: string;
  proposedImprovement: string;
  categoryName: string;
  location: string;
}

/** Preenche o formulário de solicitação (Nova solicitação ou edição de rascunho) com todos os campos. */
export async function fillDraftForm(page: Page, data: DraftFormData): Promise<void> {
  await page.getByLabel("Título").fill(data.title);
  await page.getByLabel("Descrição do problema").fill(data.problemDescription);
  await page.getByLabel("Melhoria proposta").fill(data.proposedImprovement);
  await page.getByLabel("Categoria").click();
  await page.getByRole("option", { name: data.categoryName }).click();
  await page.getByLabel("Local").fill(data.location);
}

/** Seleciona um valor num Select do Radix pelo id/label do trigger e o texto da opção. */
export async function selectRadixOption(page: Page, triggerLabel: string, optionText: string): Promise<void> {
  await page.getByLabel(triggerLabel).click();
  await page.getByRole("option", { name: optionText, exact: true }).click();
}

/** A busca é debounced em 300ms no cliente antes de refazer a consulta. */
export async function searchByTitle(page: Page, title: string): Promise<void> {
  await page.getByPlaceholder(/Buscar por título/).fill(title);
  await page.waitForTimeout(400);
}

/** Busca por título na tela de Lista e abre o resultado (título é único por teste). */
export async function openSolicitationByTitle(page: Page, title: string): Promise<void> {
  await page.goto("/solicitations");
  await searchByTitle(page, title);
  await page.getByText(title, { exact: true }).click();
}
