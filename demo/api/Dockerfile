FROM golang:1.21-alpine

ENV PORT=4000
ENV ATTACK_URL=http://192.168.2.146
ENV ATTACK_DURATION=2
ENV ATTACK_FREQ=50

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 4000

CMD ["./main"]
