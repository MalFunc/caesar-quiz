// Endpoint runtime config: dibaca browser sebelum hydration.
// Ubah env API_BASE / WS_BASE tanpa rebuild.
export const dynamic = "force-dynamic";

export function GET() {
  const cfg = {
    apiBase: process.env.API_BASE || process.env.NEXT_PUBLIC_API_BASE || "",
    wsBase: process.env.WS_BASE || process.env.NEXT_PUBLIC_WS_BASE || "",
  };
  const body = `window.__APP_CONFIG__=${JSON.stringify(cfg)};`;
  return new Response(body, {
    headers: {
      "Content-Type": "application/javascript; charset=utf-8",
      "Cache-Control": "no-store, must-revalidate",
    },
  });
}
