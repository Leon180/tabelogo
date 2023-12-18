FROM alpine:latest

RUN mkdir /app

COPY listenerAppX86 /app

CMD [ "/app/listenerAppX86"]