FROM node:22-alpine AS frontend-build
WORKDIR /frontend
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm install --silent
COPY frontend/ .
RUN npm run build

FROM golang:1.24-alpine

WORKDIR /app

RUN apk add --no-cache git postgresql-client wget

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend-build /frontend/dist ./frontend/dist

RUN go install goa.design/goa/v3/cmd/goa@latest
RUN goa gen github.com/pgquerynarrative/pgquerynarrative/api/design
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/server ./cmd/server

ENV PGQUERYNARRATIVE_HOST=0.0.0.0
ENV PGQUERYNARRATIVE_PORT=8080

COPY tools/docker/entrypoint.sh /app/tools/docker/entrypoint.sh
RUN chmod +x /app/tools/docker/entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/app/tools/docker/entrypoint.sh"]
