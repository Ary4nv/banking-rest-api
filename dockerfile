FROM golang:1.25.6

WORKDIR /app 

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o api

EXPOSE 3000

CMD ["./api"]