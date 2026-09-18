import type { ReactElement } from "react";

import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";

interface PriorityBadgeProps {
  priority: number | null;
  className?: string;
}

/** Prioridade é o produto de gravidade × urgência × tendência (1 a 125). */
function priorityTier(priority: number): "alta" | "media" | "baixa" {
  if (priority >= 60) return "alta";
  if (priority >= 20) return "media";
  return "baixa";
}

const TIER_CLASSES: Record<ReturnType<typeof priorityTier>, string> = {
  alta: "bg-red-100 text-red-900 border-red-300 dark:bg-red-900/30 dark:text-red-300 dark:border-red-800",
  media:
    "bg-amber-100 text-amber-900 border-amber-300 dark:bg-amber-900/30 dark:text-amber-300 dark:border-amber-800",
  baixa:
    "bg-emerald-100 text-emerald-900 border-emerald-300 dark:bg-emerald-900/30 dark:text-emerald-300 dark:border-emerald-800",
};

export function PriorityBadge({ priority, className }: PriorityBadgeProps): ReactElement {
  if (priority === null) {
    return (
      <Badge variant="outline" className={cn("text-muted-foreground", className)}>
        Sem prioridade
      </Badge>
    );
  }

  const tier = priorityTier(priority);

  return (
    <Badge variant="outline" className={cn("border", TIER_CLASSES[tier], className)}>
      Prioridade {priority}
    </Badge>
  );
}
