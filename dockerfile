FROM golang:1.21.5-bullseye AS build

RUN apt-get update

WORKDIR /app

COPY . .

RUN go mod download

WORKDIR /app/cmd

RUN go build -o apigateway

FROM busybox:latest

WORKDIR /apigateway/cmd

COPY --from=build /app/cmd/apigateway .

COPY --from=build /app/.env /apigateway

EXPOSE 4000

CMD ["./apigateway"]