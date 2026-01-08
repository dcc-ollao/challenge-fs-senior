import { useEffect, useMemo, useState } from "react";
import { listUsers, updateUserRole } from "./api";
import type { User, UserRole } from "./types";
import { useAuthContext } from "../auth/AuthContext";
import { useSnackbar } from "../../components/snackbar/SnackbarContext";

const ROLE_OPTIONS: { value: UserRole; label: string }[] = [
  { value: "user", label: "User" },
  { value: "admin", label: "Admin" },
];

function getErrorMessage(err: any, fallback: string) {
  return (
    err?.response?.data?.message ??
    err?.response?.data?.error ??
    err?.message ??
    fallback
  );
}

export default function UsersAdminPage() {
  const { user: sessionUser } = useAuthContext();
  const { showError, showSuccess } = useSnackbar();

  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [savingById, setSavingById] = useState<Record<string, boolean>>({});

  async function load() {
    setLoading(true);
    try {
      setUsers(await listUsers());
    } catch (e: any) {
      showError(getErrorMessage(e, "Failed to load users"));
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    void load();
  }, []);

  const sorted = useMemo(
    () => [...users].sort((a, b) => a.email.localeCompare(b.email)),
    [users]
  );

  async function onChangeRole(u: User, role: UserRole) {
    setSavingById((m) => ({ ...m, [u.id]: true }));

    try {
      await updateUserRole(u.id, { email: u.email, role });
      showSuccess("User updated");
      await load();
    } catch (e: any) {
      showError(getErrorMessage(e, "Failed to update user"));
    } finally {
      setSavingById((m) => ({ ...m, [u.id]: false }));
    }
  }

  return (
    <div className="mx-auto max-w-4xl p-4">
      <div className="flex items-start justify-between gap-4">
        <div>
          <h1 className="text-xl font-semibold">Users</h1>
          <p className="text-sm text-slate-600">Admin-only. Manage user roles.</p>
        </div>
        <button
          onClick={load}
          className="rounded-md border px-3 py-2 text-sm bg-white hover:bg-slate-50"
          disabled={loading}
        >
          Refresh
        </button>
      </div>

      {loading && (
        <div className="mt-4 rounded-lg border bg-white p-4 text-sm text-slate-600">
          Loading users…
        </div>
      )}

      {!loading && sorted.length === 0 && (
        <div className="mt-4 rounded-lg border bg-white p-6 text-sm text-slate-600">
          No users found.
        </div>
      )}

      {!loading && sorted.length > 0 && (
        <div className="mt-4 overflow-hidden rounded-lg border bg-white">
          <div className="grid grid-cols-12 gap-2 border-b bg-slate-50 px-4 py-2 text-xs font-medium text-slate-600">
            <div className="col-span-7">Email</div>
            <div className="col-span-3">Role</div>
            <div className="col-span-2 text-right">Status</div>
          </div>

          <ul className="divide-y">
            {sorted.map((u) => {
              const saving = !!savingById[u.id];
              const isSelf = sessionUser?.id === u.id;

              return (
                <li
                  key={u.id}
                  className="grid grid-cols-12 items-center gap-2 px-4 py-3 hover:bg-slate-50"
                >
                  <div className="col-span-7">
                    <div className="text-sm font-medium">
                      {u.email}{" "}
                      {isSelf ? (
                        <span className="text-xs text-slate-500">(you)</span>
                      ) : null}
                    </div>
                    <div className="text-xs text-slate-500">{u.id}</div>
                  </div>

                  <div className="col-span-3">
                    <select
                      value={u.role ?? "user"}
                      onChange={(e) => onChangeRole(u, e.target.value)}
                      disabled={saving || isSelf}
                      className="w-full rounded-md border px-2 py-1 text-sm bg-white hover:bg-slate-50 disabled:opacity-60"
                    >
                      {!ROLE_OPTIONS.some((r) => r.value === u.role) && u.role ? (
                        <option value={u.role}>{u.role}</option>
                      ) : null}
                      {ROLE_OPTIONS.map((r) => (
                        <option key={r.value} value={r.value}>
                          {r.label}
                        </option>
                      ))}
                    </select>
                    {isSelf ? (
                      <div className="mt-1 text-xs text-slate-500">
                        Cannot edit your own role.
                      </div>
                    ) : null}
                  </div>

                  <div className="col-span-2 text-right">
                    {saving ? (
                      <span className="text-xs text-slate-500">Saving…</span>
                    ) : (
                      <span className="text-xs text-slate-500">OK</span>
                    )}
                  </div>
                </li>
              );
            })}
          </ul>
        </div>
      )}
    </div>
  );
}
