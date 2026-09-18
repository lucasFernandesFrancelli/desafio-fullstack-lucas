import { useState, type ReactElement } from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import { toast } from "sonner";

import { useSetCategoryApprovers } from "@/api/admin";
import { ApiError } from "@/api/httpClient";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { ApproverSelect } from "@/features/admin/ApproverSelect";
import { setApproversSchema, type SetApproversFormValues } from "@/schemas/adminSchemas";
import type { Category, User } from "@/types/api";

interface EditApproversDialogProps {
  category: Category;
  users: User[];
}

export function EditApproversDialog({ category, users }: EditApproversDialogProps): ReactElement {
  const [open, setOpen] = useState(false);
  const setApprovers = useSetCategoryApprovers();

  const currentFirst = category.approvers.find((approver) => approver.order === 1);
  const currentSecond = category.approvers.find((approver) => approver.order === 2);

  const form = useForm<SetApproversFormValues>({
    resolver: zodResolver(setApproversSchema),
    defaultValues: {
      firstApproverId: currentFirst?.userId ?? "",
      secondApproverId: currentSecond?.userId ?? "",
    },
  });

  async function handleSubmit(values: SetApproversFormValues): Promise<void> {
    try {
      await setApprovers.mutateAsync({ id: category.id, input: values });
      toast.success("Aprovadores atualizados.");
      setOpen(false);
    } catch (error) {
      toast.error(error instanceof ApiError ? error.message : "Não foi possível atualizar os aprovadores.");
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button size="sm" variant="outline">
          Editar aprovadores
        </Button>
      </DialogTrigger>
      <DialogContent>
        <form onSubmit={form.handleSubmit(handleSubmit)} noValidate>
          <DialogHeader>
            <DialogTitle>Aprovadores — {category.name}</DialogTitle>
          </DialogHeader>

          <div className="flex flex-col gap-4 py-4">
            <div className="flex flex-col gap-1.5">
              <Label htmlFor={`edit-approvers-first-${category.id}`}>1º aprovador</Label>
              <Controller
                control={form.control}
                name="firstApproverId"
                render={({ field }) => (
                  <ApproverSelect
                    id={`edit-approvers-first-${category.id}`}
                    users={users}
                    value={field.value}
                    onChange={field.onChange}
                  />
                )}
              />
              {form.formState.errors.firstApproverId && (
                <p className="text-xs text-destructive">{form.formState.errors.firstApproverId.message}</p>
              )}
            </div>

            <div className="flex flex-col gap-1.5">
              <Label htmlFor={`edit-approvers-second-${category.id}`}>2º aprovador</Label>
              <Controller
                control={form.control}
                name="secondApproverId"
                render={({ field }) => (
                  <ApproverSelect
                    id={`edit-approvers-second-${category.id}`}
                    users={users}
                    value={field.value}
                    onChange={field.onChange}
                  />
                )}
              />
              {form.formState.errors.secondApproverId && (
                <p className="text-xs text-destructive">{form.formState.errors.secondApproverId.message}</p>
              )}
            </div>
          </div>

          <DialogFooter>
            <Button type="submit" disabled={setApprovers.isPending}>
              {setApprovers.isPending ? "Salvando…" : "Salvar"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
