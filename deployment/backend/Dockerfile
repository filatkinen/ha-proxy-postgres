FROM golang:1.21-alpine as Builder
WORKDIR /usr/src/app/

COPY . /usr/src/app/
COPY ./deployment/backend/backend.env /usr/src/app/serv.env
COPY ./deployment/backend/start.sh /usr/src/app/start.sh

RUN go build -v -o ./bin/backend  ./cmd/backend

#FROM gcr.io/distroless/base-debian12
FROM debian:stable-slim
COPY --from=builder /usr/src/app/bin/backend /app/bin/backend
COPY --from=builder /usr/src/app/serv.env /app/serv.env
COPY --from=builder /usr/src/app/start.sh /app/start.sh
WORKDIR  /app

EXPOSE 8080
ENTRYPOINT [ "./start.sh"]
#ENTRYPOINT [ "./bin/backend"]


