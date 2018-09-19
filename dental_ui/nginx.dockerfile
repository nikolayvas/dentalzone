FROM nginx:alpine
LABEL author="Nikolay Vasilev"
COPY ./config/nginx.conf /etc/nginx/conf.d/default.conf

WORKDIR /usr/share/nginx/html
COPY dist/ .

# docker build -t angular_dentalzone -f nginx.dockerfile .
# docker run -p 8080:80 -v C:/Users/nvasilev/go/src/dentalzone/dental_ui/dist:/usr/share/nginx/html nginx-angular