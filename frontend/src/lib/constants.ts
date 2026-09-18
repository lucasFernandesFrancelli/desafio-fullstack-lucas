import type { HistoryAction, Role, Status } from "@/types/api";

export const STATUS_ORDER: readonly Status[] = [
  "rascunho",
  "em_aprovacao",
  "em_analise",
  "finalizada",
  "recusada",
];

export const STATUS_LABELS: Record<Status, string> = {
  rascunho: "Rascunho",
  em_aprovacao: "Em aprovação",
  em_analise: "Em análise",
  finalizada: "Finalizada",
  recusada: "Recusada",
};

/** Classes Tailwind por status, usadas no StatusBadge e nas colunas do Kanban. */
export const STATUS_BADGE_CLASSES: Record<Status, string> = {
  rascunho: "bg-muted text-muted-foreground border-border",
  em_aprovacao: "bg-amber-100 text-amber-900 border-amber-300 dark:bg-amber-900/30 dark:text-amber-300 dark:border-amber-800",
  em_analise: "bg-sky-100 text-sky-900 border-sky-300 dark:bg-sky-900/30 dark:text-sky-300 dark:border-sky-800",
  finalizada: "bg-emerald-100 text-emerald-900 border-emerald-300 dark:bg-emerald-900/30 dark:text-emerald-300 dark:border-emerald-800",
  recusada: "bg-red-100 text-red-900 border-red-300 dark:bg-red-900/30 dark:text-red-300 dark:border-red-800",
};

export const ROLE_LABELS: Record<Role, string> = {
  colaborador: "Colaborador",
  analista: "Analista",
  gestor: "Gestor",
};

export const HISTORY_ACTION_LABELS: Record<HistoryAction, string> = {
  criada: "Criou o rascunho",
  enviada: "Enviou para aprovação",
  aprovada_etapa1: "Aprovou (1ª etapa)",
  aprovada_etapa2: "Aprovou (2ª etapa)",
  recusada: "Recusou a solicitação",
  finalizada: "Finalizou com parecer",
};
