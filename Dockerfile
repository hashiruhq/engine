FROM golang:1.16 as build
RUN mkdir -p /build/matching-engine
WORKDIR /build/matching-engine/

# Force the go compiler to use modules
ENV GO111MODULE=on

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
# COPY go.sum .

# This is the ‘magic’ step that will download all the dependencies that are specified in
# the go.mod and go.sum file.
# Because of how the layer caching system works in Docker, the  go mod download
# command will _ only_ be re-run when the go.mod or go.sum file change
# (or when we add another docker instruction below this line)
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -a -installsuffix cgo \
  --ldflags "-s -w -X 'github.com/hashiruhq/engine/version.GitCommit=$GIT_COMMIT' -X 'github.com/hashiruhq/engine/version.AppVersion=$APP_VERSION' -X 'github.com/hashiruhq/engine/version.ProductID=$PRODUCT_ID' -X 'github.com/hashiruhq/engine/version.MaxUses=$MAX_USES'" \
  -o /usr/bin/matching_engine

FROM alpine:3.9

RUN apk --no-cache add ca-certificates

COPY --from=build /usr/bin/matching_engine /root/
COPY --from=build /build/matching-engine/.engine.yml /root/

ENV LOG_LEVEL="info"
ENV LOG_FORMAT="json"

EXPOSE 6060
RUN mkdir -p /root/backups
WORKDIR /root/

CMD ["./matching_engine", "server"]
