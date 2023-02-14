FROM golang:1.19 AS build

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o ./app ./server/main.go

EXPOSE 8080

CMD [ "./app" ]