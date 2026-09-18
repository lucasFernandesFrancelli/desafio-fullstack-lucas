import {
  useMutation,
  useQuery,
  useQueryClient,
  type UseMutationResult,
  type UseQueryResult,
} from "@tanstack/react-query";

import { httpClient } from "@/api/httpClient";
import type {
  AnalysisInput,
  DraftInput,
  Solicitation,
  SolicitationDetail,
  SolicitationSummary,
  Status,
} from "@/types/api";

export interface SolicitationListFilter {
  status?: Status;
  categoryId?: string;
  q?: string;
}

export const solicitationKeys = {
  all: ["solicitations"] as const,
  list: (filter: SolicitationListFilter) => [...solicitationKeys.all, "list", filter] as const,
  detail: (id: string) => [...solicitationKeys.all, "detail", id] as const,
};

function buildListQuery(filter: SolicitationListFilter): string {
  const params = new URLSearchParams();
  if (filter.status) params.set("status", filter.status);
  if (filter.categoryId) params.set("categoryId", filter.categoryId);
  if (filter.q) params.set("q", filter.q);
  const query = params.toString();
  return query ? `?${query}` : "";
}

function fetchSolicitations(filter: SolicitationListFilter): Promise<SolicitationSummary[]> {
  return httpClient.get<SolicitationSummary[]>(`/solicitations${buildListQuery(filter)}`);
}

function fetchSolicitationDetail(id: string): Promise<SolicitationDetail> {
  return httpClient.get<SolicitationDetail>(`/solicitations/${id}`);
}

function createDraft(input: DraftInput): Promise<Solicitation> {
  return httpClient.post<Solicitation>("/solicitations", input);
}

function updateDraft(id: string, input: DraftInput): Promise<Solicitation> {
  return httpClient.patch<Solicitation>(`/solicitations/${id}`, input);
}

function submitSolicitation(id: string): Promise<Solicitation> {
  return httpClient.post<Solicitation>(`/solicitations/${id}/submit`);
}

function approveSolicitation(id: string): Promise<Solicitation> {
  return httpClient.post<Solicitation>(`/solicitations/${id}/approve`);
}

function rejectSolicitation(id: string, reason: string): Promise<Solicitation> {
  return httpClient.post<Solicitation>(`/solicitations/${id}/reject`, { reason });
}

function patchAnalysis(id: string, input: AnalysisInput): Promise<Solicitation> {
  return httpClient.patch<Solicitation>(`/solicitations/${id}/analysis`, input);
}

function finalizeSolicitation(id: string, input: AnalysisInput): Promise<Solicitation> {
  return httpClient.post<Solicitation>(`/solicitations/${id}/finalize`, input);
}

export function useSolicitations(filter: SolicitationListFilter): UseQueryResult<SolicitationSummary[], Error> {
  return useQuery({
    queryKey: solicitationKeys.list(filter),
    queryFn: () => fetchSolicitations(filter),
  });
}

export function useSolicitationDetail(id: string | undefined): UseQueryResult<SolicitationDetail, Error> {
  return useQuery({
    queryKey: solicitationKeys.detail(id ?? ""),
    queryFn: () => fetchSolicitationDetail(id ?? ""),
    enabled: Boolean(id),
  });
}

/** Invalida a lista (Lista + Kanban) e, quando houver id, o detalhe — mantém as duas visões sempre consistentes. */
function useInvalidateSolicitations() {
  const queryClient = useQueryClient();
  return (id?: string) => {
    void queryClient.invalidateQueries({ queryKey: solicitationKeys.all });
    if (id) {
      void queryClient.invalidateQueries({ queryKey: solicitationKeys.detail(id) });
    }
  };
}

export function useCreateDraft(): UseMutationResult<Solicitation, Error, DraftInput> {
  const invalidate = useInvalidateSolicitations();
  return useMutation({
    mutationFn: createDraft,
    onSuccess: () => invalidate(),
  });
}

export function useUpdateDraft(): UseMutationResult<Solicitation, Error, { id: string; input: DraftInput }> {
  const invalidate = useInvalidateSolicitations();
  return useMutation({
    mutationFn: ({ id, input }) => updateDraft(id, input),
    onSuccess: (_data, variables) => invalidate(variables.id),
  });
}

export function useSubmitSolicitation(): UseMutationResult<Solicitation, Error, string> {
  const invalidate = useInvalidateSolicitations();
  return useMutation({
    mutationFn: submitSolicitation,
    onSuccess: (_data, id) => invalidate(id),
  });
}

export function useApproveSolicitation(): UseMutationResult<Solicitation, Error, string> {
  const invalidate = useInvalidateSolicitations();
  return useMutation({
    mutationFn: approveSolicitation,
    onSuccess: (_data, id) => invalidate(id),
  });
}

export function useRejectSolicitation(): UseMutationResult<Solicitation, Error, { id: string; reason: string }> {
  const invalidate = useInvalidateSolicitations();
  return useMutation({
    mutationFn: ({ id, reason }) => rejectSolicitation(id, reason),
    onSuccess: (_data, variables) => invalidate(variables.id),
  });
}

export function usePatchAnalysis(): UseMutationResult<Solicitation, Error, { id: string; input: AnalysisInput }> {
  const invalidate = useInvalidateSolicitations();
  return useMutation({
    mutationFn: ({ id, input }) => patchAnalysis(id, input),
    onSuccess: (_data, variables) => invalidate(variables.id),
  });
}

export function useFinalizeSolicitation(): UseMutationResult<Solicitation, Error, { id: string; input: AnalysisInput }> {
  const invalidate = useInvalidateSolicitations();
  return useMutation({
    mutationFn: ({ id, input }) => finalizeSolicitation(id, input),
    onSuccess: (_data, variables) => invalidate(variables.id),
  });
}
