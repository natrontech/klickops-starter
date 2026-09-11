# Single container: the Go binary serves both the API and the built UI.
# klickops smart-deploy reads the EXPOSE port and ENV keys below to
# suggest the right port and service bindings automatically.

FROM node:26-alpine AS ui
WORKDIR /app/ui
RUN npm install -g pnpm@10
COPY ui/package.json ui/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY ui/ ./
RUN pnpm build

FROM golang:1.27-alpine AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /server ./cmd/server

FROM gcr.io/distroless/static-debian13:nonroot
WORKDIR /app

# Stamped at build time (`docker build --build-arg GIT_SHA=$(git rev-parse --short HEAD)`)
# and shown in the UI footer, so a preview environment says which commit it runs.
ARG GIT_SHA=dev
COPY --from=backend /server /app/server
COPY --from=ui /app/ui/build /app/ui/build

ENV PORT=8080 \
    UI_DIR=/app/ui/build \
    APP_ENV=production \
    GIT_SHA=${GIT_SHA} \
    DATABASE_URL="" \
    AWS_ENDPOINT_URL_S3="" \
    AWS_REGION="" \
    S3_BUCKET="" \
    AWS_ACCESS_KEY_ID="" \
    AWS_SECRET_ACCESS_KEY=""

EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/app/server"]
