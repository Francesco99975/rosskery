FROM golang:1.23.1-alpine3.19 AS build

RUN apk --no-cache add gcc g++ make git

WORKDIR /go/src/app

COPY . .

RUN go mod tidy

RUN mv .prod.env .env

RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/rosskery ./cmd/server/*.go

FROM alpine:3.19

RUN apk update && apk upgrade && apk --no-cache add ca-certificates

WORKDIR /go/bin

COPY --from=build /go/src/app/bin /go/bin
COPY --from=build /go/src/app/.env /go/bin/
COPY --from=build /go/src/app/sql /go/bin/sql
COPY --from=build /go/src/app/static /go/bin/static

EXPOSE 8078

ENTRYPOINT /go/bin/rosskery --port 8078
