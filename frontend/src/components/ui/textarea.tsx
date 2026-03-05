import { forwardRef, type TextareaHTMLAttributes } from "react";
import { cn } from "@/lib/utils";

/* Textarea: same focus ring and border contrast as Input */
const Textarea = forwardRef<HTMLTextAreaElement, TextareaHTMLAttributes<HTMLTextAreaElement>>(({ className, ...props }, ref) => (
  <textarea
    className={cn(
      "flex min-h-[120px] w-full rounded-md border border-input bg-background/80 px-3 py-2 text-sm text-foreground font-mono placeholder:text-muted-foreground resize-y",
      "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/40 focus-visible:ring-offset-2 focus-visible:ring-offset-background focus-visible:border-primary/30",
      "disabled:cursor-not-allowed disabled:opacity-50 transition-colors",
      className
    )}
    ref={ref}
    {...props}
  />
));
Textarea.displayName = "Textarea";

export { Textarea };
