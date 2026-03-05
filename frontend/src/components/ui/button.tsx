import { forwardRef, type ButtonHTMLAttributes } from "react";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "@/lib/utils";

/* Buttons: subtle cyan glow on hover/focus, active press state (translateY) */
const buttonVariants = cva(
  "inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring/60 focus-visible:ring-offset-2 focus-visible:ring-offset-background disabled:pointer-events-none disabled:opacity-50 cursor-pointer active:translate-y-px",
  {
    variants: {
      variant: {
        default:
          "bg-primary text-primary-foreground hover:bg-brand-cyan-soft hover:shadow-[0_0_16px_rgba(0,245,255,0.35)] focus-visible:shadow-[0_0_16px_rgba(0,245,255,0.3)]",
        secondary: "bg-secondary text-secondary-foreground hover:bg-secondary/80 hover:shadow-[0_0_12px_rgba(0,245,255,0.15)]",
        ghost: "hover:bg-secondary hover:text-foreground",
        destructive: "bg-destructive text-destructive-foreground hover:bg-destructive/90",
        outline:
          "border border-border bg-transparent hover:bg-secondary hover:text-foreground hover:border-primary/20 hover:shadow-[0_0_10px_rgba(0,245,255,0.1)]",
        link: "text-primary underline-offset-4 hover:underline",
      },
      size: {
        default: "h-10 px-4 py-2",
        sm: "h-8 rounded-md px-3 text-xs",
        lg: "h-11 rounded-md px-8 text-base",
        icon: "h-10 w-10",
      },
    },
    defaultVariants: { variant: "default", size: "default" },
  }
);

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement>, VariantProps<typeof buttonVariants> {}

const Button = forwardRef<HTMLButtonElement, ButtonProps>(({ className, variant, size, ...props }, ref) => (
  <button className={cn(buttonVariants({ variant, size, className }))} ref={ref} {...props} />
));
Button.displayName = "Button";

export { Button, buttonVariants };
