FROM public.ecr.aws/docker/library/golang:1.23-alpine

WORKDIR /app

# Install Air for hot reload (optional now, but keeping for dev reference)
RUN go install github.com/air-verse/air@v1.52.3

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application to a location OUTSIDE /app to avoid volume mount conflict
RUN go build -o /usr/local/bin/giftflow-backend .

EXPOSE 8080

# Run the binary from the safe location
CMD ["/usr/local/bin/giftflow-backend"]
