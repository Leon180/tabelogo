FROM alpine:latest
RUN mkdir /app
COPY loggerServiceApp /app
COPY ./cmd/api/app.env /
CMD [ "/app/loggerServiceApp" ]