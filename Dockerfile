FROM golang:latest AS go_builder
ADD . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o /server .

FROM node:alpine AS node_builder
COPY --from=go_builder /app/go2sheep_pusher /frontend
WORKDIR /frontend
RUN yarn
ARG VUE_APP_PUSHER_INSTANCE_ID
RUN echo "${VUE_APP_PUSHER_INSTANCE_ID}"
RUN yarn build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=go_builder /server .
COPY --from=node_builder /static ./static
RUN chmod +x ./server
CMD ./server