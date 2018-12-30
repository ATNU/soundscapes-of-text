FROM golang:latest as builder

WORKDIR /go/src/webserver

RUN go get github.com/spf13/viper &&        \
    go get github.com/gorilla/mux &&        \
    go get github.com/fsnotify/fsnotify &&  \
    go get github.com/aws/aws-sdk-go/aws/session && \
	go get github.com/aws/aws-sdk-go/service/polly && \
    go get github.com/aws/aws-sdk-go/aws && \
	go get github.com/aws/aws-sdk-go/service/polly/pollyiface

COPY ./server ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

FROM scratch

COPY --from=builder /app ./
COPY ./server/.cfg.json ./

ENTRYPOINT [ "./app" ]
