# Use the official Golang image
FROM golang:1.23

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./
COPY go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the main application
RUN go build -o main ./cmd/main.go

# Build the seed application
RUN go build -o seed ./script/seed.go

# Command to run the executable
CMD ["./main"]
