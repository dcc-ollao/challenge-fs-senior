import { useState } from "react";
import { changePassword, exportAdminDataZip } from "./api";
import { useAuth } from "../auth/useAuth";

export function AccountPage() {
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  
  const { user } = useAuth();
  const isAdmin = user?.role === "admin";


  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setSuccess(null);

    if (newPassword.length < 4) {
      setError("New password must be at least 4 characters.");
      return;
    }
    if (newPassword !== confirm) {
      setError("New password and confirmation do not match.");
      return;
    }

    setIsLoading(true);
    try {
      await changePassword(currentPassword, newPassword);
      setSuccess("Password updated.");
      setCurrentPassword("");
      setNewPassword("");
      setConfirm("");
    } catch (err: any) {
      const msg =
        err?.response?.data?.message ??
        "Failed to update password.";
      setError(msg);
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <div className="max-w-md">
      <h1 className="text-xl font-semibold tracking-tight">Account</h1>
      <p className="mt-1 text-sm text-slate-600">Change your password.</p>

      <form onSubmit={onSubmit} className="mt-6 space-y-4 rounded-lg border bg-white p-4">
        <div>
          <label className="text-sm font-medium">Current password</label>
          <input
            type="password"
            value={currentPassword}
            onChange={(e) => setCurrentPassword(e.target.value)}
            className="mt-1 w-full rounded-md border px-3 py-2 text-sm"
            autoComplete="current-password"
            required
          />
        </div>

        <div>
          <label className="text-sm font-medium">New password</label>
          <input
            type="password"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            className="mt-1 w-full rounded-md border px-3 py-2 text-sm"
            autoComplete="new-password"
            required
          />
        </div>

        <div>
          <label className="text-sm font-medium">Confirm new password</label>
          <input
            type="password"
            value={confirm}
            onChange={(e) => setConfirm(e.target.value)}
            className="mt-1 w-full rounded-md border px-3 py-2 text-sm"
            autoComplete="new-password"
            required
          />
        </div>

        {error && (
          <div className="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">
            {error}
          </div>
        )}
        {success && (
          <div className="rounded-md border border-green-200 bg-green-50 px-3 py-2 text-sm text-green-700">
            {success}
          </div>
        )}

        <button
          type="submit"
          disabled={isLoading}
          className="rounded-md border px-3 py-2 text-sm hover:bg-slate-50 disabled:opacity-60"
        >
          {isLoading ? "Updating..." : "Update password"}
        </button>
      </form>
      {isAdmin && (
        <div className="mt-6 rounded-lg border bg-white p-4">
          <div className="text-sm font-medium">Admin</div>
          <p className="mt-1 text-sm text-slate-600">
            Export all application data as CSV files inside a ZIP.
          </p>

          <div className="mt-3">
            <button
              type="button"
              onClick={exportAdminDataZip}
              className="rounded-md border px-3 py-2 text-sm hover:bg-slate-50"
            >
              Export data (ZIP)
            </button>
          </div>
        </div>
      )}
    </div>
    
  );
}
