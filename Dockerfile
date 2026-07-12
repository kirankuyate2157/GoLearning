# Build stage with Go 1.25
FROM golang:1.25-alpine as builder
WORKDIR /app

COPY . .

# Install swag and generate swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12
RUN swag init

RUN go build -o main .

# Final stage: Python 3.12 Alpine with Go app
FROM python:3.12-alpine
WORKDIR /app

# Copy built Go binary from builder
COPY --from=builder /app/main .

# Copy requirements.txt and install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

CMD ["./main"]