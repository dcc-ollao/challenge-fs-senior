import api from "../../lib/api";

export async function changePassword(currentPassword: string, newPassword: string) {
  await api.post("/auth/change-password", { currentPassword, newPassword });
}
