import type { ReactElement, ReactNode } from "react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, type RenderResult } from "@testing-library/react";
import { MemoryRouter } from "react-router";

import { AuthProvider } from "@/context/AuthContext";

interface RenderOptions {
  route?: string;
  withAuth?: boolean;
}

function createTestQueryClient(): QueryClient {
  return new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });
}

function Providers({ children, withAuth }: { children: ReactNode; withAuth: boolean }): ReactElement {
  const client = createTestQueryClient();
  const content = withAuth ? <AuthProvider>{children}</AuthProvider> : <>{children}</>;
  return <QueryClientProvider client={client}>{content}</QueryClientProvider>;
}

export function renderWithProviders(ui: ReactElement, options: RenderOptions = {}): RenderResult {
  const { route = "/", withAuth = true } = options;
  window.history.pushState({}, "", route);

  return render(
    <Providers withAuth={withAuth}>
      <MemoryRouter initialEntries={[route]}>{ui}</MemoryRouter>
    </Providers>,
  );
}
