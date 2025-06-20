FROM golang:1.22-alpine as build
WORKDIR /src
COPY . .
RUN go build -o app ./cmd/yourapp

FROM alpine:latest
WORKDIR /app
COPY --from=build /src/app .
EXPOSE 3000
CMD ["./app"]
