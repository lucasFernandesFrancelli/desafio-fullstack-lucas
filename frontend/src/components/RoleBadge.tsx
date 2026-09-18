import type { ReactElement } from "react";

import { Badge } from "@/components/ui/badge";
import { ROLE_BADGE_CLASSES, ROLE_LABELS } from "@/lib/constants";
import { cn } from "@/lib/utils";
import type { Role } from "@/types/api";

interface RoleBadgeProps {
  role: Role;
  className?: string;
}

export function RoleBadge({ role, className }: RoleBadgeProps): ReactElement {
  return (
    <Badge variant="outline" className={cn("border", ROLE_BADGE_CLASSES[role], className)}>
      {ROLE_LABELS[role]}
    </Badge>
  );
}
