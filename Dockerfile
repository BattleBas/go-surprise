FROM golang:1.12 AS build-env
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o surprise .

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build-env /src/surprise .
CMD ["./surprise"]