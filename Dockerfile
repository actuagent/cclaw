# CCLAW 多阶段构建镜像。
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o cclaw ./cmd/cclaw

FROM alpine:latest
WORKDIR /root
COPY --from=builder /app/cclaw /usr/local/bin/cclaw
RUN mkdir -p /root/.cclaw
VOLUME /root/.cclaw

ENTRYPOINT ["cclaw"]
CMD ["version"]
