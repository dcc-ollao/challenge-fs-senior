export type UserRole = "admin" | "user" | string;

export type User = {
  id: string;
  email: string;
  role: UserRole;
  created_at?: string;
};

export function normalizeUser(raw: any): User {
  return {
    id: raw?.id ?? raw?.ID ?? raw?.Id,
    email: raw?.email ?? raw?.Email,
    role: raw?.role ?? raw?.Role,
    created_at: raw?.created_at ?? raw?.createdAt,
  };
}
