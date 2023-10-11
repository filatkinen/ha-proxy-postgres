FROM golang:1.21-alpine as Builder
WORKDIR /usr/src/app/

COPY . /usr/src/app/
COPY ./deployment/client/start.sh /usr/src/app/start.sh

RUN go build -v -o ./bin/client  ./cmd/client

FROM debian:stable-slim
COPY --from=builder /usr/src/app/bin/client /app/bin/client
COPY --from=builder /usr/src/app/start.sh /app/start.sh
WORKDIR  /app

ENTRYPOINT [ "./start.sh"]
#ENTRYPOINT [ "./bin/client"]


