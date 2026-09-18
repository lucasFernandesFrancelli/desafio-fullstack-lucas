import type { ReactElement } from "react";
import { useNavigate } from "react-router";

import { Card, CardContent } from "@/components/ui/card";
import { PriorityBadge } from "@/features/solicitations/PriorityBadge";
import { StatusBadge } from "@/features/solicitations/StatusBadge";
import { STATUS_ACCENT_BORDER_CLASSES } from "@/lib/constants";
import { formatDateTime } from "@/lib/format";
import { cn } from "@/lib/utils";
import type { SolicitationSummary } from "@/types/api";

interface SolicitationCardProps {
  solicitation: SolicitationSummary;
  /** Esconde o badge de status quando o card já está dentro de uma coluna do Kanban que representa o status. */
  hideStatus?: boolean;
}

export function SolicitationCard({ solicitation, hideStatus = false }: SolicitationCardProps): ReactElement {
  const navigate = useNavigate();

  return (
    <Card
      role="button"
      tabIndex={0}
      data-testid={`solicitation-card-${solicitation.id}`}
      onClick={() => navigate(`/solicitations/${solicitation.id}`)}
      onKeyDown={(event) => {
        if (event.key === "Enter" || event.key === " ") {
          navigate(`/solicitations/${solicitation.id}`);
        }
      }}
      className={cn(
        "cursor-pointer border-l-4 shadow-sm transition-all hover:-translate-y-0.5 hover:shadow-md",
        STATUS_ACCENT_BORDER_CLASSES[solicitation.status],
      )}
    >
      <CardContent className="flex flex-col gap-2 p-4">
        <div className="flex items-start justify-between gap-2">
          <h3 className="text-sm font-medium leading-snug">{solicitation.title}</h3>
          {!hideStatus && <StatusBadge status={solicitation.status} />}
        </div>

        {solicitation.categoryName && (
          <p className="text-xs text-muted-foreground">{solicitation.categoryName}</p>
        )}

        <div className="flex flex-wrap items-center gap-2">
          <PriorityBadge priority={solicitation.priority} />
        </div>

        <div className="flex items-center justify-between text-xs text-muted-foreground">
          <span>Por {solicitation.requesterName}</span>
          <span>{formatDateTime(solicitation.updatedAt)}</span>
        </div>

        {solicitation.pendingActorName && (
          <p className="rounded-md bg-primary/5 px-2 py-1 text-xs text-muted-foreground">
            Aguardando: <span className="font-medium text-foreground">{solicitation.pendingActorName}</span>
          </p>
        )}
      </CardContent>
    </Card>
  );
}
