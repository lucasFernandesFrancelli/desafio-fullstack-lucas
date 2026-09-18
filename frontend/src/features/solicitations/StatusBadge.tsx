import type { ReactElement } from "react";

import { Badge } from "@/components/ui/badge";
import { STATUS_BADGE_CLASSES, STATUS_LABELS } from "@/lib/constants";
import { cn } from "@/lib/utils";
import type { Status } from "@/types/api";

interface StatusBadgeProps {
  status: Status;
  className?: string;
}

export function StatusBadge({ status, className }: StatusBadgeProps): ReactElement {
  return (
    <Badge variant="outline" className={cn("border", STATUS_BADGE_CLASSES[status], className)}>
      {STATUS_LABELS[status]}
    </Badge>
  );
}
