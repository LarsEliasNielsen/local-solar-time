# syntax=docker/dockerfile:1
FROM node:22-alpine AS web
WORKDIR /web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

FROM golang:1.27.2-alpine AS build
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/local-solar-time ./cmd/local-solar-time

FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/local-solar-time /local-solar-time
COPY --from=web /web/dist /static
ENV SOLAR_STATIC_DIR=/static
EXPOSE 8080
ENTRYPOINT ["/local-solar-time"]
