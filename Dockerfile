FROM golang:1.19 as builder

WORKDIR /go/src/app

RUN CGO_ENABLED=0 GOOS=linux go install github.com/jackc/tern@latest

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app main.go

FROM alpine:latest

WORKDIR /opt/bin

COPY --from=builder /go/src/app/deployment/migrations migrations
COPY --from=builder /go/src/app/app app
COPY --from=builder /go/bin/tern /opt/bin/tern
COPY --from=builder /go/src/app/deployment/start.sh start.sh

CMD [ "/opt/bin/start.sh" ]
