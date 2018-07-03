FROM golang:1.10.3 as build
RUN mkdir -p /go/src/trading_engine
WORKDIR /go/src/trading_engine/

COPY . .
RUN go get -v -d ./...

RUN CGO_ENABLED=0 go build -a -installsuffix cgo --ldflags "-s -w" -o /usr/bin/trading_engine

FROM alpine:3.7

COPY --from=build /usr/bin/trading_engine /root/
COPY --from=build /go/src/trading_engine/.engine.yml /root/

EXPOSE 6060
RUN mkdir -p /root/backups
WORKDIR /root/

CMD ["./trading_engine", "server"]