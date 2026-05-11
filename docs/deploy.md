# Deploy

## Build Image

```bash
./scripts/build_image.sh
```

Or:

```bash
docker build -t doctor-go:latest .
```

## Run With Docker Compose

```bash
cp .env.deploy.example .env.deploy
docker compose up -d --build
```

The app container runs database migrations before starting the API service.

Manual migration:

```bash
APP_ENV=prod ./doctor-go-migrate -action up
APP_ENV=prod ./doctor-go-migrate -action version
```

## Check Service

```bash
docker compose ps
docker compose logs -f app
curl http://localhost:8080/health
```

## Upgrade

```bash
docker compose build app
docker compose up -d app
```

## Notes

- Replace `JWT_SECRET` before deployment.
- Replace database and Redis passwords before deployment.
- Keep `configs/config.prod.yaml` free of secrets when possible.
- Use environment variables or `.env.deploy` for production secrets.
