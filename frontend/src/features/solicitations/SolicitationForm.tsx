import type { ReactElement } from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import { draftSchema, validateForSubmit, type DraftFormValues } from "@/schemas/solicitationSchemas";
import type { Category } from "@/types/api";

interface SolicitationFormProps {
  defaultValues: DraftFormValues;
  categories: Category[];
  readOnly?: boolean;
  onSaveDraft: (values: DraftFormValues) => Promise<void>;
  onSubmitForApproval: (values: DraftFormValues) => Promise<void>;
  isSavingDraft: boolean;
  isSubmittingForApproval: boolean;
}

export function SolicitationForm({
  defaultValues,
  categories,
  readOnly = false,
  onSaveDraft,
  onSubmitForApproval,
  isSavingDraft,
  isSubmittingForApproval,
}: SolicitationFormProps): ReactElement {
  const form = useForm<DraftFormValues>({
    resolver: zodResolver(draftSchema),
    defaultValues,
  });

  async function handleSaveDraft(values: DraftFormValues): Promise<void> {
    await onSaveDraft(values);
  }

  async function handleSubmitForApproval(): Promise<void> {
    const values = form.getValues();
    const fieldErrors = validateForSubmit(values);
    const entries = Object.entries(fieldErrors) as Array<[keyof DraftFormValues, string]>;

    if (entries.length > 0) {
      for (const [field, message] of entries) {
        form.setError(field, { message });
      }
      return;
    }

    await onSubmitForApproval(values);
  }

  const busy = isSavingDraft || isSubmittingForApproval;

  return (
    <form onSubmit={form.handleSubmit(handleSaveDraft)} className="flex flex-col gap-5" noValidate>
      <div className="flex flex-col gap-1.5">
        <Label htmlFor="title">Título</Label>
        <Input id="title" disabled={readOnly} {...form.register("title")} />
        {form.formState.errors.title && (
          <p className="text-xs text-destructive">{form.formState.errors.title.message}</p>
        )}
      </div>

      <div className="flex flex-col gap-1.5">
        <Label htmlFor="problemDescription">Descrição do problema</Label>
        <Textarea id="problemDescription" rows={4} disabled={readOnly} {...form.register("problemDescription")} />
        {form.formState.errors.problemDescription && (
          <p className="text-xs text-destructive">{form.formState.errors.problemDescription.message}</p>
        )}
      </div>

      <div className="flex flex-col gap-1.5">
        <Label htmlFor="proposedImprovement">Melhoria proposta</Label>
        <Textarea id="proposedImprovement" rows={4} disabled={readOnly} {...form.register("proposedImprovement")} />
        {form.formState.errors.proposedImprovement && (
          <p className="text-xs text-destructive">{form.formState.errors.proposedImprovement.message}</p>
        )}
      </div>

      <div className="grid gap-5 sm:grid-cols-2">
        <div className="flex flex-col gap-1.5">
          <Label htmlFor="categoryId">Categoria</Label>
          <Controller
            control={form.control}
            name="categoryId"
            render={({ field }) => (
              <Select disabled={readOnly} value={field.value} onValueChange={field.onChange}>
                <SelectTrigger id="categoryId" className="w-full">
                  <SelectValue placeholder="Selecione uma categoria" />
                </SelectTrigger>
                <SelectContent>
                  {categories.map((category) => (
                    <SelectItem key={category.id} value={category.id}>
                      {category.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            )}
          />
          {form.formState.errors.categoryId && (
            <p className="text-xs text-destructive">{form.formState.errors.categoryId.message}</p>
          )}
        </div>

        <div className="flex flex-col gap-1.5">
          <Label htmlFor="location">Local</Label>
          <Input id="location" disabled={readOnly} {...form.register("location")} />
          {form.formState.errors.location && (
            <p className="text-xs text-destructive">{form.formState.errors.location.message}</p>
          )}
        </div>
      </div>

      {!readOnly && (
        <div className="flex justify-end gap-2 pt-2">
          <Button type="submit" variant="outline" disabled={busy}>
            {isSavingDraft ? "Salvando…" : "Salvar rascunho"}
          </Button>
          <Button type="button" onClick={() => void handleSubmitForApproval()} disabled={busy}>
            {isSubmittingForApproval ? "Enviando…" : "Enviar para aprovação"}
          </Button>
        </div>
      )}
    </form>
  );
}
