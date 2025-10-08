
export const BACKEND_BASE_URL = 'http://localhost:8080'; // change to your deployed API base

function getIdToken(): string | null {
  try {
    return localStorage.getItem('idToken');
  } catch {
    return null;
  }
}

export async function apiFetch(path: string, data?: Record<string, unknown>, init?: RequestInit) {
  const token = getIdToken();
  const headers = new Headers(init?.headers || {});
  headers.set('Content-Type', 'application/json');
  if (token) headers.set('Authorization', `Bearer ${token}`);

  const res = await fetch(`${BACKEND_BASE_URL}${path}`, {
    method: data ? 'POST' : 'GET',
    body: data ? JSON.stringify(data) : undefined,
    headers,
    ...init
  });

  if (!res.ok) {
    let msg = `HTTP ${res.status}`;
    try {
      const j = await res.json();
      msg = j.message || msg;
    } catch {}
    throw new Error(msg);
  }
  const text = await res.text();
  try { return JSON.parse(text); } catch { return text; }
}
