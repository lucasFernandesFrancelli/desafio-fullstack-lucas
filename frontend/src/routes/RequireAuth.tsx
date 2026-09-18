import type { ReactElement } from "react";
import { Navigate } from "react-router";

import { PageSpinner } from "@/components/PageSpinner";
import { useAuth } from "@/context/AuthContext";

interface RequireAuthProps {
  children: ReactElement;
}

export function RequireAuth({ children }: RequireAuthProps): ReactElement {
  const { user, isLoading } = useAuth();

  if (isLoading) {
    return <PageSpinner fullScreen label="Carregando sessão…" />;
  }

  if (!user) {
    return <Navigate to="/" replace />;
  }

  return children;
}
