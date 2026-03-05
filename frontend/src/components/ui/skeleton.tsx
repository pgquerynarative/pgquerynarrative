import { cn } from "@/lib/utils";

/* Skeleton: theme-aware pulse, slightly darker in dark mode for contrast */
export function Skeleton({ className, ...props }: React.HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn(
        "animate-pulse rounded-md bg-muted/80 dark:bg-muted dark:bg-opacity-90",
        "ring-1 ring-inset ring-black/[0.03] dark:ring-white/[0.04]",
        className
      )}
      {...props}
    />
  );
}
