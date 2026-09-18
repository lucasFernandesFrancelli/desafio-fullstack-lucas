import type { ReactElement } from "react";

import { RoleBadge } from "@/components/RoleBadge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import type { User } from "@/types/api";

interface UsersTableProps {
  users: User[];
}

export function UsersTable({ users }: UsersTableProps): ReactElement {
  if (users.length === 0) {
    return <p className="text-sm text-muted-foreground">Nenhum usuário cadastrado.</p>;
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Nome</TableHead>
          <TableHead>E-mail</TableHead>
          <TableHead>Papel</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {users.map((user) => (
          <TableRow key={user.id}>
            <TableCell className="font-medium">{user.name}</TableCell>
            <TableCell className="text-muted-foreground">{user.email}</TableCell>
            <TableCell>
              <RoleBadge role={user.role} />
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
