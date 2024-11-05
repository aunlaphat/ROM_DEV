FROM nginx:1.13-alpine

RUN apk add --no-cache --update tzdata \
    && cp /usr/share/zoneinfo/Asia/Bangkok /etc/localtime \
    && apk del tzdata \
    && rm -rf /var/cache/apk/* /tmp/* /var/tmp/* 

RUN ls -a

COPY ./default.conf /etc/nginx/conf.d/default.conf
COPY ./build /usr/share/nginx/html

EXPOSE 3000

WORKDIR /usr/share/nginx/html
