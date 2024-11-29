FROM golang as build

WORKDIR /build

COPY ./ ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o njord

CMD ["/build/njord"]