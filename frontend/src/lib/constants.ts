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

/** Cor de destaque por status (borda de card, topo de coluna do Kanban etc.). */
export const STATUS_ACCENT_BORDER_CLASSES: Record<Status, string> = {
  rascunho: "border-l-muted-foreground/40",
  em_aprovacao: "border-l-amber-400 dark:border-l-amber-600",
  em_analise: "border-l-sky-400 dark:border-l-sky-600",
  finalizada: "border-l-emerald-500 dark:border-l-emerald-600",
  recusada: "border-l-red-500 dark:border-l-red-600",
};

export const STATUS_ACCENT_TOP_BORDER_CLASSES: Record<Status, string> = {
  rascunho: "border-t-muted-foreground/40",
  em_aprovacao: "border-t-amber-400 dark:border-t-amber-600",
  em_analise: "border-t-sky-400 dark:border-t-sky-600",
  finalizada: "border-t-emerald-500 dark:border-t-emerald-600",
  recusada: "border-t-red-500 dark:border-t-red-600",
};

/** Cor do número em destaque, usada nos cards de contagem do dashboard. */
export const STATUS_TEXT_CLASSES: Record<Status, string> = {
  rascunho: "text-foreground",
  em_aprovacao: "text-amber-600 dark:text-amber-400",
  em_analise: "text-sky-600 dark:text-sky-400",
  finalizada: "text-emerald-600 dark:text-emerald-400",
  recusada: "text-red-600 dark:text-red-400",
};

export const ROLE_LABELS: Record<Role, string> = {
  colaborador: "Colaborador",
  analista: "Analista",
  gestor: "Gestor",
};

/** Classes do avatar (iniciais) por papel, para diferenciar quem é quem no seletor de perfil. */
export const ROLE_AVATAR_CLASSES: Record<Role, string> = {
  colaborador: "bg-gradient-to-br from-slate-500 to-slate-700 text-white",
  analista: "bg-gradient-to-br from-sky-500 to-blue-700 text-white",
  gestor: "bg-gradient-to-br from-[#d02110] to-[#f2711a] text-white",
};

export const ROLE_BADGE_CLASSES: Record<Role, string> = {
  colaborador: "bg-slate-100 text-slate-700 border-slate-300 dark:bg-slate-800/40 dark:text-slate-300 dark:border-slate-700",
  analista: "bg-sky-100 text-sky-900 border-sky-300 dark:bg-sky-900/30 dark:text-sky-300 dark:border-sky-800",
  gestor: "bg-red-100 text-red-900 border-red-300 dark:bg-red-900/30 dark:text-red-300 dark:border-red-800",
};

/** Cor do marcador de linha do tempo por ação, para reforçar visualmente o tipo de decisão. */
export const HISTORY_ACTION_DOT_CLASSES: Record<HistoryAction, string> = {
  criada: "bg-muted-foreground/50",
  enviada: "bg-primary",
  aprovada_etapa1: "bg-amber-500",
  aprovada_etapa2: "bg-amber-500",
  recusada: "bg-red-500",
  finalizada: "bg-emerald-500",
};

export const HISTORY_ACTION_LABELS: Record<HistoryAction, string> = {
  criada: "Criou o rascunho",
  enviada: "Enviou para aprovação",
  aprovada_etapa1: "Aprovou (1ª etapa)",
  aprovada_etapa2: "Aprovou (2ª etapa)",
  recusada: "Recusou a solicitação",
  finalizada: "Finalizou com parecer",
};
