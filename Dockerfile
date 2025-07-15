# to use the official golang
FROM golang:1.23.3

# set the directory in the container
WORKDIR /app

LABEL authors="mmahmooda, musabt, malmannai, alimarhoon"
LABEL description="MEOW container for forum project"

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o bin/server .

EXPOSE 443

CMD ["/app/bin/server"]
