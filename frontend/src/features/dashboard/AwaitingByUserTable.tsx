import type { ReactElement } from "react";

import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import type { UserPendingCount } from "@/types/api";

interface AwaitingByUserTableProps {
  items: UserPendingCount[];
}

export function AwaitingByUserTable({ items }: AwaitingByUserTableProps): ReactElement {
  if (items.length === 0) {
    return <p className="text-sm text-muted-foreground">Ninguém com aprovações pendentes no momento.</p>;
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Pessoa</TableHead>
          <TableHead className="text-right">Pendências</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {items.map((item) => (
          <TableRow key={item.userId}>
            <TableCell className="font-medium">{item.userName}</TableCell>
            <TableCell className="text-right">
              <Badge className="border-amber-300 bg-amber-100 text-amber-900 dark:border-amber-800 dark:bg-amber-900/30 dark:text-amber-300" variant="outline">
                {item.count}
              </Badge>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
