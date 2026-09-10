# Caesar Quiz (Teka-Teki Kaisar)

Game kuis sandi Caesar realtime: host membuat room, pemain join pakai kode,
jawaban & leaderboard disiarkan lewat WebSocket.

- Backend: Go + Gin + GORM (PostgreSQL)
- Frontend: Next.js 15 (App Router, TypeScript, Tailwind v4)

```
.
├── cmd/server/        # entrypoint backend
├── internal/          # api, config, db, models, services, ws
├── migrations/        # skema SQL (opsional; app pakai AutoMigrate)
├── tests/             # go test
└── frontend/frontend/ # aplikasi Next.js
```

## Konfigurasi (dinamis, tanpa hardcode)

Tidak ada host/IP yang di-hardcode. Semua lewat environment variable.

### Backend

| Env | Default | Keterangan |
|-----|---------|------------|
| `DATABASE_DSN` | `host=localhost user=postgres password=postgres dbname=caesar_quiz port=5432 sslmode=disable` | koneksi Postgres |
| `PORT` | `8080` | port listen |
| `CORS_ORIGINS` | `*` | daftar origin browser dipisah koma; set ke domain frontend produksi |

### Frontend

| Env | Keterangan |
|-----|------------|
| `API_BASE` | URL publik backend, mis. `https://api.example.com` |
| `WS_BASE` | URL WebSocket backend, **wajib `wss://`** di produksi, mis. `wss://api.example.com` |

Nilai `API_BASE`/`WS_BASE` dibaca **runtime** oleh endpoint `/runtime-config.js`,
jadi bisa diganti tanpa rebuild. Fallback: `NEXT_PUBLIC_API_BASE` /
`NEXT_PUBLIC_WS_BASE`, lalu same-origin.

## Menjalankan lokal

Backend (butuh Postgres):

```bash
export DATABASE_DSN="host=localhost user=postgres password=postgres dbname=caesar_quiz port=5432 sslmode=disable"
export CORS_ORIGINS="http://localhost:3000"
go run ./cmd/server
```

Frontend:

```bash
cd frontend/frontend
npm install
API_BASE=http://localhost:8080 WS_BASE=ws://localhost:8080 npm run dev
```

Buka http://localhost:3000

## Deploy via Docker Compose

```bash
cp .env.example .env      # sesuaikan CORS_ORIGINS / API_BASE / WS_BASE
docker compose up -d --build
```

- Frontend: http://localhost:3000
- Backend:  http://localhost:8080
- `GET /healthz` untuk health check.

Untuk produksi di belakang HTTPS, set:

```
CORS_ORIGINS=https://quiz.example.com
API_BASE=https://api.example.com
WS_BASE=wss://api.example.com
```

> Penting: jangan pernah menjalankan `next dev` di server publik. Image
> frontend sudah memakai production build (`next build` + standalone server).

## Test

```bash
# dari root repo (lewati node_modules frontend agar tidak lambat):
go test ./cmd/... ./internal/... ./tests/...
```

