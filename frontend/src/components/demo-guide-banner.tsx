import { Button } from "@/components/ui/button";
import { X } from "lucide-react";

/** One-line tip after a guided demo — the page panels already carry the story. */
export function DemoGuideBanner({ onDismiss }: { onDismiss: () => void }) {
  return (
    <div
      className="flex items-start gap-3 rounded-lg border border-primary/25 bg-primary/5 px-4 py-3"
      data-testid="demo-guide"
    >
      <p className="flex-1 text-sm text-muted-foreground leading-relaxed">
        <span className="font-medium text-foreground">Guided demo finished. </span>
        Scroll the findings, candidate rewrite, and proof table below — <span className="font-medium text-foreground">VerifiedEqual</span>{" "}
        means every row matched; <span className="font-medium text-foreground">SampleMatch</span> means a bounded sample did.
      </p>
      <Button
        type="button"
        variant="ghost"
        size="icon"
        className="shrink-0"
        aria-label="Dismiss demo guide"
        data-testid="demo-guide-dismiss"
        onClick={onDismiss}
      >
        <X className="h-4 w-4" />
      </Button>
    </div>
  );
}
