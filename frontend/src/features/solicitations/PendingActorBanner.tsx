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
    <div className="rounded-lg border border-border bg-muted/40 px-4 py-3 text-sm">
      <span className="text-muted-foreground">{label} </span>
      <span className="font-medium">{pendingActor.name}</span>
    </div>
  );
}
