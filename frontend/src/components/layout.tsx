import { NavLink, Outlet } from "react-router-dom";
import { cn } from "@/lib/utils";
import { useTheme } from "@/contexts/theme-context";
import { LayoutDashboard, Terminal, Bookmark, FileText, Settings, PanelLeftClose, PanelLeft, Moon, Sun } from "lucide-react";
import { useState } from "react";

const navItems = [
  { to: "/", icon: LayoutDashboard, label: "Dashboard" },
  { to: "/query", icon: Terminal, label: "Query Runner" },
  { to: "/saved", icon: Bookmark, label: "Saved Queries" },
  { to: "/reports", icon: FileText, label: "Reports" },
  { to: "/settings", icon: Settings, label: "Settings" },
];

export default function Layout() {
  const [collapsed, setCollapsed] = useState(false);
  const { theme, setTheme } = useTheme();

  return (
    <div className="flex h-screen overflow-hidden relative z-10">
      {/* Sidebar: glassmorphism-lite, theme-aware */}
      <aside
        className={cn(
          "flex flex-col border-r transition-all duration-200",
          "bg-surface/95 dark:bg-surface/95 backdrop-blur-sm",
          "border-primary/10 shadow-[0_0_24px_rgba(0,0,0,0.12)] dark:shadow-[0_0_24px_rgba(0,0,0,0.2)]",
          "rounded-r-lg border-t-0 border-b-0 border-l-0",
          collapsed ? "w-16" : "w-56"
        )}
      >
        <div className="flex items-center gap-3 px-4 py-4 border-b border-border/80">
          <img src="/logo.png" alt="Logo" className="h-8 w-8 flex-shrink-0" />
          {!collapsed && <span className="text-sm font-bold text-brand-cyan tracking-wide">PgQueryNarrative</span>}
        </div>

        <nav className="flex-1 py-3 space-y-1 px-2">
          {navItems.map(({ to, icon: Icon, label }) => (
            <NavLink
              key={to}
              to={to}
              end={to === "/"}
              className={({ isActive }) =>
                cn(
                  "flex items-center gap-3 rounded-md px-3 py-2.5 text-sm font-medium transition-all duration-200",
                  "dark:hover:shadow-[inset_0_0_0_1px_rgba(0,245,255,0.08)]",
                  isActive
                    ? "bg-primary/10 text-primary shadow-[inset_2px_0_0_0] shadow-primary dark:shadow-[0_0_12px_rgba(0,245,255,0.15)]"
                    : "text-muted-foreground hover:text-foreground hover:bg-secondary/80"
                )
              }
            >
              <Icon className="h-4 w-4 flex-shrink-0" />
              {!collapsed && <span>{label}</span>}
            </NavLink>
          ))}
        </nav>

        {/* Theme toggle + collapse */}
        <div className="flex items-center border-t border-border/80">
          <button
            type="button"
            onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
            className={cn(
              "flex items-center justify-center text-muted-foreground hover:text-foreground hover:bg-secondary/50 transition-colors cursor-pointer",
              collapsed ? "flex-1 py-3" : "flex-1 gap-2 py-3"
            )}
            title={theme === "dark" ? "Switch to light mode" : "Switch to dark mode"}
          >
            {theme === "dark" ? <Sun className="h-4 w-4 shrink-0" /> : <Moon className="h-4 w-4 shrink-0" />}
            {!collapsed && <span className="text-xs font-medium">{theme === "dark" ? "Light" : "Dark"}</span>}
          </button>
          <button
            onClick={() => setCollapsed(!collapsed)}
            className="flex items-center justify-center p-3 text-muted-foreground hover:text-foreground transition-colors cursor-pointer rounded-br-lg hover:bg-secondary/50"
            title={collapsed ? "Expand sidebar" : "Collapse sidebar"}
          >
            {collapsed ? <PanelLeft className="h-4 w-4" /> : <PanelLeftClose className="h-4 w-4" />}
          </button>
        </div>
      </aside>

      {/* Main: content above background layers; subtle inner border in dark for depth */}
      <main className="flex-1 overflow-auto min-h-0 border-l border-transparent dark:border-border/30">
        <div className="max-w-6xl mx-auto px-6 py-8">
          <Outlet />
        </div>
      </main>
    </div>
  );
}
