FROM alpine:latest
RUN mkdir /app
COPY tabelogspiderApp /app
CMD [ "/app/tabelogspiderApp" ]