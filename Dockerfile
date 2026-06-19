# syntax=docker/dockerfile:1
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/local-solar-time ./cmd/local-solar-time

FROM gcr.io/distroless/static
COPY --from=build /out/local-solar-time /local-solar-time
ENTRYPOINT ["/local-solar-time"]
