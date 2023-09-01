# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.20 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /obsidian-sync ./cmd/obsidian-sync/

# Deploy the application binary into a lean image
FROM --platform=${TARGETPLATFORM} gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /obsidian-sync /obsidian-sync

# ENV HOST=localhost:3000
# ENV ADDR_HTTP=0.0.0.0:3000
# ENV DATA_DIR=/data

EXPOSE 3000

#USER nonroot:nonroot

VOLUME ["/data"]

ENTRYPOINT ["/obsidian-sync"]

CMD ["/obsidian-sync"]
