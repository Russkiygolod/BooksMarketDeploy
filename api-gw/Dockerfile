FROM golang:latest AS APP
WORKDIR /app
COPY . /app
RUN go build ./cmd/main.go

FROM ubuntu
WORKDIR /app
COPY --from=APP /app/main .
EXPOSE 80
ENTRYPOINT ["/app/main"]
