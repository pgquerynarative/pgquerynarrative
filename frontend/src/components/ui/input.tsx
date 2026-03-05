import { forwardRef, type InputHTMLAttributes } from "react";
import { cn } from "@/lib/utils";

/* Inputs: clearer focus ring (cyan low opacity), improved padding and border contrast */
const Input = forwardRef<HTMLInputElement, InputHTMLAttributes<HTMLInputElement>>(({ className, type, ...props }, ref) => (
  <input
    type={type}
    className={cn(
      "flex h-10 w-full rounded-md border border-input bg-background/80 px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground",
      "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/40 focus-visible:ring-offset-2 focus-visible:ring-offset-background focus-visible:border-primary/30",
      "disabled:cursor-not-allowed disabled:opacity-50 transition-colors",
      className
    )}
    ref={ref}
    {...props}
  />
));
Input.displayName = "Input";

export { Input };
