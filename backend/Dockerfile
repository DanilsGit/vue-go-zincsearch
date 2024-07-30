# Utiliza la imagen base de golang
FROM golang:1.23rc2-bookworm

# Establece el directorio de trabajo en /home/app
WORKDIR /home/app

# Copia el go.mod y go.sum al directorio de trabajo
COPY go.mod go.sum ./

# Descarga las dependencias necesarias
RUN go mod download

# Copia todo el código fuente al directorio de trabajo
COPY . .

# Compila la aplicación
RUN go build -o main .

# Exponer el puerto 8080 en el contenedor
EXPOSE 8080

# Comando para ejecutar la aplicación compilada
CMD ["./main"]