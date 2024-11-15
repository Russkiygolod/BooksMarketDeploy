FROM golang:latest AS APP
WORKDIR /app
COPY . /app
RUN go build ./cmd/main.go

FROM ubuntu
WORKDIR /app
COPY --from=APP /app/main .
ENV DB_USER_P = 0 DB_PASSWORD_P = 0 DB_NAME_P = 0
EXPOSE 8081
ENTRYPOINT ["/app/main"]
