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
        "text-sm font-medium transition",
        active
          ? "text-slate-900"
          : "text-slate-600 hover:text-slate-900",
      ].join(" ")}
    >
      {label}
    </Link>
  );
}

export function AppLayout() {
  const { logout, user } = useAuth();

  return (
    <div className="min-h-screen bg-slate-50 text-slate-900">
      <header className="border-b bg-white">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
          <div className="flex items-center gap-8">
            <div className="text-lg font-semibold tracking-tight">
              Task Management
            </div>

            <nav className="flex items-center gap-5">
              <NavLink to="/" label="Home" />
              <NavLink to="/tasks" label="Tasks" />
              <NavLink to="/projects" label="Projects" />
              {user?.role === "admin" && (
                <NavLink to="/admin/users" label="Users" />
              )}
            </nav>
          </div>

          <div className="flex items-center gap-4">
            <span className="text-xs text-slate-500">
              {user?.email}
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

      <main className="mx-auto max-w-6xl px-6 py-10">
        <Outlet />
      </main>
    </div>
  );
}
