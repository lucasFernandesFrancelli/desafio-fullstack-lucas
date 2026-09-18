import type { ReactElement } from "react";
import { Navigate } from "react-router";

import { useAuth } from "@/context/AuthContext";

interface RequireAuthProps {
  children: ReactElement;
}

export function RequireAuth({ children }: RequireAuthProps): ReactElement {
  const { user, isLoading } = useAuth();

  if (isLoading) {
    return (
      <div className="flex min-h-svh items-center justify-center text-muted-foreground">
        Carregando sessão…
      </div>
    );
  }

  if (!user) {
    return <Navigate to="/" replace />;
  }

  return children;
}
