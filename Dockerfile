FROM golang:1.19.2-alpine3.15 as builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o librebread

FROM alpine:3.15
WORKDIR /app
COPY --from=builder /build/librebread .
EXPOSE 443 80
CMD ["./librebread"]
