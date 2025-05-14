FROM alpine:3.18

# Устанавливаем Go, musl-dev и инструменты сборки
RUN apk add --no-cache go musl-dev gcc

WORKDIR /build
ADD client.go /build
ADD struct.go /build
ADD go.mod /build

RUN CGO_ENABLED=1 CC=cc go build -buildmode=c-shared -o multi_client.so client.go struct.go

CMD ["sh"]