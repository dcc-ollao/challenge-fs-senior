import type { JSX } from "react";
import { Navigate } from "react-router-dom";
import { useAuth } from "../features/auth/useAuth";

export default function AdminRoute({ children }: { children: JSX.Element }) {
  const { user } = useAuth();


  if (!user) return null;

  if (user.role !== "admin") {
    return <Navigate to="/" replace />;
  }

  return children;
}
