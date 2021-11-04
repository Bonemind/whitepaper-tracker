FROM node:alpine AS npmbuild
WORKDIR /build
COPY frontend/ /build
RUN npm ci
RUN npm run build

FROM golang:alpine AS appbuild
RUN apk add --no-cache build-base
WORKDIR /app
COPY ./backend /app
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o whitepaper-tracker .

FROM alpine:latest

WORKDIR /app
RUN mkdir /app/frontend

COPY --from=appbuild /app/whitepaper-tracker ./whitepaper-tracker
COPY --from=npmbuild /build/public/ ./frontend

EXPOSE 3000

ENTRYPOINT ["/app/whitepaper-tracker"]