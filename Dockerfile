FROM golang:1.19 as builder

WORKDIR /app

COPY . .

ENV GOPROXY https://goproxy.cn,direct

RUN go mod tidy && CGO_ENABLED=0 go build -o danmaku .

FROM alpine

WORKDIR /app

ENV GIN_MODE release
ENV LOG_LEVEL error

COPY --from=builder /app/danmaku ./danmaku

#RUN chmod +x ./danmaku

EXPOSE 8080

CMD ["./danmaku", " -b=true -p=8080"]
