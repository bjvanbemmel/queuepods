FROM golang:1.23.4-alpine3.21
WORKDIR /app
COPY ./* ./
RUN go mod download
RUN go install github.com/air-verse/air@latest
EXPOSE 8888
CMD [ "air", "-c", ".air.toml" ]
