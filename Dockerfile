FROM golang:1.13.2-alpine3.10 as builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o librebread

FROM alpine:3.10
WORKDIR /app
COPY --from=builder /build/librebread .
EXPOSE 443
CMD ["./librebread"]
