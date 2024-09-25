# 1. Etapa de construcción (Builder)
FROM golang:1.22.2-alpine AS builder

# Establecemos el directorio de trabajo en el contenedor
WORKDIR /app

# Copiamos los archivos de dependencias (go.mod y go.sum)
COPY go.mod go.sum ./

# Descargamos las dependencias
RUN go mod download

# Copiamos el resto del código del proyecto
COPY . .

# Entramos a la carpeta cmd donde está el archivo main.go para compilarlo
WORKDIR /app/cmd

# Compilamos la aplicación de Go de forma estática
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# 2. Etapa final (runtime)
FROM alpine:latest

# Establecemos el directorio de trabajo en el contenedor final
WORKDIR /root/

# Copiamos el binario compilado desde la etapa de construcción
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Exponemos el puerto en el que la aplicación escucha (ej: 8080)
EXPOSE 8082

# Comando para ejecutar la aplicación
CMD ["./main"]
