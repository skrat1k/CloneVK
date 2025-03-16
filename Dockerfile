FROM golang:1.23.4 as builder

WORKDIR /app

COPY ["go.mod", "go.sum", "./"]

RUN go mod download

COPY . .

RUN go build -o ./bin/app cmd/clonevk/main.go

FROM alpine
COPY --from=builder /app/bin/app /

CMD ["/app"]