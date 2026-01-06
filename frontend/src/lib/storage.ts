const TOKEN_KEY = "token";

export function setToken(token?: string | null) {
  if (!token || token === "undefined" || token === "null") {
    localStorage.removeItem(TOKEN_KEY);
    return;
  }
  localStorage.setItem(TOKEN_KEY, token);
}

export function getToken(): string | null {
  const t = localStorage.getItem(TOKEN_KEY);
  if (!t || t === "undefined" || t === "null") return null;
  return t;
}

export function clearToken() {
  localStorage.removeItem(TOKEN_KEY);
}

export const storage = {
  setToken,
  getToken,
  clearToken,
};