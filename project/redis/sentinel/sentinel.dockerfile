FROM redis:7.0.14-alpine
COPY sentinel.conf /etc/redis/sentinel.conf
ENTRYPOINT redis-server /etc/redis/sentinel.conf --sentinel