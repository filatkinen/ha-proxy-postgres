FROM golang:latest
WORKDIR /usr/src/app/

RUN  git  clone https://github.com/filatkinen/ha-proxy-postgres
WORKDIR /usr/src/app/ha-proxy-postgres
COPY loader /usr/src/app/ha-proxy-postgres/
RUN go build -v -o ./bin/loader  ./cmd/loader

# FROM scratch
# COPY --from=build /usr/src/app/socialnet/bin/socialnet /usr/src/app/socialnet/bin/socialnet
#WORKDIR /usr/src/app/socialnet

ENTRYPOINT [ "./bin/loader"]
# CMD ["-config", "./configs/server.pgsql.yaml"] 
