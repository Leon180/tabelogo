FROM alpine:latest
RUN mkdir /app
COPY spiderApp /app
CMD [ "/app/spiderApp" ]
