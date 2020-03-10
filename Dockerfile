FROM golang:1.13-alpine as builder 
RUN apk add git
WORKDIR /go/lightpoint/apiserv
COPY . .
RUN go build -o serv

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /root/
COPY --from=builder /go/lightpoint/apiserv/serv .
COPY --from=builder /go/lightpoint/apiserv/config/config.yaml .
CMD ["./serv"]
