FROM golang:1.18.0-alpine3.15 as builder

RUN apk add --update make git

WORKDIR /go/src/github.com/AltMax/unit_service
COPY go.mod .
COPY go.sum .
RUN --mount=type=cache,target=/Users/me/Library/Caches GOPRIVATE=github.com/AltMax go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux make all

FROM golang:1.18.0-alpine3.15
USER nobody
WORKDIR /app/

COPY --chown=nobody:nobody --from=builder /go/src/github.com/AltMax/unit_service/unit_service .
COPY --chown=nobody:nobody --from=builder /go/src/github.com/AltMax/unit_service/migrate_common .
COPY --chown=nobody:nobody --from=builder /go/src/github.com/AltMax/unit_service/docker-entrypoint.sh .

ENTRYPOINT [ "./docker-entrypoint.sh"]
