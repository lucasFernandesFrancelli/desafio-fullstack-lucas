import type { ReactElement } from "react";

import type { PendingActor, Status } from "@/types/api";

interface PendingActorBannerProps {
  status: Status;
  pendingActor: PendingActor;
}

export function PendingActorBanner({ status, pendingActor }: PendingActorBannerProps): ReactElement | null {
  if (pendingActor.kind === "none") {
    return null;
  }

  const label = status === "rascunho" ? "Aguardando conclusão de" : "Aguardando ação de";

  return (
    <div className="flex items-center gap-2 rounded-lg border border-primary/20 bg-primary/5 px-4 py-3 text-sm">
      <span className="size-1.5 shrink-0 animate-pulse rounded-full bg-primary" aria-hidden="true" />
      <span className="text-muted-foreground">{label} </span>
      <span className="font-medium text-foreground">{pendingActor.name}</span>
    </div>
  );
}
