# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.20 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /obsidian-sync ./cmd/obsidian-sync/

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /obsidian-sync /obsidian-sync

ENV HOST=0.0.0.0

EXPOSE 3000

#USER nonroot:nonroot

ENTRYPOINT ["/obsidian-sync"]

CMD ["/obsidian-sync"]
