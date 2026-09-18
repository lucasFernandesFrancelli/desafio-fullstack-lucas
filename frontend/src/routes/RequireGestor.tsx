import type { ReactElement } from "react";
import { Navigate } from "react-router";

import { useAuth } from "@/context/AuthContext";

interface RequireGestorProps {
  children: ReactElement;
}

export function RequireGestor({ children }: RequireGestorProps): ReactElement {
  const { user, isLoading } = useAuth();

  // Evita redirecionar durante a reidratação da sessão (antes de saber o
  // papel do usuário) — RequireGestor normalmente já roda dentro de
  // RequireAuth, mas fica seguro mesmo se usado isoladamente.
  if (isLoading) {
    return (
      <div className="flex min-h-svh items-center justify-center text-muted-foreground">
        Carregando sessão…
      </div>
    );
  }

  if (user?.role !== "gestor") {
    return <Navigate to="/solicitations" replace />;
  }

  return children;
}
