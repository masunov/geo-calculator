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


docker run -d -t -i -e REDIS_NAMESPACE='staging' \
-e POSTGRES_ENV_POSTGRES_PASSWORD='foo' \
-e POSTGRES_ENV_POSTGRES_USER='bar' \
-e POSTGRES_ENV_DB_NAME='mysite_staging' \
-e POSTGRES_PORT_5432_TCP_ADDR='docker-db-1.hidden.us-east-1.rds.amazonaws.com' \
-e SITE_URL='staging.mysite.com' \
-p 80:80 \
--link redis:redis \
--name container_name dockerhub_id/image_name