import { NavLink, Outlet } from "react-router-dom";
import { cn } from "@/lib/utils";
import { LayoutDashboard, Terminal, Bookmark, FileText, Settings, PanelLeftClose, PanelLeft } from "lucide-react";
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

  return (
    <div className="flex h-screen overflow-hidden">
      {/* Sidebar */}
      <aside className={cn("flex flex-col border-r border-border bg-surface transition-all duration-200", collapsed ? "w-16" : "w-56")}>
        <div className="flex items-center gap-3 px-4 py-4 border-b border-border">
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
                  "flex items-center gap-3 rounded-md px-3 py-2.5 text-sm font-medium transition-colors",
                  isActive
                    ? "bg-primary/10 text-primary shadow-[inset_2px_0_0_0] shadow-primary"
                    : "text-muted-foreground hover:text-foreground hover:bg-secondary"
                )
              }
            >
              <Icon className="h-4 w-4 flex-shrink-0" />
              {!collapsed && <span>{label}</span>}
            </NavLink>
          ))}
        </nav>

        <button
          onClick={() => setCollapsed(!collapsed)}
          className="flex items-center justify-center p-3 border-t border-border text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
        >
          {collapsed ? <PanelLeft className="h-4 w-4" /> : <PanelLeftClose className="h-4 w-4" />}
        </button>
      </aside>

      {/* Main */}
      <main className="flex-1 overflow-auto">
        <div className="max-w-6xl mx-auto px-6 py-8">
          <Outlet />
        </div>
      </main>
    </div>
  );
}
