FROM golang

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -v -o main

CMD [ "main" ]