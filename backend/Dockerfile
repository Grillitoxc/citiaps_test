# backend/Dockerfile

# 1) Build
FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Ajusta el main si corresponde (cmd/main.go, etc.)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

# 2) Runtime
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/server /app/server

# Variables por si quieres defaults (compose las sobreescribe)
ENV PORT=4000
ENV MONGODB_URI=mongodb://root:rootpass@mongo:27017/blog?authSource=admin
ENV MONGODB_DB=blog

EXPOSE 4000
USER nonroot:nonroot
ENTRYPOINT ["/app/server"]
