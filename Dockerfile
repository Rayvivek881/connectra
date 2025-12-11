FROM golang:1.24-alpine AS tisvc

# Install git in the build stage
RUN apk add --no-cache git build-base libxml2-dev libxslt-dev

WORKDIR /opt/app

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .
RUN touch .env
RUN go build -o connectra .

# Build from alpine
FROM alpine:3.16
RUN apk add --no-cache bash bash-doc bash-completion openssl git libxml2 libxslt
WORKDIR /opt/app



ENV RUN_COMMAND=help
RUN touch .env
RUN mkdir data
RUN echo ${RUN_COMMAND}

COPY --from=tisvc /opt/app/connectra /opt/app/connectra

CMD ["sh", "-c", "./connectra $RUN_COMMAND"]