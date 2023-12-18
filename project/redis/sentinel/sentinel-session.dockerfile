FROM redis:7.0.14-alpine
COPY sentinel-session.conf /etc/redis/sentinel-session.conf
ENTRYPOINT redis-server /etc/redis/sentinel.conf --sentinel