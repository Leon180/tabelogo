FROM alpine:latest
RUN mkdir /app
COPY loggerServiceApp /app
COPY ./config/config_docker.env /app/config/config.env
CMD [ "/app/loggerServiceApp" ]