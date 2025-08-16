FROM node:24-alpine AS frontend

RUN npm install pnpm@10 -g

WORKDIR /frontend
COPY web/package.json web/pnpm-lock.yaml ./
RUN pnpm install
COPY web/ ./
RUN pnpm run build

# Stage 2: Build Go Backend
FROM golang:1.24-alpine AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /frontend/dist ./web/dist
RUN go build -o main .

# Stage 3: Final production image
FROM alpine:3.21
WORKDIR /app
COPY --from=backend /app/main .

EXPOSE 8080
CMD ["/app/main"]

