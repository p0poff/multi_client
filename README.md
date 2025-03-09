# multi_client
## compile so
go build -o multi_client.so -buildmode=c-shared client.go struct.go

## php.ini
apk --no-cache

## dockerfile added lib
RUN apk --no-cache php82-ffi