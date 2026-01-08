import { createContext, useContext, useEffect, useState } from "react";
import api from "../../lib/api";
import { storage } from "../../lib/storage";
import type { AuthContextValue, AuthUser } from "./types";

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState(() => Boolean(storage.getToken()));
  const [isLoading, setIsLoading] = useState(true);

  async function fetchMe() {
    const res = await api.get("/auth/me");
    const me = res.data;

    setUser({
      id: me.id ?? me.userId ?? me.ID,
      email: me.email,
      role: me.role,
    });
  }

  async function setSessionFromToken(token: string) {
    storage.setToken(token);
    setIsAuthenticated(true);
    await fetchMe();
  }

  useEffect(() => {
    const token = storage.getToken();

    if (!token) {
      setIsAuthenticated(false);
      setUser(null);
      setIsLoading(false);
      return;
    }

    setIsAuthenticated(true);

    fetchMe()
      .catch(() => {
        storage.clearToken();
        setUser(null);
        setIsAuthenticated(false);
      })
      .finally(() => setIsLoading(false));
  }, []);

  async function login(email: string, password: string) {
    setIsLoading(true);
    try {
      const res = await api.post("/auth/login", { email, password });
      await setSessionFromToken(res.data.accessToken);
    } finally {
      setIsLoading(false);
    }
  }

  async function register(email: string, password: string) {
    setIsLoading(true);
    try {
      const res = await api.post("/auth/register", { email, password });
      await setSessionFromToken(res.data.accessToken);
    } finally {
      setIsLoading(false);
    }
  }

  function logout() {
    storage.clearToken();
    setUser(null);
    setIsAuthenticated(false);
    setIsLoading(false);
  }

  return (
    <AuthContext.Provider value={{ user, isAuthenticated, isLoading, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuthContext() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuthContext must be used within AuthProvider");
  return ctx;
}
