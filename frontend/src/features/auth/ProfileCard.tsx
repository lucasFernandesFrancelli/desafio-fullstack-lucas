import type { ReactElement } from "react";
import { Loader2Icon } from "lucide-react";

import { RoleBadge } from "@/components/RoleBadge";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import { ROLE_AVATAR_CLASSES } from "@/lib/constants";
import { cn } from "@/lib/utils";
import type { UserProfile } from "@/types/api";

interface ProfileCardProps {
  profile: UserProfile;
  onSelect: (userId: string) => void;
  disabled?: boolean;
  loading?: boolean;
}

function initials(name: string): string {
  const parts = name.trim().split(/\s+/);
  const first = parts[0]?.[0] ?? "";
  const last = parts.length > 1 ? (parts[parts.length - 1]?.[0] ?? "") : "";
  return (first + last).toUpperCase();
}

export function ProfileCard({ profile, onSelect, disabled = false, loading = false }: ProfileCardProps): ReactElement {
  return (
    <Card
      role="button"
      aria-disabled={disabled}
      aria-busy={loading}
      onClick={() => {
        if (!disabled) onSelect(profile.id);
      }}
      className={cn(
        "cursor-pointer border-transparent shadow-sm transition-all hover:-translate-y-0.5 hover:border-primary/40 hover:shadow-md",
        disabled && !loading && "opacity-60",
      )}
    >
      <CardContent className="flex items-start gap-3 p-4">
        <Avatar size="lg" className="relative">
          <AvatarFallback className={cn("font-semibold", ROLE_AVATAR_CLASSES[profile.role])}>
            {initials(profile.name)}
          </AvatarFallback>
          {loading && (
            <span className="absolute inset-0 flex items-center justify-center rounded-full bg-black/40">
              <Loader2Icon className="size-4 animate-spin text-white" aria-hidden="true" />
            </span>
          )}
        </Avatar>
        <div className="flex flex-1 flex-col gap-1.5">
          <div className="flex items-center justify-between gap-2">
            <span className="font-medium leading-none">{profile.name}</span>
            <RoleBadge role={profile.role} />
          </div>
          <span className="text-xs text-muted-foreground">{profile.email}</span>
          {profile.approverFor.length > 0 && (
            <div className="mt-1 flex flex-wrap gap-1">
              {profile.approverFor.map((approver) => (
                <Badge
                  key={approver.categoryId}
                  variant="outline"
                  className="border-amber-300 bg-amber-50 text-xs font-normal text-amber-800 dark:border-amber-800 dark:bg-amber-900/20 dark:text-amber-300"
                >
                  {approver.order}º aprovador · {approver.categoryName}
                </Badge>
              ))}
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
