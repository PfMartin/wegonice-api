FROM debian:latest

ARG binary_name=wegonice-api

WORKDIR /app

COPY ./${binary_name}-amd64 /app/${binary_name}-amd64
COPY ./${binary_name}-arm64 /app/${binary_name}-arm64

RUN chmod +x ./${binary_name}-amd64
RUN chmod +x ./${binary_name}-arm64

COPY ./run_api.sh ./run_api.sh

EXPOSE 8000

CMD ["./run_api.sh"]