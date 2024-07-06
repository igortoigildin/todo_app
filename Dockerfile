FROM golang:1.22.2-alpine AS build

WORKDIR /build

ENV TODO_PORT="7540"
ENV TODO_DBFILE="../scheduler.db"
ENV TODO_PASSWORD="123"

COPY . .

RUN go mod download

RUN go build ./cmd/server

FROM alpine

COPY --from=build /build/server /usr/bin
COPY --from=build /build/web /usr/bin

EXPOSE ${TODO_PORT}

CMD ["./build"]

