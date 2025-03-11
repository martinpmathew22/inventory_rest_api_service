# syntax=docker/dockerfile:1

FROM golang:1.21

# Set destination for COPY
WORKDIR /app

# Download Go modules
ENV GOPRIVATE="example.com/inventory_rest_api_service"
ENV GONOSUMDB="example.com/inventory_rest_api_service"
COPY . ./
# COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
# COPY models ./



# RUN go mod tidy
# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /inventory_service

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8000

# Run
CMD ["/inventory_service"]
