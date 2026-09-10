# Panduan Deploy Caesar Quiz

Dua layanan:
- **frontend** (Next.js) → `https://quiz.cyronethic.tech`
- **backend** (Go/Gin, HTTP + WebSocket) → `https://api.cyronethic.tech` (WS: `wss://api.cyronethic.tech`)

Cloudflare ada di depan kedua domain. Cukup satu server origin.

---

## Opsi A — Docker Compose (disarankan)

### 1. Siapkan server

```bash
# Ubuntu/Debian
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER && newgrp docker
docker compose version
```

### 2. Ambil kode & konfigurasi

```bash
cd /root/caesar-quiz      # atau path project kamu
cp .env.example .env
nano .env
```

Isi `.env` untuk domain asli:

```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=ganti-password-kuat
POSTGRES_DB=caesar_quiz
PORT=8080

CORS_ORIGINS=https://quiz.cyronethic.tech

API_BASE=https://api.cyronethic.tech
WS_BASE=wss://api.cyronethic.tech
```

### 3. Jalankan

```bash
docker compose up -d --build
docker compose ps
```

Cek cepat:
```bash
curl -s https://api.cyronethic.tech/healthz   # {"status":"ok"}
curl -s https://quiz.cyronethic.tech/         # HTML
```

Data Postgres tersimpan di volume `db-data` (aman saat rebuild container).

---

## Opsi B — Tanpa Docker (build langsung di server)

### Backend

```bash
cd /root/caesar-quiz
go build -o caesar-quiz ./cmd/server
cp deploy/caesar-quiz.service.example /etc/systemd/system/caesar-quiz.service
nano /etc/systemd/system/caesar-quiz.service   # sesuaikan DATABASE_DSN/CORS_ORIGINS
systemctl daemon-reload && systemctl enable --now caesar-quiz
```

Butuh Postgres lokal (db `caesar_quiz`). Contoh:
```bash
sudo -u postgres createdb caesar_quiz
```

### Frontend

```bash
cd /root/caesar-quiz/frontend/frontend
npm ci
npm run build
cp ../../deploy/caesar-quiz-frontend.service.example /etc/systemd/system/caesar-quiz-frontend.service
nano /etc/systemd/system/caesar-quiz-frontend.service
systemctl daemon-reload && systemctl enable --now caesar-quiz-frontend
```

> JANGAN jalankan `npm run dev` (`next dev`) di server publik — itu yang membuat
> situs bocor path & stack trace. Selalu `npm run build` lalu `npm start`.

---

## 4. Reverse proxy (nginx) + WebSocket

```bash
cp deploy/nginx.conf.example /etc/nginx/sites-available/caesar-quiz
ln -s /etc/nginx/sites-available/caesar-quiz /etc/nginx/sites-enabled/
nginx -t && systemctl reload nginx
```

Konfigurasi ini mem-proxy:
- `quiz.cyronethic.tech` → `127.0.0.1:3000`
- `api.cyronethic.tech` → `127.0.0.1:8080` (dengan header `Upgrade` untuk WebSocket)

### Cloudflare

- DNS `quiz` dan `api` → IP server, mode **Proxied** (orange cloud).
- SSL/TLS mode: **Full** atau **Full (strict)**.
- Origin hanya perlu listen `:80` karena Cloudflare menangani TLS.
- WebSocket aktif default di Cloudflare Cloud (Free). Tidak perlu setting khusus.

Kalau ingin TLS langsung di origin (tanpa Cloudflare), ubah `listen 80` → `listen 443 ssl`
dan pasang sertifikat (mis. `certbot`), lalu set SSL mode Cloudflare ke Full.

---

## 5. Verifikasi

```bash
# Health
curl -s https://api.cyronethic.tech/healthz

# CORS: harus ada Access-Control-Allow-Origin
curl -si -X OPTIONS https://api.cyronethic.tech/game \
  -H "Origin: https://quiz.cyronethic.tech" \
  -H "Access-Control-Request-Method: POST" | grep -i access-control

# Runtime config frontend
curl -s https://quiz.cyronethic.tech/runtime-config.js
# harus: window.__APP_CONFIG__={"apiBase":"https://api.cyronethic.tech","wsBase":"wss://api.cyronethic.tech"};
```

Lalu di browser: buka `https://quiz.cyronethic.tech`, coba **Host Game** → harus
langsung dapat kode room (bukan "gagal fetch").

---

## Troubleshooting

| Gejala | Penyebab | Solusi |
|--------|----------|--------|
| "Failed to fetch" saat host/join | CORS belum cocok | Set `CORS_ORIGINS=https://quiz.cyronethic.tech`, restart backend |
| WebSocket diam / soal tak muncul | `WS_BASE` masih `ws://` atau tanpa TLS | Set `WS_BASE=wss://api.cyronethic.tech`, pastikan nginx kirim header Upgrade |
| Halaman 500 + bocor path `/root/...` | Jalan `next dev` | Jalankan `npm run build && npm start` |
| `curl /runtime-config.js` kosong | Env belum di-set ke container/service | Isi `API_BASE`/`WS_BASE`, restart frontend |
| Host tak bisa mulai game | `host_token` hilang (refresh tanpa localStorage) / game lama | Buat game baru dari halaman Host |

## Update versi baru

```bash
cd /root/caesar-quiz
git pull            # atau salin ulang file
docker compose up -d --build
```
