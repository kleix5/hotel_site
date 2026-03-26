FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server ./cmd/server

FROM alpine:3.20

WORKDIR /app

COPY --from=build /bin/server /app/server
COPY --from=build /app/web /app/web
COPY --from=build /app/db /app/db

EXPOSE 8080

CMD ["/app/server"]
