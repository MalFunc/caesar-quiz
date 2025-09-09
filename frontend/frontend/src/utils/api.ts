// Utility for API calls to Go backend
export const API_BASE = 'http://139.59.217.119:8080';

export async function fetchQuestions(gameId: string) {
  const res = await fetch(`${API_BASE}/game/${gameId}/questions`);
  if (!res.ok) throw new Error('Failed to fetch questions');
  return res.json();
}

export async function submitAnswer(data: { player_id: string; question_id: string; answer: string }) {
  const res = await fetch(`${API_BASE}/answer`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error('Failed to submit answer');
  return res.json();
}

export async function fetchHallOfFame(gameId: string) {
  const res = await fetch(`${API_BASE}/hall-of-fame/${gameId}`);
  if (!res.ok) throw new Error('Failed to fetch hall of fame');
  return res.json();
}
