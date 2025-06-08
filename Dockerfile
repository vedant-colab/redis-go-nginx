FROM golang:1.23-alpine AS builder
WORKDIR /
COPY go.mod go.sum ./
RUN go mod download


COPY . .
RUN go build -o main .

FROM alpine:latest
WORKDIR /
COPY --from=builder /main .
CMD [ "./main" ]
