import type { ReactElement } from "react";

import { HISTORY_ACTION_DOT_CLASSES, HISTORY_ACTION_LABELS } from "@/lib/constants";
import { formatDateTime } from "@/lib/format";
import { cn } from "@/lib/utils";
import type { HistoryEntry } from "@/types/api";

interface HistoryTimelineProps {
  history: HistoryEntry[];
}

export function HistoryTimeline({ history }: HistoryTimelineProps): ReactElement {
  if (history.length === 0) {
    return <p className="text-sm text-muted-foreground">Ainda não há histórico.</p>;
  }

  return (
    <ol className="flex flex-col gap-4">
      {history.map((entry) => (
        <li key={entry.id} className="flex gap-3">
          <div
            className={cn("mt-1.5 h-2 w-2 shrink-0 rounded-full", HISTORY_ACTION_DOT_CLASSES[entry.action])}
            aria-hidden="true"
          />
          <div className="flex flex-1 flex-col gap-0.5">
            <div className="flex flex-wrap items-baseline justify-between gap-x-3">
              <span className="text-sm font-medium">{HISTORY_ACTION_LABELS[entry.action]}</span>
              <span className="text-xs text-muted-foreground">{formatDateTime(entry.createdAt)}</span>
            </div>
            <span className="text-xs text-muted-foreground">por {entry.actorName}</span>

            {entry.action === "recusada" && entry.comment && (
              <p className="mt-1 rounded-md bg-muted px-2 py-1 text-sm">Motivo: {entry.comment}</p>
            )}

            {entry.action === "finalizada" && (
              <div className="mt-1 flex flex-col gap-1 rounded-md bg-muted px-2 py-1.5 text-sm">
                {entry.comment && <p>Parecer: {entry.comment}</p>}
                <p className="text-xs text-muted-foreground">
                  Gravidade {entry.severity} · Urgência {entry.urgency} · Tendência {entry.trend}
                </p>
              </div>
            )}
          </div>
        </li>
      ))}
    </ol>
  );
}
