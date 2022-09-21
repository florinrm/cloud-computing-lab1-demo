FROM golang:1.19-alpine

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o /docker-gs-ping

EXPOSE 80

CMD [ "/docker-gs-ping" ]