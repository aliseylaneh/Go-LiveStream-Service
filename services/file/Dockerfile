# Start from a lightweight base image, here we're using Alpine Linux
FROM alpine:latest

WORKDIR /app

COPY ./vpeer_file /app/

ENV LISTEN_ADDRESS=0.0.0.0:8082
ENV DB_HOST=host.docker.internal
ENV DB_PORT=5432
ENV DB_NAME=file
ENV DB_USERNAME=postgres
ENV DB_PASSWORD=postgres

CMD ["sh", "-c", "./vpeer_file --listen-address=${LISTEN_ADDRESS} --db-host=${DB_HOST} --db-port=${DB_PORT} --db-name=${DB_NAME} --db-username=${DB_USERNAME} --db-password=${DB_PASSWORD}"]