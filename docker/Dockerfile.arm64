# Use golang base image on Alpine
FROM golang:alpine as base

# Set environment variables
ENV BILLDB_TEMPLATE_PATH=/billdb/templates/* \
    BILLDB_STATIC_PATH=/billdb/static \
    BILLDB_DB_PATH=/billdb/data/bills.db

# Install necessary packages
RUN apk update && apk add --no-cache \
    gcc \
    musl-dev

# Set working directory
WORKDIR /billdb

# Copy go.mod and go.sum and download dependencies
COPY ./go.mod ./go.sum ./
RUN go mod download

# Stage for building the application
FROM base as build

# Copy the source code
COPY ./cmd ./cmd
COPY ./internal ./internal

# Create necessary directory
RUN mkdir ./data

# Set environment variables and build the application
ENV CGO_ENABLED=1
RUN go build -o /billdb/server ./cmd/server/main.go

# Production environment runs this stage
FROM alpine:latest
ENV BILLDB_TEMPLATE_PATH=/server/templates/* \
    BILLDB_STATIC_PATH=/server/static \
    BILLDB_DB_PATH=/server/data/bills.db
WORKDIR /server
COPY --from=build /billdb/server ./server
#COPY --from=build /billdb/data ./data
COPY ./web/templates ./templates
COPY ./web/static ./static
EXPOSE 1323

CMD ["/server/server"]