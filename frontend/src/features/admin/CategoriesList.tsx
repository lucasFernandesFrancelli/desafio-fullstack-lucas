import type { ReactElement } from "react";

import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import { EditApproversDialog } from "@/features/admin/EditApproversDialog";
import type { Category, User } from "@/types/api";

interface CategoriesListProps {
  categories: Category[];
  users: User[];
}

export function CategoriesList({ categories, users }: CategoriesListProps): ReactElement {
  if (categories.length === 0) {
    return <p className="text-sm text-muted-foreground">Nenhuma categoria cadastrada.</p>;
  }

  return (
    <div className="flex flex-col gap-3">
      {categories.map((category) => {
        const first = category.approvers.find((approver) => approver.order === 1);
        const second = category.approvers.find((approver) => approver.order === 2);

        return (
          <Card key={category.id}>
            <CardContent className="flex flex-col gap-3 p-4">
              <div className="flex flex-wrap items-start justify-between gap-2">
                <div>
                  <h3 className="text-sm font-medium">{category.name}</h3>
                  {category.description && (
                    <p className="text-xs text-muted-foreground">{category.description}</p>
                  )}
                </div>
                <EditApproversDialog category={category} users={users} />
              </div>

              <div className="flex flex-wrap gap-2">
                <Badge variant="outline">
                  1º aprovador: {first ? first.userName : <span className="text-destructive">não definido</span>}
                </Badge>
                <Badge variant="outline">
                  2º aprovador: {second ? second.userName : <span className="text-destructive">não definido</span>}
                </Badge>
              </div>
            </CardContent>
          </Card>
        );
      })}
    </div>
  );
}
