"use client";
import Link from 'next/link';

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-8 font-pixel">
      <h1 className="text-5xl mb-8 text-black drop-shadow animate-pixel-bounce text-center">TEKA-TEKI KAISAR</h1>
      <div className="flex flex-col gap-6 w-full max-w-xs animate-pixel-fade-in">
        <Link href="/host" className="px-6 py-4 bg-pixelAccent text-black text-3xl pixel-border text-center hover:bg-pixelGreen transition-all">Host Game</Link>
        <Link href="/join" className="px-6 py-4 bg-pixelPurple text-white text-3xl pixel-border text-center hover:bg-pixelYellow hover:text-black transition-all">Join Game</Link>
      </div>
      <div className="mt-12 text-xs text-pixelAccent text-center animate-pixel-fade-in">Cyronrthic<br/>by Atmint SKS</div>
    </main>
  );
}
