// Koneksi WebSocket dinamis (otomatis ws:// atau wss:// mengikuti protokol halaman).
import { getWsBase } from "./config";
import type { WsMessage } from "./api";

export function connectWS(
  gameId: string,
  onMessage: (data: WsMessage) => void,
): WebSocket {
  const url = `${getWsBase()}/ws/${gameId}`;
  const ws = new WebSocket(url);
  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data) as WsMessage;
      onMessage(data);
    } catch {
      // ignore pesan non-JSON
    }
  };
  return ws;
}
