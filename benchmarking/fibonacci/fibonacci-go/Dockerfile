FROM golang AS builder
WORKDIR /fibonacci-go
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/fibonacci-go .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN addgroup -S fibonacci-go && adduser -S fibonacci-go -G fibonacci-go
USER fibonacci-go
WORKDIR /home/fibonacci-go
COPY --from=builder /bin/fibonacci-go ./
ENV PORT 8000
EXPOSE 8000
ENTRYPOINT ["./fibonacci-go"]