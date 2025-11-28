# Docker Setup Guide

This guide explains how to run the Spotify Playlist Sorter using Docker Compose.

## Prerequisites

- Docker Engine 20.10+
- Docker Compose 2.0+
- Spotify Developer account with app credentials

## Quick Start

### 1. Configure Environment Variables

Create a `.env` file in the project root:

```bash
# Spotify OAuth Configuration
SPOTIFY_CLIENT_ID=your_spotify_client_id_here
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret_here
SPOTIFY_REDIRECT_URL=http://localhost:3001/api/auth/callback

# Session Secret (generate a random string)
SESSION_SECRET=your_random_session_secret_here
```

**Important**: Generate a secure random string for `SESSION_SECRET`:
```bash
# On Linux/Mac
openssl rand -hex 32

# Or use an online generator
```

### 2. Configure Spotify App

1. Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Create a new app (or use existing)
3. Add Redirect URI: `http://localhost:3001/api/auth/callback`
4. Copy Client ID and Client Secret to your `.env` file

### 3. Start the Services

```bash
# Build and start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### 4. Access the Application

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:3001

## Development Mode

For development with hot-reload:

```bash
# Use the dev override file
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up
```

This will:
- Mount source code as volumes for live reloading
- Run backend with `go run` instead of compiled binary
- Run frontend with `npm run dev` instead of built static files

## Production Deployment

### Building for Production

```bash
# Build images
docker-compose build

# Start services
docker-compose up -d
```

### Environment Variables for Production

Update your `.env` file with production values:

```bash
# Use HTTPS URLs for production
SPOTIFY_REDIRECT_URL=https://yourdomain.com/api/auth/callback
FRONTEND_URL=https://yourdomain.com
CORS_ORIGINS=https://yourdomain.com

# Use a strong session secret
SESSION_SECRET=<generate-strong-random-secret>
```

### Using HTTPS

For production, you'll need to:

1. Set up a reverse proxy (nginx, Traefik, etc.) in front of the containers
2. Configure SSL certificates (Let's Encrypt recommended)
3. Update `SPOTIFY_REDIRECT_URL` to use HTTPS
4. Update `CORS_ORIGINS` to include your HTTPS domain

Example nginx configuration:

```nginx
server {
    listen 443 ssl;
    server_name yourdomain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /api {
        proxy_pass http://localhost:3001;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## Docker Compose Services

### Backend Service

- **Image**: Built from `backend/Dockerfile`
- **Port**: 3001 (mapped from container port 8080)
- **Health Check**: Checks `/api/health` endpoint
- **Volumes**: 
  - `./backend/certs:/root/certs:ro` (for SSL certificates if needed)

### Frontend Service

- **Image**: Built from `frontend/Dockerfile`
- **Port**: 3000 (mapped from container port 80)
- **Health Check**: Checks root endpoint
- **Depends on**: Backend service

## Troubleshooting

### Services won't start

```bash
# Check logs
docker-compose logs backend
docker-compose logs frontend

# Rebuild images
docker-compose build --no-cache

# Restart services
docker-compose restart
```

### Backend can't connect to Spotify

- Verify `SPOTIFY_CLIENT_ID` and `SPOTIFY_CLIENT_SECRET` are set correctly
- Check that `SPOTIFY_REDIRECT_URL` matches your Spotify app configuration
- Ensure the redirect URL uses the correct protocol (http/https)

### Frontend can't reach backend

- Verify both services are on the same Docker network (`spotify-sorter-network`)
- Check that nginx proxy configuration in `frontend/nginx.conf` points to `backend:3001`
- Ensure backend is healthy: `docker-compose ps`

### Port already in use

If ports 3000 or 3001 are already in use:

```bash
# Change ports in docker-compose.yml
ports:
  - "3002:80"  # Frontend on 3002
  - "3003:8080"  # Backend on 3003
```

Then update `FRONTEND_URL` and `CORS_ORIGINS` accordingly.

### Session not persisting

- Verify `SESSION_SECRET` is set and consistent
- Check that cookies are being set (browser dev tools)
- Ensure `CORS_ORIGINS` includes your frontend URL

## Useful Commands

```bash
# View running containers
docker-compose ps

# View logs for specific service
docker-compose logs -f backend
docker-compose logs -f frontend

# Execute command in container
docker-compose exec backend sh
docker-compose exec frontend sh

# Rebuild specific service
docker-compose build backend
docker-compose build frontend

# Stop and remove containers, networks, volumes
docker-compose down -v

# Clean up everything (including images)
docker-compose down --rmi all -v
```

## Network Architecture

```
┌─────────────────┐     ┌─────────────────┐
│                 │     │                 │
│  Frontend       │────▶│  Backend        │
│  (nginx:80)     │     │  (Go:8080)      │
│  Port 3000      │     │  Port 3001      │
│                 │◀────│                 │
└─────────────────┘     └─────────────────┘
        │                       │
        │    Docker Network     │
        │  spotify-sorter-network│
        │                       │
        └───────────────────────┘
```

## Security Notes

1. **Never commit `.env` file** - It contains sensitive credentials
2. **Use strong SESSION_SECRET** - Generate a random 32+ character string
3. **Use HTTPS in production** - OAuth requires HTTPS for callbacks
4. **Limit CORS_ORIGINS** - Only include trusted domains
5. **Keep images updated** - Regularly rebuild with latest base images

## Next Steps

After Docker setup:
1. Access http://localhost:3000
2. Click "Login with Spotify"
3. Authorize the application
4. Start analyzing and organizing your playlists!

