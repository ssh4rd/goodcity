# CI/CD Setup Guide — GoodCity

## Обзор

Проект использует **GitHub Actions** для автоматизации двух процессов:

| Workflow | Файл | Триггер |
|---|---|---|
| **CI** — проверка кода | `.github/workflows/ci.yml` | Любой push / PR в `main` |
| **CD** — сборка и деплой | `.github/workflows/cd.yml` | Push в `main` / ручной запуск |

**Схема CD:**
```
push → main
  ├─ build backend image → ghcr.io/…/goodcity-backend:latest
  ├─ build frontend image → ghcr.io/…/goodcity-frontend:latest
  └─ SSH → сервер → docker compose pull && up -d
```

---

## 1. CI — непрерывная интеграция

### Что проверяется

**Backend (Go):**
- `go mod download` — загрузка зависимостей
- `go build ./...` — компиляция
- `go vet ./...` — статический анализ

**Frontend (Node):**
- `npm ci` — воспроизводимая установка зависимостей
- `npm run build` — полная SSR-сборка (client + server bundles)

### Настройка

CI работает из коробки. Никаких секретов не требуется.

Проверить статус: вкладка **Actions** в репозитории GitHub.

---

## 2. CD — непрерывный деплой

### Архитектура

1. GitHub Actions собирает Docker-образы и публикует их в **GitHub Container Registry (ghcr.io)** — бесплатно для публичных и приватных репозиториев.
2. Деплой-шаг подключается к серверу по SSH и перезапускает сервисы.

### 2.1 Подготовка GitHub-репозитория

#### Создать секреты

Перейти: `Settings → Secrets and variables → Actions → New repository secret`

| Секрет | Описание | Пример |
|---|---|---|
| `DEPLOY_HOST` | IP или домен сервера | `185.10.20.30` |
| `DEPLOY_USER` | SSH-пользователь | `deploy` |
| `DEPLOY_KEY` | Приватный SSH-ключ (PEM, весь файл) | `-----BEGIN OPENSSH PRIVATE KEY-----…` |
| `DEPLOY_PORT` | SSH-порт (если не 22) | `2222` |

> **GITHUB_TOKEN** генерируется автоматически — добавлять не нужно.

#### Включить пакеты (если приватный репозиторий)

`Settings → Actions → General → Workflow permissions` → выбрать **Read and write permissions**.

---

### 2.2 Настройка сервера

#### Создать пользователя для деплоя

```bash
# На сервере (под root)
adduser deploy
usermod -aG docker deploy
```

#### Сгенерировать SSH-ключ

```bash
# Локально
ssh-keygen -t ed25519 -C "github-actions-deploy" -f ~/.ssh/goodcity_deploy

# Публичный ключ → на сервер
ssh-copy-id -i ~/.ssh/goodcity_deploy.pub deploy@YOUR_SERVER_IP

# Приватный ключ → в секрет DEPLOY_KEY на GitHub
cat ~/.ssh/goodcity_deploy
```

#### Установить Docker и Docker Compose

```bash
# На сервере
curl -fsSL https://get.docker.com | sh
systemctl enable --now docker
```

#### Создать директорию проекта

```bash
mkdir -p /opt/goodcity
cd /opt/goodcity
```

#### Создать `.env` файл

```bash
# /opt/goodcity/.env
POSTGRES_PASSWORD=<сложный-пароль>
JWT_SECRET=<случайная-строка-32+-символов>
GITHUB_REPOSITORY_OWNER=<ваш-github-username>

# Опционально (если отличаются от дефолтов)
POSTGRES_DB=goodcity
POSTGRES_USER=goodcity
```

Сгенерировать секрет:
```bash
openssl rand -hex 32
```

#### Скопировать docker-compose.prod.yml на сервер

```bash
# Локально
scp docker-compose.prod.yml deploy@YOUR_SERVER_IP:/opt/goodcity/
```

Или добавить в deploy-скрипт в `cd.yml`:
```yaml
- uses: actions/checkout@v4

- name: Copy compose file
  uses: appleboy/scp-action@v0.1.7
  with:
    host: ${{ secrets.DEPLOY_HOST }}
    username: ${{ secrets.DEPLOY_USER }}
    key: ${{ secrets.DEPLOY_KEY }}
    source: docker-compose.prod.yml
    target: /opt/goodcity/
```

---

### 2.3 Первый запуск на сервере

```bash
# На сервере, в /opt/goodcity
export GITHUB_REPOSITORY_OWNER=your-username

# Войти в GHCR (один раз, или генерировать токен в GitHub → Settings → Developer settings → PAT)
echo YOUR_PAT | docker login ghcr.io -u your-username --password-stdin

# Запустить
docker compose -f docker-compose.prod.yml up -d
```

После этого все последующие деплои будут автоматическими.

---

### 2.4 Проверка деплоя

```bash
# Статус контейнеров
docker compose -f docker-compose.prod.yml ps

# Логи бэкенда
docker compose -f docker-compose.prod.yml logs backend --tail=50

# Быстрый health-check
curl http://localhost:8080/api/practices
```

---

## 3. GitHub Container Registry (GHCR)

Образы хранятся по адресам:
```
ghcr.io/<owner>/goodcity-backend:latest
ghcr.io/<owner>/goodcity-backend:<sha>   # конкретный коммит

ghcr.io/<owner>/goodcity-frontend:latest
ghcr.io/<owner>/goodcity-frontend:<sha>
```

### Сделать образы публичными (опционально)

`GitHub → Packages → goodcity-backend → Package settings → Change visibility → Public`

Тогда сервер сможет делать `docker pull` без авторизации.

---

## 4. Nginx + SSL (рекомендуется для продакшена)

Установить certbot и nginx, создать конфиг:

```bash
apt install nginx certbot python3-certbot-nginx -y
certbot --nginx -d yourdomain.ru
```

Конфиг `/etc/nginx/sites-available/goodcity`:

```nginx
server {
    listen 443 ssl;
    server_name yourdomain.ru;

    # SSL сертификаты от certbot
    ssl_certificate /etc/letsencrypt/live/yourdomain.ru/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.ru/privkey.pem;

    # Фронтенд (SSR)
    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # API
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

server {
    listen 80;
    server_name yourdomain.ru;
    return 301 https://$host$request_uri;
}
```

```bash
ln -s /etc/nginx/sites-available/goodcity /etc/nginx/sites-enabled/
nginx -t && systemctl reload nginx
```

---

## 5. Переменные окружения

### Разработка (docker-compose.yml)

| Переменная | Значение по умолчанию |
|---|---|
| `DATABASE_URL` | `postgres://goodcity:goodcity@postgres:5432/goodcity` |
| `JWT_SECRET` | `dev-secret` |
| `PORT` | `8080` |
| `MIGRATIONS_DIR` | `/migrations` |
| `API_URL` | `http://backend:8080` |

### Продакшен (.env на сервере)

| Переменная | Описание |
|---|---|
| `POSTGRES_PASSWORD` | Сложный пароль БД |
| `JWT_SECRET` | Случайная строка 32+ символов |
| `GITHUB_REPOSITORY_OWNER` | GitHub username |
| `POSTGRES_DB` | Имя БД (по умолч. `goodcity`) |
| `POSTGRES_USER` | Пользователь БД (по умолч. `goodcity`) |

---

## 6. Откат (rollback)

Каждый деплой тегируется SHA коммита. Для отката:

```bash
# На сервере
cd /opt/goodcity

# Обновить .env
BACKEND_IMAGE=ghcr.io/<owner>/goodcity-backend:<sha-предыдущего-коммита>
FRONTEND_IMAGE=ghcr.io/<owner>/goodcity-frontend:<sha>

# Или напрямую через docker
docker pull ghcr.io/<owner>/goodcity-backend:<sha>
docker tag ghcr.io/<owner>/goodcity-backend:<sha> ghcr.io/<owner>/goodcity-backend:latest
docker compose -f docker-compose.prod.yml up -d backend
```

---

## 7. Структура файлов CI/CD

```
goodcity/
├── .github/
│   └── workflows/
│       ├── ci.yml               # Проверка кода на каждый push
│       └── cd.yml               # Сборка и деплой при push в main
├── backend/
│   └── Dockerfile               # Многоэтапная сборка Go-приложения
├── frontend/
│   └── Dockerfile               # Многоэтапная сборка Vue SSR
├── docker-compose.yml           # Локальная разработка (volume mounts)
└── docker-compose.prod.yml      # Продакшен (образы из GHCR)
```

---

## 8. Ручной деплой

Если нужно задеплоить без пуша в `main`:

1. Перейти: `GitHub → Actions → CD`
2. Нажать **Run workflow**
3. Выбрать ветку → **Run workflow**
