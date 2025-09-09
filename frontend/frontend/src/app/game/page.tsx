
"use client";

import { useState, useEffect } from 'react';
import { fetchQuestions, submitAnswer } from '@/utils/api';

type Question = {
  id: string;
  cipher: string;
  shift: number;
  plaintext?: string;
};


export default function GamePage() {
  const [gameId, setGameId] = useState('');
  const [playerId, setPlayerId] = useState('');
  const [questions, setQuestions] = useState<Question[]>([]);
  const [current, setCurrent] = useState(0);
  const [answer, setAnswer] = useState('');
  const [showQuestion, setShowQuestion] = useState(true);
  const [submitted, setSubmitted] = useState(false);
  const [isCorrect, setIsCorrect] = useState<boolean|null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    if (gameId) {
      setLoading(true);
      fetchQuestions(gameId)
        .then((data) => {
          setQuestions(data);
          setCurrent(0);
          setError('');
        })
        .catch(() => setError('Failed to fetch questions'))
        .finally(() => setLoading(false));
    }
  }, [gameId]);

  const handleInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    setAnswer(e.target.value.toUpperCase());
    setIsCorrect(null);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!playerId || !questions[current]) return;
    setSubmitted(true);
    try {
      const res = await submitAnswer({
        player_id: playerId,
        question_id: questions[current].id,
        answer,
      });
      setIsCorrect(res.correct ?? null);
    } catch {
      setIsCorrect(false);
    }
    setTimeout(() => {
      setShowQuestion(false);
      setTimeout(() => {
        setShowQuestion(true);
        setAnswer('');
        setSubmitted(false);
        setIsCorrect(null);
        setCurrent((prev) => (prev + 1 < questions.length ? prev + 1 : prev));
      }, 800);
    }, 1200);
  };

  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-8 font-pixel">
      <h1 className="text-4xl mb-6 text-pixelYellow drop-shadow animate-pixel-bounce">CAESAR QUIZ</h1>
      <div className="w-full max-w-md bg-pixelBg pixel-border p-6 mb-8 animate-pixel-fade-in">
        <form className="flex flex-col gap-4 mb-6" onSubmit={e => { e.preventDefault(); }}>
          <input
            type="text"
            value={gameId}
            onChange={e => setGameId(e.target.value)}
            placeholder="Game ID"
            className="px-4 py-2 text-lg bg-black/60 text-pixelAccent pixel-border outline-none tracking-widest text-center"
            style={{ fontFamily: 'inherit', letterSpacing: '0.2em' }}
          />
          <input
            type="text"
            value={playerId}
            onChange={e => setPlayerId(e.target.value)}
            placeholder="Player ID"
            className="px-4 py-2 text-lg bg-black/60 text-pixelPurple pixel-border outline-none tracking-widest text-center"
            style={{ fontFamily: 'inherit', letterSpacing: '0.2em' }}
          />
        </form>
        {loading && <div className="text-pixelAccent">Loading questions...</div>}
        {error && <div className="text-red-400">{error}</div>}
        {questions.length > 0 && showQuestion && questions[current] && (
          <div>
            <div className="text-pixelAccent text-lg mb-2">Encrypted:</div>
            <div className="text-2xl mb-4 tracking-widest text-pixelGreen">{questions[current].cipher}</div>
            <div className="text-pixelPurple mb-4">Shift: {questions[current].shift}</div>
            <form onSubmit={handleSubmit} className="flex flex-col items-center gap-4">
              <input
                type="text"
                value={answer}
                onChange={handleInput}
                maxLength={16}
                className="px-4 py-2 text-lg bg-black/60 text-pixelYellow pixel-border outline-none focus:ring-2 focus:ring-pixelAccent transition-all animate-pixel-fade-in tracking-widest text-center"
                placeholder="Type your answer..."
                autoFocus
                style={{ fontFamily: 'inherit', letterSpacing: '0.2em' }}
              />
              <button
                type="submit"
                className="bg-pixelAccent text-black px-6 py-2 rounded pixel-border hover:bg-pixelGreen transition-all"
                disabled={submitted}
              >
                Submit
              </button>
            </form>
            {isCorrect !== null && (
              <div className={`mt-4 text-xl font-bold ${isCorrect ? 'text-pixelGreen animate-pixel-bounce' : 'text-red-400 animate-pixel-bounce'}`}>
                {isCorrect ? 'Correct!' : 'Wrong!'}
              </div>
            )}
          </div>
        )}
        {questions.length > 0 && current >= questions.length && (
          <div className="text-pixelAccent text-xl mt-8 animate-pixel-bounce">Quiz Complete!</div>
        )}
      </div>
    </main>
  );
}
