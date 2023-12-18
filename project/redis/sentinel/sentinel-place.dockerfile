FROM redis:7.0.14-alpine
COPY sentinel-place.conf /etc/redis/sentinel-place.conf
ENTRYPOINT redis-server /etc/redis/sentinel.conf --sentinel