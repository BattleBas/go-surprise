FROM golang:1.12-alpine AS build-env
RUN apk --no-cache add ca-certificates && \
    apk add git
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o surprise cmd/surprise-server/main.go

FROM scratch
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-env /src/surprise /
CMD ["./surprise"]