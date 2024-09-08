# Usa una imagen base de Go
FROM golang:1.22.2 AS builder

# Establece el directorio de trabajo
WORKDIR /app

# Copia los archivos de módulo
COPY go.mod go.sum ./

# Descarga las dependencias
RUN go mod download

# Copia el código fuente
COPY . .

# Construye la aplicación
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/main .

# Usa una imagen base con herramientas de diagnóstico para verificar
FROM golang:1.22.2

# Copia el binario desde el contenedor de construcción
COPY --from=builder /app/bin/main /bin/main

# Expone el puerto que la aplicación usará
EXPOSE 8080

# Comando para ejecutar la aplicación
ENTRYPOINT ["/bin/main"]
