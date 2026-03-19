FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/tline-api ./cmd/server

FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

COPY --from=builder /bin/tline-api /usr/local/bin/tline-api
COPY configs/config.yaml /app/configs/config.yaml

WORKDIR /app
EXPOSE 8080
CMD ["tline-api"]
