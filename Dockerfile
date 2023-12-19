FROM golang:1.21

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY lib ./lib
COPY service ./service

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /lemmy-links

CMD ["/lemmy-links"]
