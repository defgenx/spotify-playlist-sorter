# Setup Guide

Complete setup instructions for the Spotify Playlist Sorter.

## Prerequisites

Before starting, ensure you have:

- **Go 1.21+**: [Download Go](https://golang.org/dl/)
- **Node.js 18+**: [Download Node.js](https://nodejs.org/)
- **ngrok account**: [Sign up for free](https://ngrok.com/)
- **Spotify Developer account**: [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)

Verify installations:
```bash
go version    # Should show go1.21 or higher
node -v       # Should show v18 or higher
npm -v        # Should show 9 or higher
```

---

## Step 1: Clone the Repository

```bash
git clone https://github.com/yourusername/spotify-playlist-sorter.git
cd spotify-playlist-sorter
```

---

## Step 2: Create a Spotify App

1. Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Log in with your Spotify account
3. Click **Create App**
4. Fill in the details:
   - **App name**: Spotify Playlist Sorter (or your choice)
   - **App description**: Organize liked songs by genre
   - **Redirect URIs**: Leave empty for now (we'll add it after ngrok setup)
5. Check the Terms of Service checkbox
6. Click **Save**
7. Note your **Client ID** (visible on the app page)
8. Click **Settings** > **View client secret** and note your **Client Secret**

---

## Step 3: Set Up ngrok

ngrok provides a secure HTTPS tunnel required for Spotify OAuth callbacks.

### Install ngrok

**macOS (Homebrew):**
```bash
brew install ngrok
```

**Linux/Windows:**
Download from [ngrok.com/download](https://ngrok.com/download)

### Configure ngrok

1. Go to [ngrok Dashboard](https://dashboard.ngrok.com/)
2. Sign up or log in
3. Copy your authtoken from the dashboard
4. Configure ngrok:

```bash
ngrok config add-authtoken YOUR_AUTHTOKEN
```

### Start ngrok Tunnel

```bash
ngrok http 3001
```

You'll see output like:
```
Session Status                online
Account                       your@email.com (Plan: Free)
Forwarding                    https://abc123.ngrok-free.app -> http://localhost:3001
```

**Important**: Copy the HTTPS URL (e.g., `https://abc123.ngrok-free.app`). This changes each time you restart ngrok.

---

## Step 4: Configure Spotify Redirect URI

1. Go back to your [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Click on your app
3. Click **Settings**
4. Under **Redirect URIs**, click **Edit**
5. Add: `https://YOUR-NGROK-URL/api/auth/callback`
   - Example: `https://abc123.ngrok-free.app/api/auth/callback`
6. Click **Add**
7. Click **Save**

---

## Step 5: Set Up the Backend

```bash
cd backend

# Copy the example environment file
cp ../.env.example .env

# Install Go dependencies
go mod tidy
```

Edit `.env` with your values:

```bash
# Server Configuration
PORT=3001
FRONTEND_URL=http://localhost:3000
CORS_ORIGINS=http://localhost:3000,https://YOUR-NGROK-URL

# Spotify Configuration
SPOTIFY_CLIENT_ID=your_client_id_here
SPOTIFY_CLIENT_SECRET=your_client_secret_here
SPOTIFY_REDIRECT_URL=https://YOUR-NGROK-URL/api/auth/callback

# Session Configuration
SESSION_SECRET=generate_a_random_string_here
```

**Generating a session secret:**
```bash
# On macOS/Linux
openssl rand -hex 32

# Or use any random string generator
```

---

## Step 6: Set Up the Frontend

```bash
cd frontend

# Install dependencies
npm install
```

No additional configuration needed - the frontend uses Vite's proxy to forward API requests to the backend.

---

## Step 7: Start the Application

You'll need **3 terminal windows**:

### Terminal 1: ngrok (if not already running)
```bash
ngrok http 3001
```

### Terminal 2: Backend
```bash
cd backend
go run ./cmd/server
```

Expected output:
```
{"level":"info","time":"...","message":"Server starting on :3001"}
{"level":"info","time":"...","message":"Environment: development"}
```

### Terminal 3: Frontend
```bash
cd frontend
npm run dev
```

Expected output:
```
  VITE v5.x.x  ready in xxx ms

  ➜  Local:   http://localhost:3000/
  ➜  Network: use --host to expose
```

---

## Step 8: Test the Application

1. Open `http://localhost:3000` in your browser
2. Click **Login with Spotify**
3. Authorize the app on Spotify's page
4. You should be redirected back to the dashboard

If login works, you're all set!

---

## Updating ngrok URL

Each time you restart ngrok, you get a new URL. You need to update:

1. **Spotify Dashboard**: Update the Redirect URI
2. **Backend `.env`**: Update `SPOTIFY_REDIRECT_URL` and `CORS_ORIGINS`
3. **Restart the backend**: `go run ./cmd/server`

**Tip**: ngrok paid plans offer fixed subdomains that don't change.

---

## Environment Variables Reference

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `PORT` | No | `3001` | Backend server port |
| `FRONTEND_URL` | No | `http://localhost:3000` | Frontend URL for redirects |
| `CORS_ORIGINS` | No | `http://localhost:3000` | Comma-separated allowed origins |
| `SPOTIFY_CLIENT_ID` | Yes | - | From Spotify Developer Dashboard |
| `SPOTIFY_CLIENT_SECRET` | Yes | - | From Spotify Developer Dashboard |
| `SPOTIFY_REDIRECT_URL` | Yes | - | ngrok HTTPS URL + `/api/auth/callback` |
| `SESSION_SECRET` | Yes | - | Random string for session encryption |

---

## Common Setup Issues

### "INVALID_CLIENT: Insecure redirect URI"

**Cause**: Spotify requires HTTPS for redirect URIs.

**Solution**: Ensure you're using the ngrok HTTPS URL (not HTTP) in both:
- Spotify Dashboard redirect URI
- `SPOTIFY_REDIRECT_URL` in `.env`

### "Invalid state parameter"

**Cause**: OAuth state mismatch, usually from browser caching or backend restart.

**Solution**:
1. Clear browser cookies for localhost
2. Ensure backend is running
3. Try logging in again

### Session not persisting after login

**Cause**: Cookie not being set correctly.

**Solution**:
1. Check browser DevTools > Application > Cookies
2. Ensure `session` cookie exists for localhost
3. If missing, check CORS configuration includes localhost

### "No songs found"

**Cause**: API permissions or empty library.

**Solution**:
1. Re-authorize the app (logout and login again)
2. Check Spotify account has liked songs
3. Check backend logs for API errors

### ngrok "Invalid Host Header"

**Cause**: ngrok security feature blocking requests.

**Solution**: Use ngrok's `--host-header` flag:
```bash
ngrok http 3001 --host-header=localhost:3001
```

### Port already in use

**Cause**: Another process using port 3001 or 3000.

**Solution**:
```bash
# Find process using port
lsof -i :3001

# Kill it
kill -9 <PID>

# Or change port in .env
```

---

## Development Tips

### Hot Reload

- **Frontend**: Vite provides hot module replacement automatically
- **Backend**: Use `air` for Go hot reload:
  ```bash
  go install github.com/air-verse/air@latest
  air
  ```

### Logging

Backend uses zerolog for structured JSON logging:
```bash
# Pretty print logs (development)
go run ./cmd/server 2>&1 | jq

# Or install pino-pretty
npm install -g pino-pretty
go run ./cmd/server 2>&1 | pino-pretty
```

### Debug Mode

Enable Gin debug mode:
```bash
GIN_MODE=debug go run ./cmd/server
```

---

## Production Considerations

For production deployment:

1. **Use a fixed domain**: Deploy behind a reverse proxy with SSL
2. **Persistent sessions**: Use Redis or database for session storage
3. **Environment**: Set `GIN_MODE=release`
4. **Secrets**: Use secret management (Vault, AWS Secrets Manager)
5. **Rate limiting**: Consider additional rate limiting at reverse proxy level

---

## Quick Reference Commands

```bash
# Start everything (run in separate terminals)
ngrok http 3001                           # Terminal 1
cd backend && go run ./cmd/server         # Terminal 2
cd frontend && npm run dev                # Terminal 3

# Build for production
cd backend && go build -o server ./cmd/server
cd frontend && npm run build

# Run tests
cd backend && go test ./...
cd frontend && npm test

# Format code
cd backend && go fmt ./...
cd frontend && npm run lint
```
