FROM golang:1.22 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/arbiter ./cmd/arbiter

FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/arbiter /bin/arbiter

EXPOSE 8443
ENTRYPOINT ["/bin/arbiter"]
