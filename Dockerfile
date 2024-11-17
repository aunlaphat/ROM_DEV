FROM alpine:3.18
WORKDIR /app
COPY ./api .
COPY ./.env .

RUN mkdir -p /app/assets
COPY ./assets/tahoma.ttf /app/assets/
COPY ./assets/tahomabd.ttf /app/assets/

EXPOSE 8080
CMD ["/app/api"]
