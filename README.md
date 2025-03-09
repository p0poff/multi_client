# multi_client
## compile so
go build -o multi_client.so -buildmode=c-shared client.go struct.go

## php.ini
apk --no-cache

## dockerfile added lib
RUN apk --no-cache php82-ffi

docker build -t multi_client .
docker run -it multi_client


docker-compose exec php-fpm apk add --no-cache gcc libc-dev
docker-compose exec php-fpm gcc -o /var/www/html/storage/app/temp/test_lib /var/www/html/storage/app/temp/multi.c -ldl
