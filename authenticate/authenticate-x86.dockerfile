FROM alpine:latest
RUN mkdir /app
COPY authenticateAppX86 /app
COPY ./cmd/api/app.env /
CMD [ "/app/authenticateAppX86" ]