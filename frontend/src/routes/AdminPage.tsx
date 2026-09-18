import type { ReactElement } from "react";

import { useUsers } from "@/api/admin";
import { useCategories } from "@/api/categories";
import { Skeleton } from "@/components/ui/skeleton";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { CategoriesList } from "@/features/admin/CategoriesList";
import { NewCategoryDialog } from "@/features/admin/NewCategoryDialog";
import { NewUserDialog } from "@/features/admin/NewUserDialog";
import { UsersTable } from "@/features/admin/UsersTable";

export function AdminPage(): ReactElement {
  const categoriesQuery = useCategories(true);
  const usersQuery = useUsers(true);

  const users = usersQuery.data ?? [];

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-1">
        <h1 className="text-xl font-semibold tracking-tight">Gestão de cadastros</h1>
        <p className="text-sm text-muted-foreground">
          Defina quem aprova cada categoria e cadastre novas pessoas — para confirmar que a regra "dois
          aprovadores por categoria, em ordem" vem de dados reais, não de valores fixos no código.
        </p>
      </div>

      <Tabs defaultValue="categories">
        <TabsList>
          <TabsTrigger value="categories">Categorias e aprovadores</TabsTrigger>
          <TabsTrigger value="users">Usuários</TabsTrigger>
        </TabsList>

        <TabsContent value="categories" className="flex flex-col gap-4 pt-4">
          <div className="flex justify-end">
            <NewCategoryDialog users={users} />
          </div>

          {categoriesQuery.isPending && (
            <div className="flex flex-col gap-3">
              <Skeleton className="h-20 w-full" />
              <Skeleton className="h-20 w-full" />
            </div>
          )}
          {categoriesQuery.isError && (
            <p className="text-sm text-destructive">Não foi possível carregar as categorias.</p>
          )}
          {categoriesQuery.data && <CategoriesList categories={categoriesQuery.data} users={users} />}
        </TabsContent>

        <TabsContent value="users" className="flex flex-col gap-4 pt-4">
          <div className="flex justify-end">
            <NewUserDialog />
          </div>

          {usersQuery.isPending && <Skeleton className="h-48 w-full" />}
          {usersQuery.isError && <p className="text-sm text-destructive">Não foi possível carregar os usuários.</p>}
          {usersQuery.data && <UsersTable users={usersQuery.data} />}
        </TabsContent>
      </Tabs>
    </div>
  );
}
