FROM golang:latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/back-end/cmd
RUN go build -o social-network .
EXPOSE 8080
CMD ["./social-network"]
