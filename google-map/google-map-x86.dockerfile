FROM alpine:latest
RUN mkdir /app
COPY googleMapAppX86 /app
COPY ./cmd/api/app.env /
CMD [ "/app/googleMapAppX86" ]