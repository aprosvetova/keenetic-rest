FROM golang:alpine as builder
RUN apk update && apk add build-base

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/app .


FROM alpine

WORKDIR /app

COPY --from=builder /go/bin/app /go/bin/app

EXPOSE 8080
ENTRYPOINT ["/go/bin/app"]