FROM alpine:latest
RUN mkdir /app
COPY loggerServiceAppX86 /app
COPY ./cmd/api/app.env /
CMD [ "/app/loggerServiceAppX86" ]