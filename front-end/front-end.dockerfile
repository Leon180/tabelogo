FROM alpine:latest
COPY frontEndApp /
COPY templates /templates
COPY static /static
COPY app.env /
CMD [ "/frontEndApp" ]