FROM golang:1.13.2-alpine3.10 as builder
WORKDIR /build
COPY . .
RUN go build -o librebread

FROM alpine:3.10
WORKDIR /app
COPY --from=builder /build/librebread .

CMD ["./librebread"]
