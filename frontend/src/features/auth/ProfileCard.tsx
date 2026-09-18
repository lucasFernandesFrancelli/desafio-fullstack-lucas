import type { ReactElement } from "react";

import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import { ROLE_LABELS } from "@/lib/constants";
import type { UserProfile } from "@/types/api";

interface ProfileCardProps {
  profile: UserProfile;
  onSelect: (userId: string) => void;
  disabled?: boolean;
}

function initials(name: string): string {
  const parts = name.trim().split(/\s+/);
  const first = parts[0]?.[0] ?? "";
  const last = parts.length > 1 ? (parts[parts.length - 1]?.[0] ?? "") : "";
  return (first + last).toUpperCase();
}

export function ProfileCard({ profile, onSelect, disabled = false }: ProfileCardProps): ReactElement {
  return (
    <Card
      role="button"
      aria-disabled={disabled}
      onClick={() => {
        if (!disabled) onSelect(profile.id);
      }}
      className="cursor-pointer transition-colors hover:border-primary/50 hover:bg-muted/40"
    >
      <CardContent className="flex items-start gap-3 p-4">
        <Avatar>
          <AvatarFallback>{initials(profile.name)}</AvatarFallback>
        </Avatar>
        <div className="flex flex-1 flex-col gap-1.5">
          <div className="flex items-center justify-between gap-2">
            <span className="font-medium leading-none">{profile.name}</span>
            <Badge variant="secondary">{ROLE_LABELS[profile.role]}</Badge>
          </div>
          <span className="text-xs text-muted-foreground">{profile.email}</span>
          {profile.approverFor.length > 0 && (
            <div className="mt-1 flex flex-wrap gap-1">
              {profile.approverFor.map((approver) => (
                <Badge key={approver.categoryId} variant="outline" className="text-xs font-normal">
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
