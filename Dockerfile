# imagen de go
FROM golang:1.20-bullseye

#variables de entorno de go
ENV GOOS linux
ENV CGO_ENABLED 1
ENV GOARCH amd64

#directorio donde se ubicara
WORKDIR /usr/src/app

#dependencias de go
COPY go.mod go.sum ./
RUN go mod download

#copiar el codigo a la carpeta base
COPY . .

#compilamos la aplicacion
RUN go build -o main .

# Comando para ejecutar la aplicación
CMD go run main.go

