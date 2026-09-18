import { z } from "zod";

const scoreField = z.coerce.number().int().min(1, "Obrigatório (1 a 5)").max(5, "Deve ser entre 1 e 5");

/** Usado ao clicar em "Finalizar": parecer e as 3 notas passam a ser obrigatórios. */
export const analysisFinalizeSchema = z.object({
  severity: scoreField,
  urgency: scoreField,
  trend: scoreField,
  analysisNotes: z.string().min(1, "Parecer é obrigatório"),
});

export type AnalysisFinalizeFormValues = z.infer<typeof analysisFinalizeSchema>;
