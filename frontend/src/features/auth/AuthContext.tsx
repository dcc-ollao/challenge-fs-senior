import { createContext, useContext, useEffect, useState } from "react";
import api from "../../lib/api";
import { storage } from "../../lib/storage";
import type { AuthContextValue, User } from "./types";

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  // bootstrap session
  useEffect(() => {
    const token = storage.getToken();
    if (token) {
      setIsAuthenticated(true);
      // opcional: setUser desde /me en el futuro
    }
  }, []);

  async function login(email: string, password: string) {
    const res = await api.post("/auth/login", { email, password });

    const { token, user } = res.data;

    storage.setToken(token);
    setUser(user ?? null);
    setIsAuthenticated(true);
  }

  function logout() {
    storage.clearToken();
    setUser(null);
    setIsAuthenticated(false);
  }

  return (
    <AuthContext.Provider
      value={{ user, isAuthenticated, login, logout }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuthContext() {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuthContext must be used within AuthProvider");
  }
  return ctx;
}
