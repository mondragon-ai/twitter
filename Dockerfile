FROM golang:1.22.0

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
RUN go build -o ./out/dist .

# Expose the port your application listens on
EXPOSE 8080

CMD ./out/dist
