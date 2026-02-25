FROM public.ecr.aws/docker/library/golang:1.23-alpine

WORKDIR /app

# Install Air for hot reload (pin version to avoid Go version mismatch)
RUN go install github.com/air-verse/air@v1.52.3

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["air"]
