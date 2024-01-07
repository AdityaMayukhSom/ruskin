FROM golang:1.21.5-alpine3.19 AS builder

RUN apk update

WORKDIR /usr/src/ruskin

COPY . .

RUN go build -o bin/gstream

FROM scratch 
COPY --from=builder /usr/src/ruskin .
ENTRYPOINT [ "./bin/ruskin" ]