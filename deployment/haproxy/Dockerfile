FROM haproxy:2.8-alpine
COPY /deployment/haproxy/haproxy.cfg /usr/local/etc/haproxy/haproxy.cfg


COPY ./deployment/haproxy/start.sh /usr/src/start.sh

ENTRYPOINT [ "/usr/src/start.sh"]

