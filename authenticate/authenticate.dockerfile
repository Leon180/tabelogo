FROM alpine:latest
RUN mkdir -p /app/config
COPY authenticateApp /app
COPY ./config/config_docker.env /app/config/config.env
WORKDIR /app
CMD [ "./authenticateApp" ]
