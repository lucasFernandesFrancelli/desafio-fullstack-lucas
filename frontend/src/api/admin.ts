import {
  useMutation,
  useQuery,
  useQueryClient,
  type UseMutationResult,
  type UseQueryResult,
} from "@tanstack/react-query";

import { httpClient } from "@/api/httpClient";
import type {
  Category,
  CreateCategoryInput,
  CreateUserInput,
  SetApproversInput,
  UpdateCategoryInput,
  User,
} from "@/types/api";

function fetchUsers(): Promise<User[]> {
  return httpClient.get<User[]>("/users");
}

function createUser(input: CreateUserInput): Promise<User> {
  return httpClient.post<User>("/users", input);
}

function createCategory(input: CreateCategoryInput): Promise<Category> {
  return httpClient.post<Category>("/categories", input);
}

function updateCategory(id: string, input: UpdateCategoryInput): Promise<Category> {
  return httpClient.patch<Category>(`/categories/${id}`, input);
}

function setCategoryApprovers(id: string, input: SetApproversInput): Promise<Category> {
  return httpClient.patch<Category>(`/categories/${id}/approvers`, input);
}

export function useUsers(enabled: boolean): UseQueryResult<User[], Error> {
  return useQuery({
    queryKey: ["admin", "users"],
    queryFn: fetchUsers,
    enabled,
  });
}

export function useCreateUser(): UseMutationResult<User, Error, CreateUserInput> {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: createUser,
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["admin", "users"] });
      void queryClient.invalidateQueries({ queryKey: ["auth", "profiles"] });
    },
  });
}

function useInvalidateCategories() {
  const queryClient = useQueryClient();
  return () => {
    void queryClient.invalidateQueries({ queryKey: ["categories"] });
    void queryClient.invalidateQueries({ queryKey: ["auth", "profiles"] });
  };
}

export function useCreateCategory(): UseMutationResult<Category, Error, CreateCategoryInput> {
  const invalidate = useInvalidateCategories();
  return useMutation({
    mutationFn: createCategory,
    onSuccess: () => invalidate(),
  });
}

export function useUpdateCategory(): UseMutationResult<Category, Error, { id: string; input: UpdateCategoryInput }> {
  const invalidate = useInvalidateCategories();
  return useMutation({
    mutationFn: ({ id, input }) => updateCategory(id, input),
    onSuccess: () => invalidate(),
  });
}

export function useSetCategoryApprovers(): UseMutationResult<
  Category,
  Error,
  { id: string; input: SetApproversInput }
> {
  const invalidate = useInvalidateCategories();
  return useMutation({
    mutationFn: ({ id, input }) => setCategoryApprovers(id, input),
    onSuccess: () => invalidate(),
  });
}
