FROM golang:1.18-alpine

WORKDIR /app
COPY ./bin/backend /app

EXPOSE 4000
CMD ["./backend"]