import api from "../../lib/api";
import type { User } from "./types";
import { normalizeUser } from "./types";

export async function listUsers(): Promise<User[]> {
  const { data } = await api.get("/users");
  const items = Array.isArray(data) ? data : data?.items ?? [];
  return items.map(normalizeUser);
}

export async function updateUserRole(id: string, payload: { email: string; role: string }): Promise<void> {
  await api.put(`/users/${id}`, payload);
}
