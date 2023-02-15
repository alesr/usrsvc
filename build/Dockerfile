FROM alpine:latest

RUN apk add --update go

WORKDIR /app

COPY usrsvc .

EXPOSE 50051

ENTRYPOINT ["./usrsvc"]
