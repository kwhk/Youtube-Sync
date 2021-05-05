FROM arm32v7/golang:alpine
ENV GOOS=linux \
    GOARCH=arm

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o ./build ./cmd/server

CMD ["./build/server"]

