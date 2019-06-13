FROM nginx:alpine
LABEL author="Nikolay Vasilev"
COPY ./config/nginx.conf /etc/nginx/conf.d/default.conf

WORKDIR /usr/share/nginx/html
COPY dist/ .
 
# ng build --prod --configuration=bg-prod
# docker build -t nikolyvas/angular_dentalzone:1.0.4 -f nginx.dockerfile .
# docker push nikolyvas/angular_dentalzone:1.0.4
# docker run -p 8080:80 -v C:/Users/nvasilev/go/src/dentalzone/dental_ui/dist:/usr/share/nginx/html nginx-angular

# http-server dist/. - local