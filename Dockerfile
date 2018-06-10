FROM golang
WORKDIR /go/src/trading_engine
COPY . .
RUN go install ./...
ENTRYPOINT /go/bin/trading_engine

# FROM scratch
# COPY --from=0 /go/bin/trading_engine /trading_engine
# ENTRYPOINT /go/bin/trading_engine
