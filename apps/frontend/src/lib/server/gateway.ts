import { env } from '$env/dynamic/private';

export type GatewayInit = RequestInit & {
  token?: string;
  raw?: boolean;
};

export async function gatewayFetch(path: string, init: GatewayInit = {}) {
  const GATEWAY_URL = env.NOVELHIVE_GATEWAY_URL;
  const INTERNAL_KEY = env.NOVELHIVE_INTERNAL_API_KEY;

  if (!GATEWAY_URL) throw new Error("NOVELHIVE_GATEWAY_URL is not set");
  if (!INTERNAL_KEY) throw new Error("NOVELHIVE_INTERNAL_API_KEY is not set");
 
  const headers = new Headers(init.headers);
  if (!headers.has("Content-Type") && init.body && !(init.body instanceof FormData)) {
    headers.set("Content-Type", "application/json");
  }
  headers.set("X-Internal-Key", INTERNAL_KEY);
  if (init.token) headers.set("Authorization", `Bearer ${init.token}`);

  const url = `${GATEWAY_URL.replace(/\/$/, "")}${path}`;

  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), 20000);
  try {
    const res = await fetch(url, {
      ...init,
      headers,
      signal: controller.signal,
    });
    if (init.raw) return res;
    const text = await res.text();
    let body: any = null;
    try { body = text ? JSON.parse(text) : null; } catch { body = text; }
    if (!res.ok) {
      const msg =
        (body && typeof body === "object" && (body.error || body.message)) ||
        `Request failed: ${res.status}`;
      throw new Error(typeof msg === "string" ? msg : `HTTP ${res.status}`);
    }
    return body;
  } catch (e: any) {
    if (e?.name === "AbortError") throw new Error("Request timed out");
    throw e;
  } finally {
    clearTimeout(timeout);
  }
}
