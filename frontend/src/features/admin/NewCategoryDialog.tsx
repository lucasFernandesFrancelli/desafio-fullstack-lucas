import { useState, type ReactElement } from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import { toast } from "sonner";

import { useCreateCategory } from "@/api/admin";
import { ApiError } from "@/api/httpClient";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { ApproverSelect } from "@/features/admin/ApproverSelect";
import { createCategorySchema, type CreateCategoryFormValues } from "@/schemas/adminSchemas";
import type { User } from "@/types/api";

interface NewCategoryDialogProps {
  users: User[];
}

export function NewCategoryDialog({ users }: NewCategoryDialogProps): ReactElement {
  const [open, setOpen] = useState(false);
  const createCategory = useCreateCategory();
  const form = useForm<CreateCategoryFormValues>({
    resolver: zodResolver(createCategorySchema),
    defaultValues: { name: "", description: "", firstApproverId: "", secondApproverId: "" },
  });

  async function handleSubmit(values: CreateCategoryFormValues): Promise<void> {
    try {
      await createCategory.mutateAsync(values);
      toast.success("Categoria criada.");
      form.reset();
      setOpen(false);
    } catch (error) {
      toast.error(error instanceof ApiError ? error.message : "Não foi possível criar a categoria.");
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button size="sm">Nova categoria</Button>
      </DialogTrigger>
      <DialogContent>
        <form onSubmit={form.handleSubmit(handleSubmit)} noValidate>
          <DialogHeader>
            <DialogTitle>Nova categoria</DialogTitle>
            <DialogDescription>
              Toda categoria precisa dos dois aprovadores definidos já na criação.
            </DialogDescription>
          </DialogHeader>

          <div className="flex flex-col gap-4 py-4">
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="new-category-name">Nome</Label>
              <Input id="new-category-name" {...form.register("name")} />
              {form.formState.errors.name && (
                <p className="text-xs text-destructive">{form.formState.errors.name.message}</p>
              )}
            </div>

            <div className="flex flex-col gap-1.5">
              <Label htmlFor="new-category-description">Descrição</Label>
              <Textarea id="new-category-description" rows={2} {...form.register("description")} />
            </div>

            <div className="flex flex-col gap-1.5">
              <Label htmlFor="new-category-first-approver">1º aprovador</Label>
              <Controller
                control={form.control}
                name="firstApproverId"
                render={({ field }) => (
                  <ApproverSelect id="new-category-first-approver" users={users} value={field.value} onChange={field.onChange} />
                )}
              />
              {form.formState.errors.firstApproverId && (
                <p className="text-xs text-destructive">{form.formState.errors.firstApproverId.message}</p>
              )}
            </div>

            <div className="flex flex-col gap-1.5">
              <Label htmlFor="new-category-second-approver">2º aprovador</Label>
              <Controller
                control={form.control}
                name="secondApproverId"
                render={({ field }) => (
                  <ApproverSelect id="new-category-second-approver" users={users} value={field.value} onChange={field.onChange} />
                )}
              />
              {form.formState.errors.secondApproverId && (
                <p className="text-xs text-destructive">{form.formState.errors.secondApproverId.message}</p>
              )}
            </div>
          </div>

          <DialogFooter>
            <Button type="submit" disabled={createCategory.isPending}>
              {createCategory.isPending ? "Criando…" : "Criar categoria"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
