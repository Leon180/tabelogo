FROM alpine:latest
RUN mkdir -p /app/config
COPY googlemapApp /app/
COPY ./config/config.env /app/config/
WORKDIR /app
CMD [ "./googlemapApp" ]
