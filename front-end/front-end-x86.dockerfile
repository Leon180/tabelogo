FROM alpine:latest
COPY frontEndAppX86 /
COPY templates /templates
COPY static /static
COPY app.env /
CMD [ "/frontEndAppX86" ]