// Utility for WebSocket connection to backend
export function connectWS(gameId: string, onMessage: (data: any) => void): WebSocket {
  const ws = new WebSocket(`ws://139.59.217.119:8080/ws/${gameId}`);
  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      onMessage(data);
    } catch (e) {
      // ignore
    }
  };
  return ws;
}
