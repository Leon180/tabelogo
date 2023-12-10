FROM alpine:latest
RUN mkdir /app
COPY tabelogspiderApp /app
COPY ./cmd/api/app.env /
CMD [ "/app/tabelogspiderApp" ]