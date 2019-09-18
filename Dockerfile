FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o healthchain . 

EXPOSE 3000

EXPOSE 8080

CMD [ './healthchain' ]
