FROM golang:1.10.3 as build
RUN mkdir -p /go/src/gitlab.com/around25/products/matching-engine
WORKDIR /go/src/gitlab.com/around25/products/matching-engine/

COPY . .
RUN go get -v -d ./...

RUN CGO_ENABLED=0 go build -a -installsuffix cgo --ldflags "-s -w" -o /usr/bin/matching_engine

FROM alpine:3.7

COPY --from=build /usr/bin/matching_engine /root/
COPY --from=build /go/src/gitlab.com/around25/products/matching-engine/.engine.yml /root/

EXPOSE 6060
RUN mkdir -p /root/backups
WORKDIR /root/

CMD ["./matching_engine", "server"]