# frontend/Dockerfile  (Nuxt 4)
FROM node:20-alpine AS builder
WORKDIR /app

ENV NODE_ENV=production
ENV NUXT_TELEMETRY_DISABLED=1

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

# --- Runtime stage ---
FROM node:20-alpine
WORKDIR /app

# Nitro usa HOST/PORT; PORT se respeta por compatibilidad
ENV NODE_ENV=production
ENV HOST=0.0.0.0
ENV NITRO_HOST=0.0.0.0
ENV PORT=3000
ENV NITRO_PORT=3000
# Para ver logs de Nitro
ENV NITRO_LOG_LEVEL=debug

COPY --from=builder /app/.output /app/.output

EXPOSE 3000
CMD ["node", ".output/server/index.mjs"]

