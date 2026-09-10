// Konfigurasi runtime agar host API/WS tidak di-hardcode.
// Prioritas: window.__APP_CONFIG__ (di-inject runtime-config.js) -> NEXT_PUBLIC_* -> same-origin.
type AppConfig = { apiBase?: string; wsBase?: string };

declare global {
  interface Window {
    __APP_CONFIG__?: AppConfig;
  }
}

function runtimeConfig(): AppConfig {
  if (typeof window === "undefined") return {};
  return window.__APP_CONFIG__ ?? {};
}

function pick(...values: Array<string | undefined>): string {
  for (const v of values) {
    if (v && v.trim()) return v.trim().replace(/\/+$/, "");
  }
  return "";
}

export function getApiBase(): string {
  const base = pick(
    runtimeConfig().apiBase,
    process.env.NEXT_PUBLIC_API_BASE,
  );
  if (base) return base;
  // Fallback: same-origin (cocok bila backend di-proxy di domain yang sama).
  if (typeof window !== "undefined") return window.location.origin;
  return "";
}

export function getWsBase(): string {
  const base = pick(
    runtimeConfig().wsBase,
    process.env.NEXT_PUBLIC_WS_BASE,
  );
  if (base) return base;
  if (typeof window !== "undefined") {
    const proto = window.location.protocol === "https:" ? "wss:" : "ws:";
    return `${proto}//${window.location.host}`;
  }
  return "";
}

export function getApiUrl(path: string): string {
  const p = path.startsWith("/") ? path : `/${path}`;
  return `${getApiBase()}${p}`;
}
