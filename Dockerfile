FROM golang:alpine
RUN apk add git
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main .
EXPOSE 8000
CMD ["/app/main"]