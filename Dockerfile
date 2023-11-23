FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk add --update make

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN make build
WORKDIR /dist
RUN cp /build/main .

FROM scratch
COPY --from=builder /dist/main .

ENV DB_CONN=mongodb://mongodb-0.mongodb.default.svc.cluster.local:27017/?ssl=false&authSource=admin

EXPOSE 8080
ENTRYPOINT ["/main"]