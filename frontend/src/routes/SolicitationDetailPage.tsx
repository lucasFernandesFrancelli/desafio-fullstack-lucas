import type { ReactElement } from "react";
import { useParams } from "react-router";
import { toast } from "sonner";

import { useCategories } from "@/api/categories";
import { ApiError } from "@/api/httpClient";
import {
  useApproveSolicitation,
  useFinalizeSolicitation,
  usePatchAnalysis,
  useRejectSolicitation,
  useSolicitationDetail,
  useSubmitSolicitation,
  useUpdateDraft,
} from "@/api/solicitations";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";
import { AnalysisPanel } from "@/features/solicitations/AnalysisPanel";
import { ApproveRejectPanel } from "@/features/solicitations/ApproveRejectPanel";
import { HistoryTimeline } from "@/features/solicitations/HistoryTimeline";
import { PendingActorBanner } from "@/features/solicitations/PendingActorBanner";
import { PriorityBadge } from "@/features/solicitations/PriorityBadge";
import { SolicitationForm } from "@/features/solicitations/SolicitationForm";
import { StatusBadge } from "@/features/solicitations/StatusBadge";
import type { DraftFormValues } from "@/schemas/solicitationSchemas";
import type { AnalysisInput, DraftInput } from "@/types/api";

function errorMessage(error: unknown, fallback: string): string {
  return error instanceof ApiError ? error.message : fallback;
}

export function SolicitationDetailPage(): ReactElement {
  const { id } = useParams<{ id: string }>();
  const detailQuery = useSolicitationDetail(id);
  const categoriesQuery = useCategories(true);

  const updateDraft = useUpdateDraft();
  const submitSolicitation = useSubmitSolicitation();
  const approveSolicitation = useApproveSolicitation();
  const rejectSolicitation = useRejectSolicitation();
  const patchAnalysis = usePatchAnalysis();
  const finalizeSolicitation = useFinalizeSolicitation();

  if (!id) {
    return <p className="text-sm text-destructive">Solicitação inválida.</p>;
  }

  if (detailQuery.isPending || categoriesQuery.isPending) {
    return (
      <div className="flex flex-col gap-4">
        <Skeleton className="h-8 w-2/3" />
        <Skeleton className="h-40 w-full" />
        <Skeleton className="h-40 w-full" />
      </div>
    );
  }

  if (detailQuery.isError) {
    return (
      <p className="text-sm text-destructive">
        Não foi possível carregar esta solicitação (ela pode não existir, ou você não tem permissão para vê-la).
      </p>
    );
  }

  if (categoriesQuery.isError) {
    return <p className="text-sm text-destructive">Não foi possível carregar as categorias. Recarregue a página.</p>;
  }

  const detail = detailQuery.data;
  const solicitationId = id;

  function toDraftInput(values: DraftFormValues): DraftInput {
    return {
      title: values.title,
      problemDescription: values.problemDescription,
      proposedImprovement: values.proposedImprovement,
      categoryId: values.categoryId || undefined,
      location: values.location,
    };
  }

  async function handleSaveDraft(values: DraftFormValues): Promise<void> {
    try {
      await updateDraft.mutateAsync({ id: solicitationId, input: toDraftInput(values) });
      toast.success("Rascunho salvo.");
    } catch (error) {
      toast.error(errorMessage(error, "Não foi possível salvar o rascunho."));
    }
  }

  async function handleSubmitForApproval(values: DraftFormValues): Promise<void> {
    try {
      await updateDraft.mutateAsync({ id: solicitationId, input: toDraftInput(values) });
      await submitSolicitation.mutateAsync(solicitationId);
      toast.success("Solicitação enviada para aprovação.");
    } catch (error) {
      toast.error(errorMessage(error, "Não foi possível enviar a solicitação."));
    }
  }

  async function handleApprove(): Promise<void> {
    try {
      await approveSolicitation.mutateAsync(solicitationId);
      toast.success("Solicitação aprovada.");
    } catch (error) {
      toast.error(errorMessage(error, "Não foi possível aprovar a solicitação."));
    }
  }

  async function handleReject(reason: string): Promise<void> {
    try {
      await rejectSolicitation.mutateAsync({ id: solicitationId, reason });
      toast.success("Solicitação recusada.");
    } catch (error) {
      toast.error(errorMessage(error, "Não foi possível recusar a solicitação."));
    }
  }

  async function handleSaveAnalysisPartial(input: AnalysisInput): Promise<void> {
    try {
      await patchAnalysis.mutateAsync({ id: solicitationId, input });
      toast.success("Parecer salvo.");
    } catch (error) {
      toast.error(errorMessage(error, "Não foi possível salvar o parecer."));
    }
  }

  async function handleFinalize(input: AnalysisInput): Promise<void> {
    try {
      await finalizeSolicitation.mutateAsync({ id: solicitationId, input });
      toast.success("Solicitação finalizada.");
    } catch (error) {
      toast.error(errorMessage(error, "Não foi possível finalizar a solicitação."));
    }
  }

  return (
    <div className="mx-auto flex max-w-3xl flex-col gap-6">
      <div className="flex flex-col gap-2">
        <div className="flex flex-wrap items-center justify-between gap-2">
          <h1 className="text-xl font-semibold tracking-tight">{detail.title || "Rascunho sem título"}</h1>
          <div className="flex items-center gap-2">
            <StatusBadge status={detail.status} />
            <PriorityBadge priority={detail.priority} />
          </div>
        </div>
        <p className="text-sm text-muted-foreground">Solicitado por {detail.requesterName}</p>
      </div>

      <PendingActorBanner status={detail.status} pendingActor={detail.pendingActor} />

      {detail.status === "recusada" && detail.rejectionReason && (
        <div className="rounded-lg border border-destructive/30 bg-destructive/10 px-4 py-3 text-sm text-destructive">
          Motivo da recusa: {detail.rejectionReason}
        </div>
      )}

      <SolicitationForm
        defaultValues={{
          title: detail.title,
          problemDescription: detail.problemDescription,
          proposedImprovement: detail.proposedImprovement,
          categoryId: detail.categoryId ?? "",
          location: detail.location,
        }}
        categories={categoriesQuery.data ?? []}
        readOnly={!detail.permissions.canEdit}
        onSaveDraft={handleSaveDraft}
        onSubmitForApproval={handleSubmitForApproval}
        isSavingDraft={updateDraft.isPending && !submitSolicitation.isPending}
        isSubmittingForApproval={submitSolicitation.isPending}
      />

      {detail.permissions.canApprove && (
        <ApproveRejectPanel
          onApprove={handleApprove}
          onReject={handleReject}
          isApproving={approveSolicitation.isPending}
          isRejecting={rejectSolicitation.isPending}
        />
      )}

      {detail.permissions.canAnalyze && (
        <AnalysisPanel
          defaultValues={{
            severity: detail.severity,
            urgency: detail.urgency,
            trend: detail.trend,
            analysisNotes: detail.analysisNotes,
          }}
          onSavePartial={handleSaveAnalysisPartial}
          onFinalize={handleFinalize}
          isSavingPartial={patchAnalysis.isPending}
          isFinalizing={finalizeSolicitation.isPending}
        />
      )}

      <Separator />

      <div className="flex flex-col gap-3">
        <h2 className="text-sm font-semibold">Histórico</h2>
        <HistoryTimeline history={detail.history} />
      </div>
    </div>
  );
}
