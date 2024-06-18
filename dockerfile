FROM golang:1.22.4-alpine

WORKDIR /app

COPY ./src/go.mod ./src
COPY ./src/go.sum ./src
RUN cd src && go mod download

COPY ./src ./src
COPY ./config/dev.config.yml ./config/dev.config.yml

RUN cd src && go build -o health-checker

WORKDIR /app/src

EXPOSE 8080

CMD ["./health-checker", "serve"]

