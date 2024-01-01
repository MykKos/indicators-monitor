# Build app image
FROM golang:1.18.6 AS compile-image

WORKDIR /src
COPY ./go.mod ./go.sum ./

COPY . .

# Установка гита
RUN apt update && apt install git -y

# Добавляем приватный репо в гошную
ENV GOPRIVATE=github.com

# Логин/пасс для доступа к приватной либе
ENV USERNAME
ENV PASSWORD

# Запись в .netrc для получения доступа в приватный репо
RUN echo "machine github.com login $USERNAME password $PASSWORD" > /root/.netrc


RUN go mod download

ARG TARGET
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -installsuffix cgo -o ./app ./cmd/app/main.go
