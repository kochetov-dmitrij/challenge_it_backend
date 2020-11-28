FROM golang:1.15

ARG PORT=3000
EXPOSE ${PORT}

WORKDIR /app
COPY . .

RUN go get github.com/swaggo/swag/cmd/swag &&\
    swag init -d echo_server --output docs/echo_server

RUN go mod tidy
RUN cd echo_server && go build

CMD ["/app/echo_server/echo_server"]