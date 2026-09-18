import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { http, HttpResponse } from "msw";
import { MemoryRouter, Route, Routes } from "react-router";
import { describe, expect, it, vi } from "vitest";

import { Toaster } from "@/components/ui/sonner";
import { SolicitationDetailPage } from "@/routes/SolicitationDetailPage";
import { server } from "@/tests/mocks/server";
import {
  completeDraftSolicitation,
  draftSolicitation,
  inAnalysisSolicitation,
  pendingApprovalSolicitation,
} from "@/tests/fixtures";
import type { SolicitationDetail } from "@/types/api";

const BASE = "http://localhost:8080/api/v1";

function renderDetail(detail: SolicitationDetail) {
  server.use(http.get(`${BASE}/solicitations/:id`, () => HttpResponse.json(detail)));

  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={client}>
      <MemoryRouter initialEntries={[`/solicitations/${detail.id}`]}>
        <Routes>
          <Route path="/solicitations/:id" element={<SolicitationDetailPage />} />
        </Routes>
        <Toaster />
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe("SolicitationDetailPage", () => {
  it("mostra o formulário editável e sem painéis de ação quando é o dono do rascunho", async () => {
    renderDetail(draftSolicitation);

    expect(await screen.findByLabelText("Título")).toBeEnabled();
    expect(screen.getByRole("button", { name: "Salvar rascunho" })).toBeInTheDocument();
    expect(screen.queryByText("É a sua vez de decidir esta solicitação.")).not.toBeInTheDocument();
    expect(screen.queryByText("Registre o parecer e as notas para finalizar a análise.")).not.toBeInTheDocument();
  });

  it("mostra o painel de aprovação apenas quando permissions.canApprove é true", async () => {
    renderDetail(pendingApprovalSolicitation);

    expect(await screen.findByText("É a sua vez de decidir esta solicitação.")).toBeInTheDocument();
    expect(screen.getByLabelText("Título")).toBeDisabled();
    expect(screen.queryByRole("button", { name: "Salvar rascunho" })).not.toBeInTheDocument();
  });

  it("mostra o painel de análise apenas quando permissions.canAnalyze é true", async () => {
    renderDetail(inAnalysisSolicitation);

    expect(await screen.findByText("Registre o parecer e as notas para finalizar a análise.")).toBeInTheDocument();
    expect(screen.queryByText("É a sua vez de decidir esta solicitação.")).not.toBeInTheDocument();
  });

  it("mostra quem precisa agir a partir de pendingActor", async () => {
    renderDetail(pendingApprovalSolicitation);

    expect(await screen.findByText("Ricardo Nogueira")).toBeInTheDocument();
  });

  it("aprova a solicitação ao clicar em Aprovar", async () => {
    const approveHandler = vi.fn(() => HttpResponse.json(pendingApprovalSolicitation));
    server.use(http.post(`${BASE}/solicitations/:id/approve`, approveHandler));

    const user = userEvent.setup();
    renderDetail(pendingApprovalSolicitation);

    await user.click(await screen.findByRole("button", { name: "Aprovar" }));

    await waitFor(() => expect(approveHandler).toHaveBeenCalledTimes(1));
    expect(await screen.findByText("Solicitação aprovada.")).toBeInTheDocument();
  });

  it("mostra erro ao falhar a aprovação", async () => {
    server.use(
      http.post(`${BASE}/solicitations/:id/approve`, () =>
        HttpResponse.json({ error: { code: "forbidden", message: "não é a sua vez" } }, { status: 403 }),
      ),
    );

    const user = userEvent.setup();
    renderDetail(pendingApprovalSolicitation);

    await user.click(await screen.findByRole("button", { name: "Aprovar" }));

    expect(await screen.findByText("não é a sua vez")).toBeInTheDocument();
  });

  it("recusa a solicitação com o motivo informado", async () => {
    const rejectHandler = vi.fn(() => HttpResponse.json(pendingApprovalSolicitation));
    server.use(http.post(`${BASE}/solicitations/:id/reject`, rejectHandler));

    const user = userEvent.setup();
    renderDetail(pendingApprovalSolicitation);

    await user.click(await screen.findByRole("button", { name: "Recusar" }));
    await user.type(screen.getByLabelText("Motivo"), "Fora do escopo");
    await user.click(screen.getByRole("button", { name: "Confirmar recusa" }));

    await waitFor(() => expect(rejectHandler).toHaveBeenCalledTimes(1));
    expect(await screen.findByText("Solicitação recusada.")).toBeInTheDocument();
  });

  it("salva o parecer parcial e finaliza a análise", async () => {
    const finalizeHandler = vi.fn(() => HttpResponse.json(inAnalysisSolicitation));
    server.use(http.post(`${BASE}/solicitations/:id/finalize`, finalizeHandler));

    const user = userEvent.setup();
    renderDetail(inAnalysisSolicitation);

    await user.type(await screen.findByLabelText("Parecer"), "Ainda em avaliação");
    await user.click(screen.getByRole("button", { name: "Salvar parecer" }));
    expect(await screen.findByText("Parecer salvo.")).toBeInTheDocument();
  });

  it("edita o rascunho e salva", async () => {
    const patchHandler = vi.fn(() => HttpResponse.json(draftSolicitation));
    server.use(http.patch(`${BASE}/solicitations/:id`, patchHandler));

    const user = userEvent.setup();
    renderDetail(draftSolicitation);

    await user.type(await screen.findByLabelText("Local"), "Refeitório");
    await user.click(screen.getByRole("button", { name: "Salvar rascunho" }));

    await waitFor(() => expect(patchHandler).toHaveBeenCalledTimes(1));
    expect(await screen.findByText("Rascunho salvo.")).toBeInTheDocument();
  });

  it("envia para aprovação quando todos os campos já estão preenchidos", async () => {
    const submitHandler = vi.fn(() => HttpResponse.json(completeDraftSolicitation));
    server.use(
      http.patch(`${BASE}/solicitations/:id`, () => HttpResponse.json(completeDraftSolicitation)),
      http.post(`${BASE}/solicitations/:id/submit`, submitHandler),
    );

    const user = userEvent.setup();
    renderDetail(completeDraftSolicitation);

    await screen.findByLabelText("Título");
    await user.click(screen.getByRole("button", { name: "Enviar para aprovação" }));

    await waitFor(() => expect(submitHandler).toHaveBeenCalledTimes(1));
    expect(await screen.findByText("Solicitação enviada para aprovação.")).toBeInTheDocument();
  });
});
