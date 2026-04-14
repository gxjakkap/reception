FROM golang:alpine AS builder

RUN apk update --no-cache && apk upgrade --no-cache \
    && apk add --no-cache mailcap git tzdata ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -a -installsuffix cgo -o reception .

FROM gcr.io/distroless/static-debian12

COPY --from=builder /app/reception /reception

USER nonroot:nonroot

EXPOSE 8000
EXPOSE 8080

ENTRYPOINT ["/reception"]