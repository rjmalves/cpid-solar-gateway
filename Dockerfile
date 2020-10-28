FROM golang:1.15-alpine as builder

WORKDIR /go/src/cpid-solar-gateway
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
 
RUN go build

FROM alpine
WORKDIR /app
COPY --from=builder /go/src/cpid-solar-gateway .

CMD ["./cpid-solar-gateway"]
