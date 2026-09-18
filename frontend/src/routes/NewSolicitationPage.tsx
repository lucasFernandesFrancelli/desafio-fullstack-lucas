import type { ReactElement } from "react";
import { useNavigate } from "react-router";
import { toast } from "sonner";

import { useCategories } from "@/api/categories";
import { useCreateDraft, useSubmitSolicitation } from "@/api/solicitations";
import { ApiError } from "@/api/httpClient";
import { Skeleton } from "@/components/ui/skeleton";
import { SolicitationForm } from "@/features/solicitations/SolicitationForm";
import { emptyDraftFormValues, type DraftFormValues } from "@/schemas/solicitationSchemas";
import type { DraftInput } from "@/types/api";

export function NewSolicitationPage(): ReactElement {
  const navigate = useNavigate();
  const categoriesQuery = useCategories(true);
  const createDraft = useCreateDraft();
  const submitSolicitation = useSubmitSolicitation();

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
      const created = await createDraft.mutateAsync(toDraftInput(values));
      toast.success("Rascunho salvo.");
      navigate(`/solicitations/${created.id}`, { replace: true });
    } catch (error) {
      toast.error(error instanceof ApiError ? error.message : "Não foi possível salvar o rascunho.");
    }
  }

  async function handleSubmitForApproval(values: DraftFormValues): Promise<void> {
    try {
      const created = await createDraft.mutateAsync(toDraftInput(values));
      await submitSolicitation.mutateAsync(created.id);
      toast.success("Solicitação enviada para aprovação.");
      navigate(`/solicitations/${created.id}`, { replace: true });
    } catch (error) {
      toast.error(error instanceof ApiError ? error.message : "Não foi possível enviar a solicitação.");
    }
  }

  return (
    <div className="mx-auto flex max-w-2xl flex-col gap-6">
      <div className="flex flex-col gap-1">
        <h1 className="text-xl font-semibold tracking-tight">Nova solicitação</h1>
        <p className="text-sm text-muted-foreground">
          Descreva o problema e a melhoria proposta. Você pode salvar como rascunho e continuar depois.
        </p>
      </div>

      {categoriesQuery.isPending && (
        <div className="flex flex-col gap-5">
          <Skeleton className="h-9 w-full" />
          <Skeleton className="h-24 w-full" />
          <Skeleton className="h-24 w-full" />
          <div className="grid gap-5 sm:grid-cols-2">
            <Skeleton className="h-9 w-full" />
            <Skeleton className="h-9 w-full" />
          </div>
        </div>
      )}

      {categoriesQuery.isError && (
        <p className="text-sm text-destructive">Não foi possível carregar as categorias. Recarregue a página.</p>
      )}

      {categoriesQuery.data && (
        <SolicitationForm
          defaultValues={emptyDraftFormValues}
          categories={categoriesQuery.data}
          onSaveDraft={handleSaveDraft}
          onSubmitForApproval={handleSubmitForApproval}
          isSavingDraft={createDraft.isPending && !submitSolicitation.isPending}
          isSubmittingForApproval={submitSolicitation.isPending}
        />
      )}
    </div>
  );
}
