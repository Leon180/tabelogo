FROM alpine:latest
RUN mkdir /app
RUN mkdir /app/config
COPY googlemapApp /app
COPY ./config/app.env /app/config
CMD [ "/app/googlemapApp" ]
