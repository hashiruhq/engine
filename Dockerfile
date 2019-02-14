FROM golang:1.11.5 as build
RUN mkdir -p /build/matching-engine
WORKDIR /build/matching-engine/

COPY . .
RUN go get -v -d ./...

RUN CGO_ENABLED=0 go build -a -installsuffix cgo --ldflags "-s -w" -o /usr/bin/matching_engine

FROM alpine:3.7

COPY --from=build /usr/bin/matching_engine /root/
COPY --from=build /build/matching-engine/.engine.yml /root/

EXPOSE 6060
RUN mkdir -p /root/backups
WORKDIR /root/

CMD ["./matching_engine", "server"]