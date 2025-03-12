FROM golang:1.23.3-alpine AS builder

ENV CGO_ENABLED=0\
  GOOS=linux\
  GOARCH=amd64

WORKDIR /build

COPY . .

RUN go build -o /build/cmd/restapi/app ./cmd/restapi/main.go

#
FROM gcr.io/distroless/static-debian12

COPY --from=builder /build/cmd/restapi/app /app/main
COPY --from=builder /build/docs /app/docs/
COPY --from=builder /build/.env /app/.env
WORKDIR /app

ENTRYPOINT [ "./main" ]
