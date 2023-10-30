# Use an official Go runtime as a parent image
FROM golang:1.21.3

# Set the working directory to /dist
WORKDIR /dist

# Copy the current directory contents into the container at /dist
COPY . ./

# Download Go Modules
RUN go mod download

# # Copy all files
# COPY *.go ./

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /chat-be

# Expose port 5050
EXPOSE 5051

# Run the application
CMD ["/chat-be"]
