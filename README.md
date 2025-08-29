# Blog – Go + MongoDB + Nuxt

Stack completo para gestionar publicaciones de un blog:

- **Backend:** Go (Gin) + MongoDB
- **DB:** MongoDB (+ Mongo Express)
- **Frontend:** Nuxt (SSR)
- **Orquestación:** Docker Compose  
- **Seed:** se ejecuta **automáticamente** al levantar Mongo (no necesitas hacer nada extra)

> También se incluyen `.env` para correr en local sin Docker si quisieras, pero el camino recomendado es **Docker**.

---

## 1) Requisitos

- [Docker](https://www.docker.com/) 20+
- [Docker Compose](https://docs.docker.com/compose/) v2+

---

## 2) Estructura del proyecto

```bash
├─ backend/ # API Go (Gin)
│ ├─ controllers/
│ ├─ routes/
│ ├─ services/
│ ├─ models/
│ ├─ main.go
│ └─ .env # opcional para correr en local (no Docker)
├─ frontend/ # Nuxt (SSR)
│ ├─ app/
│ ├─ nuxt.config.ts
│ ├─ Dockerfile
│ └─ .env # opcional para correr en local (no Docker)
├─ db/
│ ├─ init/ # <— seed automático (JS de inicialización)
│ │ └─ seed.js
│ └─ dump/ # opcional: BSON dump (si quieres cargarlo manualmente)
├─ docker-compose.yml
└─ README.md
```

### Semillas (seed) de datos

- La carpeta `db/init/` contiene **scripts de inicialización** (p.ej. `seed.js`) que Mongo ejecuta en el primer arranque.  
- Se cargan **12 posts**, **3–4 autores**, **5–6 tags**, mezcla **publicados/borradores**.  
- No se necesita correr nada manualmente.

---

## 3) Levantar todo con Docker

Desde la raíz del repo:

```bash
docker compose up -d --build
```

Es todo, los puertos son:

- Frontend (Nuxt): http://localhost:3000

- Backend (API Go): http://localhost:4000/api

- Mongo Express: http://localhost:8081

(usuario/clave definidos en docker-compose.yml)

---

## 4) Borrar/Levantar nuevamente
```bash
docker compose down
docker compose up -d
docker compose up -d --build
```

---

## 5) Resetear la base y re-ejecutar el seed
```bash
docker compose down -v
docker compose up -d --build
```

---

## 6) Rutas principales (prefijo /api):

- GET /api/posts – listado con filtros q, tag, published, page, limit, sort

- POST /api/posts – crear post

- GET /api/posts/:id – obtener por id

- PUT /api/posts/:id – actualizar

- DELETE /api/posts/:id – eliminar

- GET /api/posts/metrics/by-tag?limit=10&published=true – top tags (solo publicados)