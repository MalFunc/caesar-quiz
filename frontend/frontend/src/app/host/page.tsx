"use client";
import { useState, useRef, useEffect } from 'react';
import { connectWS } from '@/utils/ws';
import { API_BASE } from '@/utils/api';

export default function HostPage() {
  const [step, setStep] = useState<'form'|'waiting'|'game'>('form');
  const [name, setName] = useState('');
  const [shift, setShift] = useState(0);
  const [customShift, setCustomShift] = useState('');
  const [target, setTarget] = useState('');
  const [gameId, setGameId] = useState('');
  const [code, setCode] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [question, setQuestion] = useState<string | null>(null);
  const [plaintext, setPlaintext] = useState<string | null>(null);
  const [shiftVal, setShiftVal] = useState<number | null>(null);
  const [hints, setHints] = useState<string[]>([]);
  const [leaderboard, setLeaderboard] = useState<any[]>([]);
  const [players, setPlayers] = useState<{ id: string; name: string }[]>([]);
  const wsRef = useRef<WebSocket | null>(null);

  // Clean up WebSocket on unmount
  useEffect(() => {
    return () => {
      if (wsRef.current) wsRef.current.close();
    };
  }, []);

  // Ambil soal/jawaban/shift dari localStorage saat page load
useEffect(() => {
  const savedQuestion = localStorage.getItem('host_question');
  const savedPlaintext = localStorage.getItem('host_plaintext');
  const savedShift = localStorage.getItem('host_shift');

  if (savedQuestion) setQuestion(savedQuestion);
  if (savedPlaintext) setPlaintext(savedPlaintext);
  if (savedShift) setShiftVal(Number(savedShift));
}, []);


  useEffect(() => {
  if (step === 'game' && gameId) {
    if (wsRef.current) wsRef.current.close();

    wsRef.current = connectWS(gameId, (data: any) => {
      // Abaikan type 'question', host sudah punya data lokal
      if (data.type === 'leaderboard') {
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

  // Broadcast hint handler
  const broadcastHint = async () => {
    if (!gameId) return;
    await fetch(`${API_BASE}/game/${gameId}/hint`, { method: 'POST' });
  };

  // Poll player list every 2s while waiting
  useEffect(() => {
    let interval: NodeJS.Timeout;
    if (step === 'waiting' && gameId) {
      const fetchPlayers = async () => {
        try {
          const res = await fetch(`${API_BASE}/game/${gameId}/players`);
          if (res.ok) {
            setPlayers(await res.json());
          }
        } catch {}
      };
      fetchPlayers();
      interval = setInterval(fetchPlayers, 2000);
    }
    return () => { if (interval) clearInterval(interval); };
  }, [step, gameId]);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    try {
      const payload = {
        name,
        shift: shift === 0 ? 0 : Number(customShift),
        target: target.trim(),
      };
      const res = await fetch(`${API_BASE}/game`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });
      if (!res.ok) throw new Error('Gagal membuat game');
      const data = await res.json();
      setGameId(data.game_id);
      setCode(data.code);
      setStep('waiting');
    } catch (err: any) {
      setError(err.message || 'Gagal membuat game');
    } finally {
      setLoading(false);
    }
  };
// Simpan soal/jawaban/shift ke localStorage saat berubah
useEffect(() => {
  if (question && plaintext && shiftVal !== null) {
    localStorage.setItem('host_question', question);
    localStorage.setItem('host_plaintext', plaintext);
    localStorage.setItem('host_shift', shiftVal.toString());
  }
}, [question, plaintext, shiftVal]);

  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-8 font-pixel">
      <h1 className="text-4xl mb-6 text-black drop-shadow animate-pixel-bounce">Host Game</h1>
      <div className="bg-pixelBg pixel-border p-6 w-full max-w-md animate-pixel-fade-in">
        {step === 'form' && (
          <form onSubmit={handleCreate} className="flex flex-col gap-4">
            <input
              type="text"
              value={name}
              onChange={e => setName(e.target.value)}
              placeholder="Nama Host"
              className="px-4 py-2 text-lg bg-black/60 text-pixelAccent pixel-border outline-none tracking-widest text-center"
              required
            />
            <div className="flex gap-2 items-center">
              <label className="text-black">
                <input type="radio" checked={shift === 0} onChange={() => setShift(0)} /> Random Shift
              </label>
              <label className="text-pixelPurple">
                <input type="radio" checked={shift === 1} onChange={() => setShift(1)} /> Custom Shift
              </label>
              {shift === 1 && (
                <input
                  type="number"
                  min={1}
                  max={25}
                  value={customShift}
                  onChange={e => setCustomShift(e.target.value)}
                  placeholder="Shift"
                  className="w-16 px-2 py-1 text-center bg-black/60 text-pixelPurple pixel-border"
                  required
                />
              )}
            </div>
            <input
              type="text"
              value={target}
              onChange={e => setTarget(e.target.value)}
              placeholder="Kata Target (Soal)"
              className="px-4 py-2 text-lg bg-black/60 text-pixelGreen pixel-border outline-none tracking-widest text-center"
              required
            />
            <button
              type="submit"
              className="bg-pixelAccent text-black px-6 py-2 rounded pixel-border hover:bg-pixelGreen transition-all"
              disabled={loading}
            >
              {loading ? 'Membuat...' : 'Buat Game'}
            </button>
            {error && <div className="text-red-400 text-center">{error}</div>}
          </form>
        )}
        {step === 'waiting' && (
          <div className="flex flex-col items-center gap-4">
            <div className="text-pixelAccent text-lg">Kode Room:</div>
            <div className="text-3xl font-bold text-pixelYellow pixel-border px-6 py-2 tracking-widest">{code}</div>
            <div className="w-full bg-black/40 pixel-border p-2 mt-2">
            <div className="text-pixelGreen text-sm mb-1">
                Peserta yang sudah join ({Array.isArray(players) ? players.length : 0}):
            </div>
            <ul className="list-disc pl-6">
                {(!players || players.length === 0) && <li className="text-gray-400">Belum ada peserta</li>}
                {Array.isArray(players) && players.map((p) => (
                <li key={p.id} className="text-pixelAccent">{p.name}</li>
                ))}
            </ul>
            </div>
            <button
              className="bg-pixelGreen text-black px-6 py-2 rounded pixel-border hover:bg-pixelAccent transition-all mt-4"
              onClick={async () => {
                if (!gameId) return;
                setLoading(true);
                setError('');
                try {
                  const res = await fetch(`${API_BASE}/game/${gameId}/start`, { method: 'POST' });
                  if (!res.ok) throw new Error('Gagal mulai game');
                  setStep('game');
                } catch (err: any) {
                  setError(err.message || 'Gagal mulai game');
                } finally {
                  setLoading(false);
                }
              }}
              disabled={loading}
            >
              {loading ? 'Memulai...' : 'Mulai Game'}
            </button>
            {error && <div className="text-red-400 text-center">{error}</div>}
          </div>
        )}
        {step === 'game' && (
          <div className="flex flex-col gap-6 items-center">
            <div className="w-full bg-black/60 pixel-border p-4">
              <div className="text-pixelAccent text-lg mb-2">Soal (Cipher):</div>
              <div className="text-2xl text-pixelYellow text-center min-h-[2em]">{question || 'Menunggu soal...'}</div>
              {/* Jawaban dan shift */}
              <div className="mt-4 text-pixelGreen text-lg">Jawaban: <span className="text-pixelYellow">{plaintext || '-'}</span></div>
              <div className="text-pixelPurple text-lg">Shift: <span className="text-pixelYellow">{shiftVal ?? '-'}</span></div>
              {/* Area hint */}
              <div className="mt-4">
                <div className="text-pixelAccent mb-2">Hint:</div>
                <ul className="text-pixelGreen">
                  {hints.length === 0 && <li>Belum ada hint</li>}
                  {hints.map((h, i) => <li key={i}>{h}</li>)}
                </ul>
                <button
                  className="mt-2 bg-pixelAccent text-black px-4 py-1 rounded pixel-border hover:bg-pixelGreen transition-all"
                  onClick={broadcastHint}
                >
                  Broadcast Hint
                </button>
              </div>
            </div>
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
