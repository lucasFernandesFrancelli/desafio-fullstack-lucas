import { z } from "zod";

export const rejectSchema = z.object({
  reason: z.string().min(1, "Motivo da recusa é obrigatório"),
});

export type RejectFormValues = z.infer<typeof rejectSchema>;
