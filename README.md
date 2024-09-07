# Inventory App

Este proyecto es una API sencilla de un sistema de CRUD (Crear, Leer, Actualizar, Eliminar) desarrollada en Go (Golang). A través de este proyecto, aprendí los conceptos fundamentales para crear una API en Go, conectar una base de datos PostgreSQL, crear y relacionar tablas, y configurar un servidor HTTP.

## 🛠️ Características

- **API RESTful**: La API permite realizar operaciones CRUD sobre los recursos definidos.
- **Conexión a PostgreSQL**: Se utiliza PostgreSQL como base de datos relacional para almacenar y gestionar datos.
- **Migración automática**: Las tablas de la base de datos se crean y relacionan automáticamente usando GORM.
- **Servidor HTTP**: Implementación de un servidor HTTP para manejar las solicitudes a la API.

## 🧰 Tecnologías Utilizadas

- **Go (Golang)**: Lenguaje de programación utilizado para desarrollar la API.
- **GORM**: ORM (Object-Relational Mapping) para Go que facilita la interacción con la base de datos.
- **PostgreSQL**: Sistema de gestión de bases de datos relacional utilizado para almacenar los datos.
- **gorilla/mux**: Router HTTP para manejar las rutas de la API de manera eficiente.

## 📁 Estructura del Proyecto

La estructura del proyecto es la siguiente:

```plaintext
ecommerce-go/
├── db/                 # Configuración y manejo de la base de datos
│   └── connection.go           # Conexión a la base de datos y funciones relacionadas
├── models/             # Definición de modelos y relaciones
│   ├── user.go         # Modelo y migración de la tabla 'users'
│   └── task.go         # Modelo y migración de la tabla 'tasks'
├── routes/             # Definición de rutas de la API
│   └── routes.go       # Configuración de rutas
├── handlers/           # Lógica de las rutas y controladores
│   ├── user_handler.go # Lógica de negocio para las rutas de 'users'
│   └── task_handler.go # Lógica de negocio para las rutas de 'tasks'
├── main.go             # Punto de entrada de la aplicación
└── go.mod              # Archivo de módulos de Go (dependencias)
```
