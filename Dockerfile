FROM alpine:3.18
WORKDIR /app
COPY ./api .
COPY ./.env .


EXPOSE 8080
CMD ["/app/main"]
