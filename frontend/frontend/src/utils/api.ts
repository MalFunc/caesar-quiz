// Helper API terpusat. Semua path relatif ke getApiBase() (dinamis).
import { getApiUrl } from "./config";

export type Question = { id: string; cipher: string };
export type LeaderboardEntry = {
  rank: number;
  player: string;
  time_ms: number;
  player_id: string;
};
export type PlayerLite = { id: string; name: string };

export type WsMessage =
  | { type: "question"; id: string; cipher: string }
  | { type: "leaderboard"; leaderboard: LeaderboardEntry[] }
  | { type: "hint"; hint: string };

async function requestJson<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(getApiUrl(path), init);
  if (!res.ok) {
    let message = `Request gagal (${res.status})`;
    try {
      const data = await res.json();
      if (data?.error) message = data.error;
    } catch {
      // biarkan pesan default
    }
    throw new Error(message);
  }
  return (await res.json()) as T;
}

export function fetchQuestions(gameId: string) {
  return requestJson<Question[]>(`/game/${gameId}/questions`);
}

export function submitAnswer(
  gameId: string,
  data: { player_id: string; question_id: string; answer: string },
) {
  return requestJson<{ correct: boolean; time_ms?: number; duplicate?: boolean }>(
    `/game/${gameId}/answer`,
    {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    },
  );
}

export function fetchHallOfFame(gameId: string) {
  return requestJson<LeaderboardEntry[]>(`/hall-of-fame/${gameId}`);
}

export function createGame(payload: {
  name?: string;
  shift?: number;
  target?: string;
}) {
  return requestJson<{
    game_id: string;
    code: string;
    shift: number;
    host_token: string;
    question_id: string;
    cipher: string;
  }>(`/game`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
}

export function joinGame(code: string, name: string) {
  return requestJson<{ game_id: string; player_id: string }>(`/game/join`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ code, name }),
  });
}

export function startGame(gameId: string, hostToken: string) {
  return requestJson<{ message: string; question_id: string }>(
    `/game/${gameId}/start`,
    { method: "POST", headers: { "X-Host-Token": hostToken } },
  );
}

export function broadcastHint(gameId: string, hostToken: string) {
  return requestJson<{ hint: string }>(`/game/${gameId}/hint`, {
    method: "POST",
    headers: { "X-Host-Token": hostToken },
  });
}

export function fetchPlayers(gameId: string) {
  return requestJson<PlayerLite[]>(`/game/${gameId}/players`);
}
