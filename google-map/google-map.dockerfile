FROM alpine:latest
RUN mkdir /app
COPY googleMapApp /app
COPY ./cmd/api/app.env /
CMD [ "/app/googleMapApp" ]