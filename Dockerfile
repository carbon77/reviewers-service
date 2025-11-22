# Build
FROM golang:1.25 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Runtime
FROM scratch

COPY --from=build /app/main /main

ENTRYPOINT [ "/main" ]
