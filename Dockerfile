FROM golang:latest

WORKDIR /app
ADD . /app
RUN cd /app && make

CMD ["./main"]