FROM nginx:alpine
LABEL author="Nikolay Vasilev"
COPY ./config/nginx.conf /etc/nginx/conf.d/default.conf