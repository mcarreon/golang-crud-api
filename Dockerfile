FROM golang:alpine as builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test server_test.go asserts.go models.go server.go utils.go 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/main .

EXPOSE 8080

FROM scratch

COPY --from=builder /app/bin/main .

CMD ["./main"]