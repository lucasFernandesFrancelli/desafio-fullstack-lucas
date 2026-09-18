import { useMutation, useQuery, type UseMutationResult, type UseQueryResult } from "@tanstack/react-query";

import { httpClient } from "@/api/httpClient";
import type { LoginResponse, User, UserProfile } from "@/types/api";

function fetchProfiles(): Promise<UserProfile[]> {
  return httpClient.get<UserProfile[]>("/auth/profiles", { auth: false });
}

function login(userId: string): Promise<LoginResponse> {
  return httpClient.post<LoginResponse>("/auth/login", { userId }, { auth: false });
}

function fetchMe(): Promise<User> {
  return httpClient.get<User>("/auth/me");
}

export function useProfiles(): UseQueryResult<UserProfile[], Error> {
  return useQuery({
    queryKey: ["auth", "profiles"],
    queryFn: fetchProfiles,
  });
}

export function useLoginMutation(): UseMutationResult<LoginResponse, Error, string> {
  return useMutation({
    mutationFn: login,
  });
}

export function useMe(enabled: boolean): UseQueryResult<User, Error> {
  return useQuery({
    queryKey: ["auth", "me"],
    queryFn: fetchMe,
    enabled,
    retry: false,
  });
}
