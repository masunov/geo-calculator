# builder image
FROM golang:1.26 as builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY .env.example .env
RUN CGO_ENABLED=0 GOOS=linux go build -a -o GeoCalc ./cmd/geocalc

# generate clean, final image for end users
FROM alpine:latest
COPY --from=builder /build/GeoCalc .
COPY --from=builder /build/.env .
EXPOSE ${APP_PORT}
ENTRYPOINT [ "./GeoCalc" ]
