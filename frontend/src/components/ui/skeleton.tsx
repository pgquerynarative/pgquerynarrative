import { cn } from "@/lib/utils";

/* Skeleton: theme-aware pulse with smooth animation, visible in both themes */
export function Skeleton({ className, ...props }: React.HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn(
        "rounded-md bg-muted/80 dark:bg-muted/90",
        "ring-1 ring-inset ring-black/[0.04] dark:ring-white/[0.05]",
        "animate-pulse motion-reduce:animate-none",
        className
      )}
      {...props}
    />
  );
}
