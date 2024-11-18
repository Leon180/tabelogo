FROM alpine:latest
RUN mkdir /app
RUN mkdir /app/config
COPY authenticateApp /app
COPY ./config/config.env /app/config
CMD [ "/app/authenticateApp" ]
