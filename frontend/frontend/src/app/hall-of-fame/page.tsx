

import { useEffect, useState } from 'react';
import { fetchHallOfFame } from '@/utils/api';

type Player = { rank: number; player: string; time_ms: number; player_id: string };

export default function HallOfFamePage() {
  const [players, setPlayers] = useState<Player[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [gameId, setGameId] = useState('');
  const [submitted, setSubmitted] = useState(false);

  useEffect(() => {
    if (!submitted || !gameId) return;
    setLoading(true);
    fetchHallOfFame(gameId)
      .then((data) => {
        setPlayers(data);
        setError('');
      })
      .catch(() => setError('Failed to fetch leaderboard'))
      .finally(() => setLoading(false));
  }, [submitted, gameId]);

  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-8 font-pixel">
      <h1 className="text-4xl mb-6 text-pixelYellow drop-shadow animate-pixel-bounce">HALL OF FAME</h1>
      <div className="w-full max-w-md bg-pixelBg pixel-border p-6 animate-pixel-fade-in">
        {!submitted && (
          <form
            onSubmit={e => {
              e.preventDefault();
              setSubmitted(true);
            }}
            className="flex flex-col gap-4 mb-6"
          >
            <input
              type="text"
              value={gameId}
              onChange={e => setGameId(e.target.value)}
              placeholder="Kode Room/Game ID"
              className="px-4 py-2 text-lg bg-black/60 text-pixelAccent pixel-border outline-none tracking-widest text-center"
              required
            />
            <button
              type="submit"
              className="bg-pixelAccent text-black px-6 py-2 rounded pixel-border hover:bg-pixelGreen transition-all"
            >
              Lihat Hall of Fame Sesi Ini
            </button>
          </form>
        )}
        {submitted && loading && <div className="text-pixelAccent">Loading leaderboard...</div>}
        {submitted && error && <div className="text-red-400">{error}</div>}
        {submitted && players.length > 0 && (
          <table className="w-full text-left">
            <thead>
              <tr className="text-pixelAccent text-lg">
                <th className="pb-2">Rank</th>
                <th className="pb-2">Name</th>
                <th className="pb-2 text-right">Waktu (ms)</th>
              </tr>
            </thead>
            <tbody>
              {players.map((player) => (
                <tr key={player.player_id}>
                  <td>{player.rank}</td>
                  <td>{player.player}</td>
                  <td className="text-right">{player.time_ms}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
                <tr
                  key={player.name}
                  className={
                    i < 3
                      ? `animate-pixel-bounce ${i === 0 ? 'text-pixelYellow' : i === 1 ? 'text-pixelPurple' : 'text-pixelGreen'} font-bold`
                      : 'text-white'
                  }
                >
                  <td className="py-2 pr-4">{i + 1}</td>
                  <td className="py-2 pr-4">{player.name}</td>
                  <td className="py-2 text-right">{player.score}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
        <div className="mt-4 text-xs text-pixelAccent text-center animate-pixel-fade-in">
          <span>Leaderboard auto-update with pixel animation!</span>
        </div>
      </div>
    </main>
  );
}
