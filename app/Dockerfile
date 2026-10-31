# Stage 1: build the binary
FROM golang:1.26 AS build
WORKDIR /src

# Copy dependency files first so Docker can cache the download step
COPY go.mod go.sum ./
RUN go mod download

COPY app/ ./app/
# CGO off = a static binary that runs without system C libraries
RUN CGO_ENABLED=0 go build -o /departures-api ./app

# Stage 2: tiny runtime image, no shell, runs as non-root
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /departures-api /departures-api
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/departures-api"]