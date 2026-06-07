FROM golang:1.26-alpine

WORKDIR /app

# 1. Copy over your dependency configurations and download them
COPY go.mod go.sum ./
RUN go mod download

# 2. Copy the rest of your app codebase (cmd/, web/, content/)
COPY . .

# 3. Build the Go binary directly inside the container footprint
RUN go build -o /server ./cmd/*.go

EXPOSE 8080

# 4. Boot up the freshly compiled application
CMD ["/server"]
