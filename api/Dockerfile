FROM golang:1.18-alpine as build

RUN mkdir /app

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main



FROM alpine:latest

EXPOSE 8080

COPY --from=build /app/main .

CMD ["./main"]