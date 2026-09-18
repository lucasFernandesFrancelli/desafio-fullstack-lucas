// Tipos espelhando exatamente os JSON tags da API Go (backend/internal/models).
// Mantidos manualmente em vez de gerados: o contrato é pequeno e estável.

export type Role = "colaborador" | "analista" | "gestor";

export type Status =
  | "rascunho"
  | "em_aprovacao"
  | "em_analise"
  | "finalizada"
  | "recusada";

export type HistoryAction =
  | "criada"
  | "enviada"
  | "aprovada_etapa1"
  | "aprovada_etapa2"
  | "recusada"
  | "finalizada";

export interface User {
  id: string;
  name: string;
  email: string;
  role: Role;
  createdAt: string;
}

export interface ApproverFor {
  categoryId: string;
  categoryName: string;
  order: 1 | 2;
}

export interface UserProfile extends User {
  approverFor: ApproverFor[];
}

export interface CategoryApprover {
  userId: string;
  userName: string;
  order: 1 | 2;
}

export interface Category {
  id: string;
  name: string;
  description: string;
  approvers: CategoryApprover[];
}

export interface Solicitation {
  id: string;
  title: string;
  problemDescription: string;
  proposedImprovement: string;
  categoryId: string | null;
  categoryName?: string;
  location: string;
  requesterId: string;
  requesterName?: string;
  status: Status;
  currentApprovalStep: 1 | 2 | null;
  severity: number | null;
  urgency: number | null;
  trend: number | null;
  priority: number | null;
  analysisNotes: string;
  rejectionReason: string;
  lastTransitionAt: string;
  createdAt: string;
  updatedAt: string;
}

export interface SolicitationSummary {
  id: string;
  title: string;
  categoryName?: string;
  status: Status;
  currentApprovalStep: 1 | 2 | null;
  requesterName: string;
  priority: number | null;
  pendingActorName: string;
  createdAt: string;
  updatedAt: string;
}

export type PendingActorKind = "user" | "role" | "none";

export interface PendingActor {
  kind: PendingActorKind;
  id?: string;
  name: string;
}

export interface Permissions {
  canEdit: boolean;
  canSubmit: boolean;
  canApprove: boolean;
  canReject: boolean;
  canAnalyze: boolean;
  canFinalize: boolean;
}

export interface HistoryEntry {
  id: string;
  solicitationId: string;
  fromStatus: Status | null;
  toStatus: Status;
  action: HistoryAction;
  actorId: string;
  actorName: string;
  approvalStep?: 1 | 2;
  comment?: string;
  severity?: number;
  urgency?: number;
  trend?: number;
  createdAt: string;
}

export interface SolicitationDetail extends Solicitation {
  pendingActor: PendingActor;
  permissions: Permissions;
  history: HistoryEntry[];
}

export interface OldestPendingItem {
  id: string;
  title: string;
  status: Status;
  categoryName: string;
  pendingActorName: string;
  daysSinceLastMove: number;
}

export interface UserPendingCount {
  userId: string;
  userName: string;
  count: number;
}

export type CountsByStatus = Partial<Record<Status, number>>;

export interface DashboardSummary {
  countsByStatus: CountsByStatus;
  oldestPending: OldestPendingItem[];
  awaitingByUser: UserPendingCount[];
}

// ---- Request DTOs ----

export interface DraftInput {
  title?: string;
  problemDescription?: string;
  proposedImprovement?: string;
  categoryId?: string;
  location?: string;
}

export interface AnalysisInput {
  severity?: number;
  urgency?: number;
  trend?: number;
  analysisNotes?: string;
}

export interface RejectInput {
  reason: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface CreateUserInput {
  name: string;
  email: string;
  role: Role;
}

export interface CreateCategoryInput {
  name: string;
  description?: string;
  firstApproverId: string;
  secondApproverId: string;
}

export interface UpdateCategoryInput {
  name?: string;
  description?: string;
}

export interface SetApproversInput {
  firstApproverId: string;
  secondApproverId: string;
}

// ---- Erros ----

export interface ApiErrorBody {
  error: {
    code: string;
    message: string;
    fields?: Record<string, string>;
  };
}
