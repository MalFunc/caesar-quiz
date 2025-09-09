"use client";
import { useState, useRef, useEffect } from 'react';
import { connectWS } from '@/utils/ws';
import { API_BASE } from '@/utils/api';

export default function JoinPage() {
  const [step, setStep] = useState<'form'|'waiting'|'game'>('form');
  const [name, setName] = useState('');
  const [code, setCode] = useState('');
  const [gameId, setGameId] = useState('');
  const [playerId, setPlayerId] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [question, setQuestion] = useState<string | null>(null);
  const [questionId, setQuestionId] = useState<string | null>(null);
  const [hints, setHints] = useState<string[]>([]);
  const [leaderboard, setLeaderboard] = useState<any[]>([]);
  const wsRef = useRef<WebSocket | null>(null);

  const handleJoin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    try {
      const res = await fetch(`${API_BASE}/game/join`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ code, name }),
      });
      if (!res.ok) throw new Error('Gagal join game');
      const data = await res.json();
      setGameId(data.game_id);
      setPlayerId(data.player_id);
      setStep('waiting');
    } catch (err: any) {
      setError(err.message || 'Gagal join game');
    } finally {
      setLoading(false);
    }
  };

  // Connect WebSocket when gameId is set (after join)
  useEffect(() => {
    if ((step === 'waiting' || step === 'game') && gameId) {
      if (wsRef.current) wsRef.current.close();
      wsRef.current = connectWS(gameId, (data: any) => {
        if (data.type === 'question') {
            setQuestion(data.cipher || '');
            setQuestionId(data.id || null);
          
          setStep('game');
        } else if (data.type === 'leaderboard') {
          setLeaderboard(data.leaderboard);
        } else if (data.type === 'hint') {
          setHints(hs => [...hs, data.hint]);
        }
      });
    }
    return () => {
      if (wsRef.current) wsRef.current.close();
    };
  }, [step, gameId]);

  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-8 font-pixel">
      <h1 className="text-4xl mb-6 text-black drop-shadow animate-pixel-bounce">Join Game</h1>
      <div className="bg-pixelBg pixel-border p-6 w-full max-w-md animate-pixel-fade-in">
        {step === 'form' && (
          <form onSubmit={handleJoin} className="flex flex-col gap-4">
            <input
              type="text"
              value={code}
              onChange={e => setCode(e.target.value.toUpperCase())}
              placeholder="Kode Room"
              className="px-4 py-2 text-lg bg-black/60 text-pixelAccent pixel-border outline-none tracking-widest text-center"
              required
            />
            <input
              type="text"
              value={name}
              onChange={e => setName(e.target.value)}
              placeholder="Nama Kamu"
              className="px-4 py-2 text-lg bg-black/60 text-pixelPurple pixel-border outline-none tracking-widest text-center"
              required
            />
            <button
              type="submit"
              className="bg-pixelAccent text-black px-6 py-2 rounded pixel-border hover:bg-pixelGreen transition-all"
              disabled={loading}
            >
              {loading ? 'Gabung...' : 'Gabung Game'}
            </button>
            {error && <div className="text-red-400 text-center">{error}</div>}
          </form>
        )}
        {step === 'waiting' && (
          <div className="flex flex-col items-center gap-4">
            <div className="text-pixelAccent text-lg">Menunggu host mulai game...</div>
            <div className="text-3xl font-bold text-pixelYellow pixel-border px-6 py-2 tracking-widest">{code}</div>
          </div>
        )}
        {step === 'game' && (
          <div className="flex flex-col gap-6 items-center">
            <div className="w-full bg-black/60 pixel-border p-4">
              <div className="text-pixelAccent text-lg mb-2">Soal:</div>
              <div className="text-2xl text-pixelYellow text-center min-h-[2em]">{question || 'Menunggu soal...'}</div>
            </div>
            {/* Hints area */}
            <div className="w-full bg-black/60 pixel-border p-4 mb-2">
              <div className="text-pixelAccent mb-2">Hint:</div>
              <ul className="text-pixelGreen">
                {hints.length === 0 && <li>Belum ada hint</li>}
                {hints.map((h, i) => <li key={i}>{h}</li>)}
              </ul>
            </div>
            <form
              className="flex flex-col gap-2 w-full bg-black/60 pixel-border p-4"
              onSubmit={async (e) => {
                e.preventDefault();
                const form = e.target as HTMLFormElement;
                const answerInput = form.elements.namedItem('answer') as HTMLInputElement;
                const answer = answerInput.value.trim();
                if (!answer || !gameId || !playerId || !questionId) return;

                try {
                  const res = await fetch(`${API_BASE}/game/${gameId}/answer`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ player_id: playerId, question_id: questionId, answer }),
                  });

                  const data = await res.json();

                  // Feedback jawaban
                  if (data.correct === true) {
                    setError('✅ Jawaban Benar!');
                  } else {
                    setError('❌ Jawaban Salah!');
                  }

                  // Reset input
                  answerInput.value = '';
                } catch (err: any) {
                  setError(err.message || 'Terjadi kesalahan jaringan');
                }
              }}
            >
              <input
                type="text"
                name="answer"
                placeholder="Jawaban kamu..."
                className="px-4 py-2 text-lg bg-black/60 text-pixelGreen pixel-border outline-none tracking-widest text-center"
                autoComplete="off"
                required
              />
              <button
                type="submit"
                className="bg-pixelAccent text-black px-6 py-2 rounded pixel-border hover:bg-pixelGreen transition-all"
              >
                Kirim Jawaban
              </button>
            </form>

           {error && (
            <div
              className={`
                text-center mt-4
                text-2xl font-bold
                px-4 py-2
                ${error.includes('Benar') ? 'text-pixelGreen bg-black' : 'text-red-500 bg-black'}
                border-2 border-black
                rounded
                drop-shadow
              `}
            >
              {error}
            </div>
          )}

            <div className="w-full bg-black/60 pixel-border p-4">
              <div className="text-pixelGreen text-lg mb-2">Leaderboard:</div>
              <ol className="list-decimal pl-6">
                {leaderboard.length === 0 && <li className="text-gray-400">Belum ada skor</li>}
                {leaderboard.map((p, i) => (
                <li key={p.player_id || i} className="text-pixelPurple">
                  {p.player}: <span className="text-pixelYellow">{p.time_ms} detik</span>
                </li>
              ))}
              </ol>
            </div>
          </div>
        )}
      </div>
    </main>
  );
}
