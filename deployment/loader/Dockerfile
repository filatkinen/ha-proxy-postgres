FROM golang:1.21-alpine
WORKDIR /usr/src/app/

COPY . /usr/src/app/
COPY ./deployment/loader/loader.env /usr/src/app/loader.env
COPY ./deployment/loader/start.sh /usr/src/app/start.sh

RUN go build -v -o ./bin/loader  ./cmd/loader


ENTRYPOINT [ "./start.sh"]

