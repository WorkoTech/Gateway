# build stage
FROM golang as builder

ENV GO111MODULE=on

WORKDIR /app

RUN go get github.com/githubnemo/CompileDaemon

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 make build

CMD /app/bin/gateway
