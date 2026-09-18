import { useState, type ReactElement } from "react";
import { Navigate, useNavigate } from "react-router";
import { toast } from "sonner";

import { useProfiles } from "@/api/auth";
import { PageSpinner } from "@/components/PageSpinner";
import { Skeleton } from "@/components/ui/skeleton";
import { useAuth } from "@/context/AuthContext";
import { ProfileCard } from "@/features/auth/ProfileCard";

export function ProfileSelectPage(): ReactElement {
  const { user, isLoading: isSessionLoading, login } = useAuth();
  const navigate = useNavigate();
  const profilesQuery = useProfiles();
  const [loggingInId, setLoggingInId] = useState<string | null>(null);

  if (isSessionLoading) {
    return <PageSpinner fullScreen label="Carregando sessão…" />;
  }

  if (user) {
    return <Navigate to="/solicitations" replace />;
  }

  async function handleSelect(userId: string): Promise<void> {
    setLoggingInId(userId);
    try {
      await login(userId);
      navigate("/solicitations", { replace: true });
    } catch {
      toast.error("Não foi possível entrar com este perfil. Tente novamente.");
    } finally {
      setLoggingInId(null);
    }
  }

  return (
    <div className="mx-auto flex min-h-svh max-w-3xl flex-col justify-center gap-6 px-4 py-10">
      <div className="overflow-hidden rounded-2xl bg-brand-gradient px-6 py-10 text-center shadow-lg">
        <span className="text-xs font-semibold tracking-wide text-white/80 uppercase">eKaizen · Melhoria contínua</span>
        <h1 className="mt-1 text-2xl font-semibold tracking-tight text-white sm:text-3xl">Da ideia à decisão</h1>
        <p className="mx-auto mt-2 max-w-lg text-sm text-white/90">
          Escolha um perfil para entrar. Cada um representa uma pessoa diferente no processo de melhoria
          contínua — solicitante, aprovador, analista ou gestor.
        </p>
      </div>

      {profilesQuery.isPending && (
        <div className="grid gap-3 sm:grid-cols-2">
          {Array.from({ length: 6 }, (_, index) => (
            <Skeleton key={`profile-skeleton-${index}`} className="h-24 w-full rounded-xl" />
          ))}
        </div>
      )}

      {profilesQuery.isError && (
        <p className="text-center text-sm text-destructive">
          Não foi possível carregar os perfis. Verifique se a API está no ar e recarregue a página.
        </p>
      )}

      {profilesQuery.data && (
        <div className="grid gap-3 sm:grid-cols-2">
          {profilesQuery.data.map((profile) => (
            <ProfileCard
              key={profile.id}
              profile={profile}
              onSelect={handleSelect}
              disabled={loggingInId !== null}
              loading={loggingInId === profile.id}
            />
          ))}
        </div>
      )}
    </div>
  );
}
