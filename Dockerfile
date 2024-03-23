FROM docker.io/golang:1.22

ENV SESSION_KEY="T0PS3CR3T"

WORKDIR /app

# Download Go modules
COPY app/go.mod app/go.sum ./
RUN go mod download

COPY app/*.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /injectionapp

EXPOSE 8080

# Run
CMD ["/injectionapp"]
