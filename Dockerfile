# Use an appropriate base image for your service (e.g., golang:latest)
FROM golang:latest

# Set the working directory in the container
WORKDIR /app

# Copy your Go application files to the container
COPY . .

# Build your Go application
RUN go build -o main .

# Expose the port your service listens on
EXPOSE 8080

# Run your application
CMD ["./main"]
