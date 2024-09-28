# Use Alpine Linux as the base image
FROM alpine:latest

# Set the working directory in the container
WORKDIR /app

# Copy the compiled binary into the container
COPY ./vpeer_usergw /app/

# Create the required directories
RUN mkdir -p /app/file_storage /app/minio-storage

# Install ca-certificates in case your application makes HTTPS requests, and any other dependencies
RUN apk --no-cache add ca-certificates && update-ca-certificates

# Set environment variables for the application's configuration
ENV LISTEN_ADDRESS="0.0.0.0:8080"
ENV ROOM_SERVER_ADDRESS="127.0.0.1:8081"
ENV FILE_SERVER_ADDRESS="127.0.0.1:8082"
ENV MINIO_SERVER_ADDRESS="185.36.231.135:9000"
ENV MINIO_ACCESS_KEY="t1Ke5iYwhOzrQmEdQsQg"
ENV MINIO_SECRET_KEY="vbRyrUz5B1banJYGUWKhTN5o1n0Sl7HVgsTefHXD"
ENV FILE_STORAGE_PATH="./file_storage"
ENV MINIO_DOWNLOADED_FILE_STORAGE_PATH="./minio-storage"

CMD ./vpeer_usergw --listen-address=${LISTEN_ADDRESS} \
    --room-server-address=${ROOM_SERVER_ADDRESS} \
    --file-server-address=${FILE_SERVER_ADDRESS} \
    --minio-server-address=${MINIO_SERVER_ADDRESS} \
    --minio-access-key=${MINIO_ACCESS_KEY} \
    --minio-secret-key=${MINIO_SECRET_KEY} \
    --file-storage-path=${FILE_STORAGE_PATH} \
    --minio-downloaded-file-storage-path=${MINIO_DOWNLOADED_FILE_STORAGE_PATH}
