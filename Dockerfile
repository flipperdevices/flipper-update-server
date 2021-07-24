FROM golang:alpine as builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/app .


FROM alpine

COPY --from=builder /go/bin/app /go/bin/app

EXPOSE 8080
ENTRYPOINT ["/go/bin/app"]