FROM golang:alpine AS build

WORKDIR /build

ENV TODO_PORT="7540"
ENV TODO_DBFILE="../scheduler.db"
ENV TODO_PASSWORD="123"

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build ./cmd/server 

EXPOSE 7540:7540

FROM alpine

COPY --from=build /build/server /usr/bin
COPY --from=build /build/web /usr/bin

CMD ["server"]