FROM nginx:stable-alpine-slim
COPY  ./deployment/nginx/default.conf /etc/nginx/conf.d/default.conf
COPY ./deployment/nginx/start.sh /usr/src/start.sh

EXPOSE 8080
EXPOSE 8090

ENTRYPOINT [ "/usr/src/start.sh"]