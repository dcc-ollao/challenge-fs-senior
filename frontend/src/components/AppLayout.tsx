import { Link, Outlet, useLocation } from "react-router-dom";
import { useAuth } from "../features/auth/useAuth";

function NavLink({
  to,
  label,
}: {
  to: string;
  label: string;
}) {
  const { pathname } = useLocation();
  const active = pathname === to;

  return (
    <Link
      to={to}
      className={[
        "text-sm font-medium",
        active ? "underline" : "hover:underline",
      ].join(" ")}
    >
      {label}
    </Link>
  );
}

export function AppLayout() {
  const { logout, user } = useAuth();

  return (
    <div className="min-h-screen bg-white text-slate-900">
      <header className="border-b">
        <div className="mx-auto flex max-w-5xl items-center justify-between px-4 py-3">
          <div className="flex items-center gap-6">
            <div className="text-base font-semibold">Task Management</div>

            <nav className="flex items-center gap-4">
              <NavLink to="/" label="Home" />
              <NavLink to="/tasks" label="Tasks" />
              <NavLink to="/projects" label="Projects" />
            </nav>
          </div>

          <div className="flex items-center gap-3">
            <span className="text-xs text-slate-600">
              {user?.email ?? ""}
            </span>
            <button
              onClick={logout}
              className="rounded-md border px-3 py-1.5 text-sm hover:bg-slate-50"
              type="button"
            >
              Logout
            </button>
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-5xl px-4 py-8">
        <Outlet />
      </main>
    </div>
  );
}
