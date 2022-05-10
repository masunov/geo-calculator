# builder image
FROM golang:1.17 as builder
RUN mkdir /build
ADD *.go /build/
ADD go.mod /build/
ADD go.sum /build/
COPY .env.example /build/.env
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o GeoCalc .


# generate clean, final image for end users
FROM alpine:latest
COPY --from=builder /build/GeoCalc .
COPY --from=builder /build/.env .

EXPOSE ${APP_PORT}

ENTRYPOINT [ "./GeoCalc" ]