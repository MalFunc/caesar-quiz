module.exports = {
  content: [
    './src/**/*.{js,ts,jsx,tsx}',
    './src/app/**/*.{js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {
      fontFamily: {
        pixel: ['"Press Start 2P"', 'monospace'],
      },
      colors: {
        pixelBg: '#181825',
        pixelAccent: '#00ffd0',
        pixelYellow: '#ffe066',
        pixelPurple: '#a259f7',
        pixelGreen: '#00ff85',
      },
      animation: {
        'pixel-fade-in': 'pixelFadeIn 0.6s cubic-bezier(.68,-0.55,.27,1.55) both',
        'pixel-bounce': 'pixelBounce 0.7s cubic-bezier(.68,-0.55,.27,1.55) both',
      },
      keyframes: {
        pixelFadeIn: {
          '0%': { opacity: '0', transform: 'scale(0.95) translateY(20px)', filter: 'blur(2px)' },
          '100%': { opacity: '1', transform: 'scale(1) translateY(0)', filter: 'blur(0)' },
        },
        pixelBounce: {
          '0%, 100%': { transform: 'translateY(0)' },
          '50%': { transform: 'translateY(-12px) scale(1.1)' },
        },
      },
    },
  },
  plugins: [],
};
