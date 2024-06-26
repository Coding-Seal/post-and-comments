FROM golang:alpine
RUN apk update && \
    apk add --no-cache make
LABEL authors="coding_seal"
WORKDIR /app
COPY . .
RUN make
EXPOSE 8080

ENTRYPOINT ["./main"]