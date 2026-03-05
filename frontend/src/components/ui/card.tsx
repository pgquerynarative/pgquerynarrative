import { forwardRef, type HTMLAttributes } from "react";
import { cn } from "@/lib/utils";

/* Glassmorphism-lite: theme-aware surface, border, shadow; hover lift + subtle glow in dark */
const Card = forwardRef<HTMLDivElement, HTMLAttributes<HTMLDivElement>>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn(
      "rounded-lg text-card-foreground transition-all duration-200",
      "bg-card/90 dark:bg-card/90 backdrop-blur-sm",
      "border border-primary/10 shadow-md dark:shadow-md dark:shadow-black/20",
      "ring-1 ring-inset ring-black/[0.02] dark:ring-white/[0.03]",
      "hover:shadow-lg dark:hover:shadow-[0_0_28px_rgba(0,0,0,0.25)] hover:border-primary/15",
      "hover:translate-y-px",
      className
    )}
    {...props}
  />
));
Card.displayName = "Card";

const CardHeader = forwardRef<HTMLDivElement, HTMLAttributes<HTMLDivElement>>(({ className, ...props }, ref) => (
  <div ref={ref} className={cn("flex flex-col space-y-1.5 p-4 sm:p-5", className)} {...props} />
));
CardHeader.displayName = "CardHeader";

const CardTitle = forwardRef<HTMLHeadingElement, HTMLAttributes<HTMLHeadingElement>>(({ className, ...props }, ref) => (
  <h3 ref={ref} className={cn("text-base font-semibold leading-tight tracking-tight", className)} {...props} />
));
CardTitle.displayName = "CardTitle";

const CardDescription = forwardRef<HTMLParagraphElement, HTMLAttributes<HTMLParagraphElement>>(({ className, ...props }, ref) => (
  <p ref={ref} className={cn("text-sm text-muted-foreground leading-snug", className)} {...props} />
));
CardDescription.displayName = "CardDescription";

const CardContent = forwardRef<HTMLDivElement, HTMLAttributes<HTMLDivElement>>(({ className, ...props }, ref) => (
  <div ref={ref} className={cn("p-4 sm:p-5 pt-0", className)} {...props} />
));
CardContent.displayName = "CardContent";

export { Card, CardHeader, CardTitle, CardDescription, CardContent };
