FROM alpine:latest
RUN mkdir /app
COPY authenticateApp /app
COPY ./cmd/api/app.env /
CMD [ "/app/authenticateApp" ]