FROM golang:1.12.0-alpine3.9 AS builder
WORKDIR /go/src/github.com/BiteLikeASnake/products_microservise
COPY . .
RUN go install ./...

FROM alpine:3.9 AS production
#COPY --from=builder /go/bin/demo-docker ./app
COPY --from=builder /go/bin/cmd .
ENV DB="user=postgres password=example dbname=online_shop sslmode=disable port=5432 host=db" \
ADDRESS=":8081" \
ADMIN_TOKEN="adminpass" \
USER_TOKEN="userpass"
EXPOSE 8081