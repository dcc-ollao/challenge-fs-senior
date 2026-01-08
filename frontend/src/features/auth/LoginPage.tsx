import { useState } from "react";
import { useAuth } from "./useAuth";
import { useNavigate } from "react-router-dom";
import { useSnackbar } from "../../components/snackbar/SnackbarContext";

type Mode = "login" | "register";

export default function LoginPage() {
  const { login, register, isLoading } = useAuth();
  const navigate = useNavigate();
  const { showError, showSuccess } = useSnackbar();

  const [mode, setMode] = useState<Mode>("login");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  function isValidEmail(value: string) {
    return /\S+@\S+\.\S+/.test(value);
  }
  
  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();

    if (!isValidEmail(email)) {
    showError("Please enter a valid email address");
    return;
  }

    try {
      if (mode === "login") {
        await login(email.trim(), password);
        navigate("/");
      } else {
        await register(email.trim(), password);
        showSuccess("Account created");
        navigate("/");
      }
    } catch (err: any) {
      showError(
        err?.response?.data?.message ??
          err?.response?.data?.error ??
          err?.message ??
          (mode === "login" ? "Login failed" : "Register failed")
      );
    }
  }

  return (
    <div className="p-6 max-w-sm mx-auto">
      <h1 className="text-xl font-semibold mb-1">
        {mode === "login" ? "Login" : "Create account"}
      </h1>

      <form onSubmit={handleSubmit} className="space-y-3">
        <input
          type="email"
          className="border p-2 w-full rounded-md"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          disabled={isLoading}
          required
        />

        <input
          className="border p-2 w-full rounded-md"
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          disabled={isLoading}
        />

        <button
          className="bg-black text-white px-4 py-2 w-full rounded-md disabled:opacity-60"
          disabled={isLoading || !email.trim() || !password}
          type="submit"
        >
          {isLoading ? "Loading…" : mode === "login" ? "Login" : "Create account"}
        </button>

        <button
          type="button"
          className="w-full text-sm text-slate-600 hover:text-slate-900"
          onClick={() => setMode((m) => (m === "login" ? "register" : "login"))}
          disabled={isLoading}
        >
          {mode === "login" ? "No account? Create one" : "Already have an account? Login"}
        </button>
      </form>
    </div>
  );
}
