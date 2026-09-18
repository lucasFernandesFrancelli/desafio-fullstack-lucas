import type {
  Category,
  DashboardSummary,
  Solicitation,
  SolicitationDetail,
  SolicitationSummary,
  User,
  UserProfile,
} from "@/types/api";

export const collaboratorUser: User = {
  id: "user-collaborator",
  name: "Ana Beatriz Costa",
  email: "ana.costa@ekaizen.example",
  role: "colaborador",
  createdAt: "2026-01-01T00:00:00Z",
};

export const approverStep1User: User = {
  id: "user-approver-1",
  name: "Ricardo Nogueira",
  email: "ricardo.nogueira@ekaizen.example",
  role: "colaborador",
  createdAt: "2026-01-01T00:00:00Z",
};

export const approverStep2User: User = {
  id: "user-approver-2",
  name: "Fernanda Albuquerque",
  email: "fernanda.albuquerque@ekaizen.example",
  role: "colaborador",
  createdAt: "2026-01-01T00:00:00Z",
};

export const analystUser: User = {
  id: "user-analyst",
  name: "Renata Souza",
  email: "renata.souza@ekaizen.example",
  role: "analista",
  createdAt: "2026-01-01T00:00:00Z",
};

export const managerUser: User = {
  id: "user-manager",
  name: "Sérgio Andrade",
  email: "sergio.andrade@ekaizen.example",
  role: "gestor",
  createdAt: "2026-01-01T00:00:00Z",
};

export const profiles: UserProfile[] = [
  { ...collaboratorUser, approverFor: [] },
  {
    ...approverStep1User,
    approverFor: [{ categoryId: "cat-seguranca", categoryName: "Segurança do Trabalho", order: 1 }],
  },
  { ...analystUser, approverFor: [] },
  { ...managerUser, approverFor: [] },
];

export const securityCategory: Category = {
  id: "cat-seguranca",
  name: "Segurança do Trabalho",
  description: "Riscos e melhorias de segurança.",
  approvers: [
    { userId: "user-approver-1", userName: "Ricardo Nogueira", order: 1 },
    { userId: "user-approver-2", userName: "Fernanda Albuquerque", order: 2 },
  ],
};

export const categories: Category[] = [
  securityCategory,
  {
    id: "cat-qualidade",
    name: "Qualidade",
    description: "Não conformidades de processo.",
    approvers: [
      { userId: "user-approver-3", userName: "Camila Duarte", order: 1 },
      { userId: "user-approver-4", userName: "Marcos Vieira", order: 2 },
    ],
  },
];

export const draftSolicitation: SolicitationDetail = {
  id: "sol-draft",
  title: "Piso escorregadio no refeitório",
  problemDescription: "",
  proposedImprovement: "",
  categoryId: null,
  location: "",
  requesterId: collaboratorUser.id,
  requesterName: collaboratorUser.name,
  status: "rascunho",
  currentApprovalStep: null,
  severity: null,
  urgency: null,
  trend: null,
  priority: null,
  analysisNotes: "",
  rejectionReason: "",
  lastTransitionAt: "2026-01-02T10:00:00Z",
  createdAt: "2026-01-02T10:00:00Z",
  updatedAt: "2026-01-02T10:00:00Z",
  pendingActor: { kind: "user", id: collaboratorUser.id, name: collaboratorUser.name },
  permissions: {
    canEdit: true,
    canSubmit: true,
    canApprove: false,
    canReject: false,
    canAnalyze: false,
    canFinalize: false,
  },
  history: [
    {
      id: "hist-1",
      solicitationId: "sol-draft",
      fromStatus: null,
      toStatus: "rascunho",
      action: "criada",
      actorId: collaboratorUser.id,
      actorName: collaboratorUser.name,
      createdAt: "2026-01-02T10:00:00Z",
    },
  ],
};

/** Rascunho com todos os campos já preenchidos — usado para testar o envio
 * para aprovação sem precisar interagir com o Select de categoria (o Radix
 * Select não abre de forma confiável em jsdom; essa interação real é
 * coberta pelos testes e2e). */
export const completeDraftSolicitation: SolicitationDetail = {
  ...draftSolicitation,
  id: "sol-draft-completo",
  problemDescription: "Risco de queda no refeitório",
  proposedImprovement: "Instalar piso antiderrapante",
  categoryId: "cat-seguranca",
  categoryName: "Segurança do Trabalho",
  location: "Refeitório",
};

export const pendingApprovalSolicitation: SolicitationDetail = {
  id: "sol-em-aprovacao",
  title: "Retrabalho recorrente na linha de embalagem",
  problemDescription: "Itens saem com etiqueta trocada.",
  proposedImprovement: "Checklist de dupla checagem.",
  categoryId: "cat-seguranca",
  categoryName: "Segurança do Trabalho",
  location: "Linha 3",
  requesterId: collaboratorUser.id,
  requesterName: collaboratorUser.name,
  status: "em_aprovacao",
  currentApprovalStep: 1,
  severity: null,
  urgency: null,
  trend: null,
  priority: null,
  analysisNotes: "",
  rejectionReason: "",
  lastTransitionAt: "2026-01-02T11:00:00Z",
  createdAt: "2026-01-02T10:00:00Z",
  updatedAt: "2026-01-02T11:00:00Z",
  pendingActor: { kind: "user", id: approverStep1User.id, name: approverStep1User.name },
  permissions: {
    canEdit: false,
    canSubmit: false,
    canApprove: true,
    canReject: true,
    canAnalyze: false,
    canFinalize: false,
  },
  history: [
    {
      id: "hist-2",
      solicitationId: "sol-em-aprovacao",
      fromStatus: null,
      toStatus: "rascunho",
      action: "criada",
      actorId: collaboratorUser.id,
      actorName: collaboratorUser.name,
      createdAt: "2026-01-02T10:00:00Z",
    },
    {
      id: "hist-3",
      solicitationId: "sol-em-aprovacao",
      fromStatus: "rascunho",
      toStatus: "em_aprovacao",
      action: "enviada",
      actorId: collaboratorUser.id,
      actorName: collaboratorUser.name,
      createdAt: "2026-01-02T11:00:00Z",
    },
  ],
};

export const inAnalysisSolicitation: SolicitationDetail = {
  id: "sol-em-analise",
  title: "Tempo de setup elevado na troca de molde",
  problemDescription: "Troca de molde leva 90 minutos.",
  proposedImprovement: "Aplicar técnica SMED.",
  categoryId: "cat-qualidade",
  categoryName: "Qualidade",
  location: "Prensa 4",
  requesterId: collaboratorUser.id,
  requesterName: collaboratorUser.name,
  status: "em_analise",
  currentApprovalStep: null,
  severity: null,
  urgency: null,
  trend: null,
  priority: null,
  analysisNotes: "",
  rejectionReason: "",
  lastTransitionAt: "2026-01-02T12:00:00Z",
  createdAt: "2026-01-02T10:00:00Z",
  updatedAt: "2026-01-02T12:00:00Z",
  pendingActor: { kind: "role", name: "Qualquer analista" },
  permissions: {
    canEdit: false,
    canSubmit: false,
    canApprove: false,
    canReject: false,
    canAnalyze: true,
    canFinalize: true,
  },
  history: [],
};

export function toSummary(detail: Solicitation): SolicitationSummary {
  return {
    id: detail.id,
    title: detail.title,
    categoryName: detail.categoryName,
    status: detail.status,
    currentApprovalStep: detail.currentApprovalStep,
    requesterName: detail.requesterName ?? "",
    priority: detail.priority,
    pendingActorName: "",
    createdAt: detail.createdAt,
    updatedAt: detail.updatedAt,
  };
}

export const solicitationSummaries: SolicitationSummary[] = [
  toSummary(draftSolicitation),
  toSummary(pendingApprovalSolicitation),
  toSummary(inAnalysisSolicitation),
];

export const dashboardSummary: DashboardSummary = {
  countsByStatus: {
    rascunho: 1,
    em_aprovacao: 1,
    em_analise: 1,
    finalizada: 0,
    recusada: 0,
  },
  oldestPending: [
    {
      id: draftSolicitation.id,
      title: draftSolicitation.title,
      status: draftSolicitation.status,
      categoryName: "",
      pendingActorName: collaboratorUser.name,
      daysSinceLastMove: 3,
    },
  ],
  awaitingByUser: [{ userId: approverStep1User.id, userName: approverStep1User.name, count: 1 }],
};

export const allUsers: User[] = [collaboratorUser, approverStep1User, approverStep2User, analystUser, managerUser];
