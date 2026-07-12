# Single container: the Go binary serves both the API and the built UI.
# klickops smart-deploy reads the EXPOSE port and ENV keys below to
# suggest the right port and service bindings automatically.

FROM node:24-alpine AS ui
WORKDIR /app/ui
RUN npm install -g pnpm@10
COPY ui/package.json ui/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY ui/ ./
RUN pnpm build

FROM golang:1.26-alpine AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /server ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=backend /server /app/server
COPY --from=ui /app/ui/build /app/ui/build

ENV PORT=8080 \
    UI_DIR=/app/ui/build \
    DATABASE_URL="" \
    S3_ENDPOINT="" \
    S3_REGION="" \
    S3_BUCKET="" \
    S3_ACCESS_KEY="" \
    S3_SECRET_KEY=""

EXPOSE 8080
USER nonroot
ENTRYPOINT ["/app/server"]
