import type { ApiErrorBody } from "@/types/api";

const AUTH_TOKEN_STORAGE_KEY = "ekaizen.token";

/** Erro tipado lançado para qualquer resposta HTTP não-2xx da API. */
export class ApiError extends Error {
  readonly code: string;
  readonly status: number;
  readonly fields?: Record<string, string>;

  constructor(status: number, body: ApiErrorBody) {
    super(body.error.message);
    this.name = "ApiError";
    this.code = body.error.code;
    this.status = status;
    this.fields = body.error.fields;
  }
}

function getBaseUrl(): string {
  const base = import.meta.env.VITE_API_BASE_URL;
  if (!base) {
    throw new Error("VITE_API_BASE_URL não configurada");
  }
  return base.replace(/\/$/, "");
}

export function getStoredToken(): string | null {
  return localStorage.getItem(AUTH_TOKEN_STORAGE_KEY);
}

export function setStoredToken(token: string | null): void {
  if (token) {
    localStorage.setItem(AUTH_TOKEN_STORAGE_KEY, token);
  } else {
    localStorage.removeItem(AUTH_TOKEN_STORAGE_KEY);
  }
}

type HttpMethod = "GET" | "POST" | "PATCH";

interface RequestOptions {
  method?: HttpMethod;
  body?: unknown;
  /** Envia o Bearer token, se houver um armazenado. Padrão: true. */
  auth?: boolean;
}

async function parseBody(response: Response): Promise<unknown> {
  const text = await response.text();
  if (!text) {
    return undefined;
  }
  return JSON.parse(text) as unknown;
}

function isApiErrorBody(value: unknown): value is ApiErrorBody {
  return (
    typeof value === "object" &&
    value !== null &&
    "error" in value &&
    typeof (value as { error: unknown }).error === "object"
  );
}

async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const { method = "GET", body, auth = true } = options;

  const headers: Record<string, string> = {};
  if (body !== undefined) {
    headers["Content-Type"] = "application/json";
  }
  if (auth) {
    const token = getStoredToken();
    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }
  }

  const response = await fetch(`${getBaseUrl()}${path}`, {
    method,
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  });

  const data = await parseBody(response);

  if (!response.ok) {
    const errorBody: ApiErrorBody = isApiErrorBody(data)
      ? data
      : { error: { code: "unknown_error", message: "Erro inesperado na comunicação com a API" } };
    throw new ApiError(response.status, errorBody);
  }

  return data as T;
}

export const httpClient = {
  get: <T>(path: string, options?: Omit<RequestOptions, "method" | "body">): Promise<T> =>
    request<T>(path, { ...options, method: "GET" }),
  post: <T>(path: string, body?: unknown, options?: Omit<RequestOptions, "method" | "body">): Promise<T> =>
    request<T>(path, { ...options, method: "POST", body }),
  patch: <T>(path: string, body?: unknown, options?: Omit<RequestOptions, "method" | "body">): Promise<T> =>
    request<T>(path, { ...options, method: "PATCH", body }),
};
