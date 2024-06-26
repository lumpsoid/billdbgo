FROM golang:alpine as base
ENV BILLDB_TEMPLATE_PATH=/billdb/templates/* \
    BILLDB_STATIC_PATH=/billdb/static \
    BILLDB_DB_PATH=/billdb/data/bills.db
WORKDIR /billdb
COPY ./go.mod ./
COPY ./go.sum ./
RUN apk add gcc musl-dev
RUN go mod download

# Dev environment doesn't run this stage or beyond
FROM base as build
COPY ./cmd ./cmd
COPY ./internal ./internal
RUN mkdir ./data
ENV CGO_ENABLED=1
RUN go build -o /billdb/server ./cmd/server/main.go

# Production environment runs this stage
FROM alpine:latest
ENV BILLDB_TEMPLATE_PATH=/server/templates/* \
    BILLDB_STATIC_PATH=/server/static \
    BILLDB_DB_PATH=/server/data/bills.db
WORKDIR /server
COPY --from=build /billdb/server ./server
COPY --from=build /billdb/data ./data
COPY ./web/templates ./templates
COPY ./web/static ./static
EXPOSE 1323

CMD ["/server/server"]
