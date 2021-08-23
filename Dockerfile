FROM golang:alpine as builder
COPY go.mod go.sum /go/src/github.com/mcarreon/golang-crud-api/
WORKDIR /go/src/github.com/mcarreon/golang-crud-api
RUN go mod download
COPY . /go/src/github.com/mcarreon/golang-crud-api
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/golang-crud-api github.com/mcarreon/golang-crud-api/


FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/mcarreon/golang-crud-api/golang-crud-api.exe /usr/bin/golang-crud-api
EXPOSE 3000 3000
RUN chmod a+x /usr/bin/golang-crud-api
ENTRYPOINT ["/usr/bin/golang-crud-api"]