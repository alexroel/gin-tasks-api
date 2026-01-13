# ============================================
# STAGE 1: Build
# ============================================
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copiar dependencias primero (mejor cache)
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar binario estático
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/main.go

# ============================================
# STAGE 2: Runtime (imagen mínima)
# ============================================
FROM alpine:3.19

# Agregar certificados SSL y timezone
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copiar solo el binario compilado
COPY --from=builder /app/main .

# Usuario no-root para seguridad
RUN adduser -D -g '' appuser
USER appuser

EXPOSE 8080

CMD ["./main"]

