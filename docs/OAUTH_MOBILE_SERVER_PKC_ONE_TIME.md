# OAuth untuk Keerja Backend — Server + Mobile (PKCE + Kode Sekali-Pakai)

Dokumen ini menjelaskan cara Keerja backend mengimplementasikan OAuth Google untuk klien web dan mobile, serta bagaimana mengintegrasikan aplikasi Flutter menggunakan PKCE + mekanisme pertukaran kode sekali-pakai. Di dalamnya juga terdapat konfigurasi, endpoint, langkah pengujian catatan keamanan.

## Ringkasan

- Alur server (web): mengikuti Authorization Code Flow standar. Server menggunakan `GOOGLE_CLIENT_SECRET` untuk menukar authorization code menjadi token, lalu membuat JWT aplikasi (APP_JWT) untuk sistem.
- Alur mobile (direkomendasikan): memakai PKCE agar aman. Setelah login, backend akan mengembalikan kode sekali-pakai melalui deep-link, lalu aplikasi mobile menukar kode itu ke endpoint backend untuk mendapatkan APP_JWT. Backend tidak pernah mengirim `GOOGLE_CLIENT_SECRET` ke aplikasi mobile.
- State disimpan di Redis (atau fallback memory). Redirect URI mobile harus di-whitelist melalui `ALLOWED_MOBILE_REDIRECT_URIS`.

Key files (implemented in this repo):
- internal/service/oauth_service.go (business logic, PKCE, exchange, one-time codes)
- internal/service/oauth_state_store.go (Redis / in-memory state storage)
- internal/handler/http/auth_handler.go (HTTP handlers for auth endpoints)
- internal/routes/auth_routes.go (routes)
- internal/config/config.go (configuration / env parsing — supports JSON credentials file)

## Variabel lingkungan yang diperlukan
Tambahkan variabel berikut ke `.env` (pengembangan) atau simpan di secret manager platform (production):

Wajib (minimum):
- GOOGLE_CLIENT_ID — client id OAuth Google (type: Web application)
- GOOGLE_CLIENT_SECRET — simpan rahasia di server (jangan di-embed di aplikasi)
- GOOGLE_REDIRECT_URI — redirect URI server (mis. `http://localhost:8080/api/v1/auth/oauth/google/callback`)

Opsional / direkomendasikan:
- GOOGLE_CREDENTIALS_FILE — path menuju file JSON credentials dari Google (untuk dev lokal; repo sudah mendukung pembacaan di LoadConfig())
- ALLOWED_MOBILE_REDIRECT_URIS — daftar scheme deep-link mobile (dipisah koma), misal `myapp://oauth-callback`
- REDIS_URL (atau REDIS_HOST/REDIS_PORT/REDIS_PASSWORD) — untuk penyimpanan state dan kode sekali-pakai

Contoh `.env` untuk development:

```
GOOGLE_CLIENT_ID=728006377220-...apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=<secret-from-gcp-console>
GOOGLE_REDIRECT_URI=http://localhost:8080/api/v1/auth/oauth/google/callback
ALLOWED_MOBILE_REDIRECT_URIS=myapp://oauth-callback
REDIS_URL=redis://:mypassword@localhost:6379/0
GOOGLE_CREDENTIALS_FILE=./config/google/client_secret_...json   # optional for dev
```

> PENTING: Jangan pernah commit `GOOGLE_CLIENT_SECRET` atau file JSON credentials ke repositori. Selalu gunakan `.gitignore` dan secret manager di production.

---

## Rute / Endpoint

Semua rute berada di bawah `/api/v1`. Endpoint OAuth yang penting:

- GET /api/v1/auth/oauth/google
  - Query params: client (web|mobile), redirect_uri (optional), post_login_redirect_uri (optional, for deep-link), code_challenge (S256), code_challenge_method
  - Response: JSON { auth_url, state, expires_in }

- GET /api/v1/auth/oauth/google/callback
  - Server handles Google callback (exchange code server-side), creates an app JWT, and either:
    - redirects to `post_login_redirect_uri` with a single-use code (safe) or with `#token=<APP_JWT>` as fallback
    - or returns a JSON token response for web flows

- POST /api/v1/auth/oauth/google/exchange
  - Body (mobile PKCE): { code, code_verifier, state, redirect_uri }
  - Server validates state, verifies code_verifier (if required), exchanges for tokens, saves provider and returns app JWT.

- POST /api/v1/auth/oauth/google/exchange-one-time
  - Body: { code }
  - Server consumes the single-use one-time code created at callback -> returns app JWT and expiry

---

## Alur Mobile PKCE + Kode Sekali-Pakai (direkomendasikan)

1. Aplikasi mobile membuat PKCE pair (code_verifier, code_challenge S256).
2. Mobile minta URL otorisasi dari backend:
  GET /api/v1/auth/oauth/google?client=mobile&code_challenge=<challenge>&redirect_uri=myapp://oauth-callback&post_login_redirect_uri=myapp://oauth-callback
  - Backend menyimpan state dan mengembalikan `auth_url` + `state`.
3. Mobile membuka `auth_url` di browser; user login via Google.
4. Google meng-redirect ke `GOOGLE_REDIRECT_URI` (server callback). Server menukar code menggunakan `GOOGLE_CLIENT_SECRET` dan kemudian:
  - membuat kode sekali-pakai dan meng-redirect ke `myapp://oauth-callback?code=<ONE_TIME_CODE>` (direkomendasikan), atau
  - fallback meng-redirect dengan fragment `myapp://oauth-callback#token=<APP_JWT>`.
5. Mobile menerima deep-link dengan kode sekali-pakai atau token langsung.
6. Jika mobile menerima kode sekali-pakai, mobile POST ke `POST /api/v1/auth/oauth/google/exchange-one-time` dengan body `{ code: "<ONE_TIME_CODE>" }` lalu menerima APP_JWT.
7. Simpan APP_JWT dengan aman di device (mis. `flutter_secure_storage`) dan jangan pernah mengirim `GOOGLE_CLIENT_SECRET` ke mobile.

### Catatan
- State disimpan di Redis (TTL 5 menit) sehingga restart server tidak membuat state hilang.
- Kode sekali-pakai disimpan di Redis (atau fallback in-memory) dan hanya bisa dipakai sekali, dengan TTL singkat (default: 2 menit).

---

## Contoh request

1) Memulai PKCE (panggilan ke backend agar mendapat auth URL)

GET /api/v1/auth/oauth/google?client=mobile&code_challenge=<S256>&redirect_uri=myapp://oauth-callback&post_login_redirect_uri=myapp://oauth-callback

Response:
```json
{
  "data": {
    "auth_url": "https://accounts.google.com/..",
    "state": "random-state-token",
    "expires_in": 300
  }
}
```

2) Mobile menerima deep-link `myapp://oauth-callback?code=<ONE_TIME_CODE>` lalu menukar kode:

POST /api/v1/auth/oauth/google/exchange-one-time
Body: { "code": "<ONE_TIME_CODE>" }

Response:
```json
{
  "data": {
    "access_token": "<APP_JWT>",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

3) Alur PKCE (mobile menerima code dan state) lalu POST ke /api/v1/auth/oauth/google/exchange
Body: { "code":"<code>", "code_verifier":"<verifier>", "state":"<state>", "redirect_uri":"myapp://oauth-callback" }

---

## Penyimpanan refresh_token & masa berlaku token di server

- When exchanging for tokens, the server stores `refresh_token` and `token_expiry` in `oauth_providers` DB table (see migrations).
- A refresh helper exists in `internal/service/oauth_service.go` which can be used by background jobs to refresh tokens when expired. Current implementation stores new access token + expiry.

## Cara menguji secara lokal

1. Add credentials to `.env` or `GOOGLE_CREDENTIALS_FILE` — see above.
2. Start Redis (optional but recommended): `docker run -d --name keerja-redis -p 6379:6379 redis:7-alpine`
3. Run your server: `go run ./cmd` or use your usual dev script.
4. Use the sequence described in "Mobile PKCE flow". You can manually test `exchange-one-time` using curl.

## Contoh singkat untuk Flutter (konsep)

Use `flutter_appauth` to handle PKCE and deep-linking. Example steps:

1) Create `code_verifier` and `code_challenge` using S256.
2) Call backend to get `auth_url` and `state` (as described above).
3) Launch auth_url in external browser with `url_launcher`.
4) Listen for the deep-link `myapp://oauth-callback`. If you have `code` + `state` -> `POST /exchange` with `code_verifier` and `state`. If you have one-time code, `POST /exchange-one-time`.
5) Store received `APP_JWT` securely with `flutter_secure_storage`.

Jika Anda ingin, saya bisa menambahkan snippet Flutter yang lengkap — katakan saja dan saya akan tambahkan.

## Catatan keamanan

- Never expose `GOOGLE_CLIENT_SECRET` to the client app.
- Use PKCE for mobile apps and one-time codes for deep-link exchanges to keep tokens secure.
- Use HTTPS in production and add all production redirect URIs to Google Cloud Console.
- Use a secret manager (GCP Secret Manager, etc.) in production for `GOOGLE_CLIENT_SECRET` and other secrets.

---

Jika Anda mau, saya bisa menambahkan contoh Flutter singkat, atau membuat integration test yang mem-mock token endpoint Google untuk menguji alur penuh di CI. Mana yang Anda pilih untuk saya tambahkan berikutnya?
