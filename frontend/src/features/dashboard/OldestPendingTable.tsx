import type { ReactElement } from "react";
import { useNavigate } from "react-router";

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { StatusBadge } from "@/features/solicitations/StatusBadge";
import type { OldestPendingItem } from "@/types/api";

interface OldestPendingTableProps {
  items: OldestPendingItem[];
}

export function OldestPendingTable({ items }: OldestPendingTableProps): ReactElement {
  const navigate = useNavigate();

  if (items.length === 0) {
    return <p className="text-sm text-muted-foreground">Nenhuma solicitação parada no momento.</p>;
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Título</TableHead>
          <TableHead>Categoria</TableHead>
          <TableHead>Estado</TableHead>
          <TableHead>Aguardando</TableHead>
          <TableHead className="text-right">Dias parado</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {items.map((item) => (
          <TableRow
            key={item.id}
            className="cursor-pointer"
            onClick={() => navigate(`/solicitations/${item.id}`)}
          >
            <TableCell className="font-medium">{item.title}</TableCell>
            <TableCell>{item.categoryName || "—"}</TableCell>
            <TableCell>
              <StatusBadge status={item.status} />
            </TableCell>
            <TableCell>{item.pendingActorName || "—"}</TableCell>
            <TableCell className="text-right tabular-nums">{item.daysSinceLastMove}</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
