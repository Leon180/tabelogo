FROM alpine:latest
RUN mkdir -p /app/config
COPY spiderApp /app
COPY ./config/config.env /app/config/
WORKDIR /app
CMD [ "./spiderApp" ]
