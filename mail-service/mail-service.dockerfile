FROM alpine:latest
RUN mkdir /app
COPY mailServiceApp /app
COPY templates /templates
COPY ./cmd/api/app.env /
CMD [ "/app/mailServiceApp"]