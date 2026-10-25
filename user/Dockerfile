# --- build ---
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.work go.work.sum ./
COPY shared/ ./shared/
COPY user/ ./user/
COPY url/ ./url/
COPY gateway/ ./gateway/
COPY analytics/ ./analytics/
RUN cd user && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /bin/app ./cmd

# --- run ---
FROM alpine:3.20
RUN adduser -D -u 10001 app
USER app
COPY --from=build /bin/app /bin/app
COPY --from=build /src/user/migrations /migrations
ENTRYPOINT ["/bin/app"]