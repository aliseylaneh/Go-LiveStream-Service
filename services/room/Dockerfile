# Start from a lightweight base image, here we're using Alpine Linux
FROM alpine:latest

WORKDIR /app

COPY ./vpeer_room /app/

ENV LISTEN_ADDRESS=0.0.0.0:8081
ENV DB_HOST=172.17.0.2
ENV DB_PORT=5432
ENV DB_NAME=room
ENV DB_USERNAME=postgres
ENV DB_PASSWORD=postgres

CMD ["sh", "-c", "./vpeer_room --listen-address=${LISTEN_ADDRESS} --db-host=${DB_HOST} --db-port=${DB_PORT} --db-name=${DB_NAME} --db-username=${DB_USERNAME} --db-password=${DB_PASSWORD}"]