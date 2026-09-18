import { useState, type ReactElement } from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";

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
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { rejectSchema, type RejectFormValues } from "@/schemas/rejectSchema";

interface ApproveRejectPanelProps {
  onApprove: () => Promise<void>;
  onReject: (reason: string) => Promise<void>;
  isApproving: boolean;
  isRejecting: boolean;
}

export function ApproveRejectPanel({
  onApprove,
  onReject,
  isApproving,
  isRejecting,
}: ApproveRejectPanelProps): ReactElement {
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const form = useForm<RejectFormValues>({
    resolver: zodResolver(rejectSchema),
    defaultValues: { reason: "" },
  });

  async function handleReject(values: RejectFormValues): Promise<void> {
    await onReject(values.reason);
    setIsDialogOpen(false);
    form.reset();
  }

  return (
    <div className="flex flex-col gap-3 rounded-lg border border-border p-4">
      <p className="text-sm font-medium">É a sua vez de decidir esta solicitação.</p>
      <div className="flex gap-2">
        <Button onClick={() => void onApprove()} disabled={isApproving || isRejecting}>
          {isApproving ? "Aprovando…" : "Aprovar"}
        </Button>

        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogTrigger asChild>
            <Button variant="destructive" disabled={isApproving || isRejecting}>
              Recusar
            </Button>
          </DialogTrigger>
          <DialogContent>
            <form onSubmit={form.handleSubmit(handleReject)} noValidate>
              <DialogHeader>
                <DialogTitle>Recusar solicitação</DialogTitle>
                <DialogDescription>
                  A recusa encerra o fluxo desta solicitação. Informe o motivo para o solicitante.
                </DialogDescription>
              </DialogHeader>

              <div className="flex flex-col gap-1.5 py-4">
                <Label htmlFor="reject-reason">Motivo</Label>
                <Textarea id="reject-reason" rows={4} {...form.register("reason")} />
                {form.formState.errors.reason && (
                  <p className="text-xs text-destructive">{form.formState.errors.reason.message}</p>
                )}
              </div>

              <DialogFooter>
                <Button type="submit" variant="destructive" disabled={isRejecting}>
                  {isRejecting ? "Recusando…" : "Confirmar recusa"}
                </Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      </div>
    </div>
  );
}
