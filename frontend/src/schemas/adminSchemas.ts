import { z } from "zod";

export const createUserSchema = z.object({
  name: z.string().min(1, "Nome é obrigatório"),
  email: z.string().min(1, "E-mail é obrigatório").email("E-mail inválido"),
  role: z.enum(["colaborador", "analista", "gestor"], { message: "Selecione um papel" }),
});

export type CreateUserFormValues = z.infer<typeof createUserSchema>;

export const createCategorySchema = z
  .object({
    name: z.string().min(1, "Nome é obrigatório"),
    description: z.string(),
    firstApproverId: z.string().min(1, "Selecione o 1º aprovador"),
    secondApproverId: z.string().min(1, "Selecione o 2º aprovador"),
  })
  .refine((values) => values.firstApproverId !== values.secondApproverId, {
    message: "Os dois aprovadores devem ser pessoas distintas",
    path: ["secondApproverId"],
  });

export type CreateCategoryFormValues = z.infer<typeof createCategorySchema>;

export const setApproversSchema = z
  .object({
    firstApproverId: z.string().min(1, "Selecione o 1º aprovador"),
    secondApproverId: z.string().min(1, "Selecione o 2º aprovador"),
  })
  .refine((values) => values.firstApproverId !== values.secondApproverId, {
    message: "Os dois aprovadores devem ser pessoas distintas",
    path: ["secondApproverId"],
  });

export type SetApproversFormValues = z.infer<typeof setApproversSchema>;
