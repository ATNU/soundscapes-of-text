FROM golang:latest as builder

WORKDIR /go/src/github.com/mattnolf/polly

COPY . .

RUN go get ./...

RUN CGO_ENABLED=0 GOOS=linux go install ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest 

COPY --from=builder /go/src/github.com/mattnolf/polly/polly bin/polly
