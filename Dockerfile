FROM nginx:1.27-alpine
COPY webapp/default.conf /etc/nginx/conf.d/default.conf
COPY webapp/index.html /usr/share/nginx/html/index.html
EXPOSE 80