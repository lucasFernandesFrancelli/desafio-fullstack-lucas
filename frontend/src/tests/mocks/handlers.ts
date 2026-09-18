import { http, HttpResponse } from "msw";

import {
  categories,
  collaboratorUser,
  dashboardSummary,
  draftSolicitation,
  profiles,
  solicitationSummaries,
} from "@/tests/fixtures";
import type { ApiErrorBody, LoginResponse } from "@/types/api";

const BASE = "http://localhost:8080/api/v1";

function errorBody(code: string, message: string, fields?: Record<string, string>): ApiErrorBody {
  return { error: { code, message, fields } };
}

export const handlers = [
  http.get(`${BASE}/auth/profiles`, () => HttpResponse.json(profiles)),

  http.post(`${BASE}/auth/login`, () => {
    const response: LoginResponse = { token: "fake-jwt-token", user: collaboratorUser };
    return HttpResponse.json(response);
  }),

  http.get(`${BASE}/auth/me`, () => HttpResponse.json(collaboratorUser)),

  http.get(`${BASE}/categories`, () => HttpResponse.json(categories)),

  http.get(`${BASE}/solicitations`, () => HttpResponse.json(solicitationSummaries)),

  http.get(`${BASE}/solicitations/:id`, () => HttpResponse.json(draftSolicitation)),

  http.post(`${BASE}/solicitations`, () => HttpResponse.json(draftSolicitation, { status: 201 })),

  http.patch(`${BASE}/solicitations/:id`, () => HttpResponse.json(draftSolicitation)),

  http.post(`${BASE}/solicitations/:id/submit`, () =>
    HttpResponse.json(
      errorBody("validation_error", "preencha todos os campos antes de enviar", { categoryId: "obrigatório" }),
      { status: 422 },
    ),
  ),

  http.post(`${BASE}/solicitations/:id/approve`, () => HttpResponse.json(draftSolicitation)),

  http.post(`${BASE}/solicitations/:id/reject`, () => HttpResponse.json(draftSolicitation)),

  http.patch(`${BASE}/solicitations/:id/analysis`, () => HttpResponse.json(draftSolicitation)),

  http.post(`${BASE}/solicitations/:id/finalize`, () => HttpResponse.json(draftSolicitation)),

  http.get(`${BASE}/dashboard`, () => HttpResponse.json(dashboardSummary)),
];
