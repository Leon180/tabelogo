FROM alpine:latest
RUN mkdir /app
COPY tabelogspiderAppX86 /app
CMD [ "/app/tabelogspiderAppX86" ]