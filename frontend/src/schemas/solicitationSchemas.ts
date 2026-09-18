import { z } from "zod";

/** Rascunho: só o título é exigido (regra do backend), o resto pode ficar incompleto. */
export const draftSchema = z.object({
  title: z.string().min(1, "Título é obrigatório"),
  problemDescription: z.string(),
  proposedImprovement: z.string(),
  categoryId: z.string(),
  location: z.string(),
});

export type DraftFormValues = z.infer<typeof draftSchema>;

export const emptyDraftFormValues: DraftFormValues = {
  title: "",
  problemDescription: "",
  proposedImprovement: "",
  categoryId: "",
  location: "",
};

/** Envio: todos os campos viram obrigatórios — mesma regra aplicada no backend. */
export const submitSchema = z.object({
  title: z.string().min(1, "Título é obrigatório"),
  problemDescription: z.string().min(1, "Descrição do problema é obrigatória"),
  proposedImprovement: z.string().min(1, "Melhoria proposta é obrigatória"),
  categoryId: z.string().min(1, "Categoria é obrigatória"),
  location: z.string().min(1, "Local é obrigatório"),
});

export function validateForSubmit(values: DraftFormValues): Partial<Record<keyof DraftFormValues, string>> {
  const result = submitSchema.safeParse(values);
  if (result.success) {
    return {};
  }
  const fieldErrors: Partial<Record<keyof DraftFormValues, string>> = {};
  for (const issue of result.error.issues) {
    const key = issue.path[0];
    if (typeof key === "string") {
      fieldErrors[key as keyof DraftFormValues] = issue.message;
    }
  }
  return fieldErrors;
}
