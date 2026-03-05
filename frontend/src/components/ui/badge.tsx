import { type HTMLAttributes } from "react";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "@/lib/utils";

/* Badge: subtle border and shadow for depth, theme-aware */
const badgeVariants = cva(
  "inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-all duration-200 shadow-sm",
  {
    variants: {
      variant: {
        default:
          "border-primary/20 bg-primary/15 text-primary dark:border-primary/25 dark:shadow-[0_0_8px_rgba(0,245,255,0.08)]",
        secondary: "border-border/80 bg-secondary text-secondary-foreground",
        success: "border-success/25 bg-success/15 text-success",
        warning: "border-warning/25 bg-warning/15 text-warning",
        destructive: "border-destructive/25 bg-destructive/15 text-destructive",
        outline: "border-border text-foreground bg-transparent hover:bg-secondary/50",
      },
    },
    defaultVariants: { variant: "default" },
  }
);

export interface BadgeProps extends HTMLAttributes<HTMLDivElement>, VariantProps<typeof badgeVariants> {}

export function Badge({ className, variant, ...props }: BadgeProps) {
  return <div className={cn(badgeVariants({ variant }), className)} {...props} />;
}
