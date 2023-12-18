FROM alpine:latest
RUN mkdir /app
COPY brokerAppX86 /app
CMD [ "/app/brokerAppX86" ]