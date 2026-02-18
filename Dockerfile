FROM golang:1.24-alpine

WORKDIR /app

RUN apk add --no-cache git postgresql-client

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install goa.design/goa/v3/cmd/goa@latest
RUN goa gen github.com/pgquerynarrative/pgquerynarrative/api/design
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/server ./cmd/server

ENV PGQUERYNARRATIVE_HOST=0.0.0.0
ENV PGQUERYNARRATIVE_PORT=8080

COPY tools/docker/entrypoint.sh /app/tools/docker/entrypoint.sh
RUN chmod +x /app/tools/docker/entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/app/tools/docker/entrypoint.sh"]
