
FROM golang:1.24-alpine


LABEL description="Forum App"


RUN apk add --no-cache gcc musl-dev


WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN CGO_ENABLED=1 go build -o main .

EXPOSE 8080

CMD ["./main"]