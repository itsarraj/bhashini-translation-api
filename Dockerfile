# builder stage
FROM golang:alpine as builder

## ENV GO111MODULE=on

## it is important that these ARG's are defined after the FROM statement
# ARG PRIVATE_GIT_USER="<PRIVATE_GIT_USER>"
# ARG PRIVATE_GIT_PASS="<PRIVATE_GIT_PASS>"

## required for fetching private modules
# ENV GOPRIVATE <PRIVATE_GIT_URL>/module

## git is required to fetch go dependencies
RUN apk update \
    && apk add --no-cache ca-certificates \
    && apk add --update gcc musl-dev \
    && apk add --no-cache git \
    && update-ca-certificates

## create a netrc file using the credentials specified using --build-arg (for private modules)
# RUN echo "machine <PRIVATE_GIT_URL> login ${PRIVATE_GIT_USER} password ${PRIVATE_GIT_PASS}" > ~/.netrc
# RUN cat ~/.netrc

## set the working directory
WORKDIR /app

## fetch dependencies
COPY ./go.mod ./go.sum ./
RUN go mod download

## copy the source from the current directory to the working directory inside the container
COPY cmd cmd
COPY internal internal

## build the Go apps
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o app ./cmd

# start a new stage from scratch
FROM alpine:latest AS app

RUN apk add --no-cache ca-certificates \
    && update-ca-certificates

COPY --from=builder /app/app .

EXPOSE 3001

CMD ["./app"]
