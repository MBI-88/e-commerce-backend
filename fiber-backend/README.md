# Nombre de la Aplicación

Fiber-Backend (Go)

## Requisitos Previos

- Go (versión 1.20.5)
- Otros requisitos previos específicos, si los hay.

## Instalación

1. Clona el repositorio usando ssh: `git clone git@github.com:Proyectos-MBI-YAF/fiber-backend.git`
2. Ingresa al directorio del proyecto: `cd fiber-backend`
3. Instala las dependencias: `go mod tidy`

## Configuración

1. Abre el archivo de configuración: `config/config.yaml`.
2. Modifica los valores de configuración según sea necesario.

## Uso

1. Ejecuta la aplicación: `go run main.go`.
2. La aplicación estará disponible en `http://localhost:3000`.

## Estructura del Proyecto

- `main.go`: Archivo principal que inicializa la aplicación y configura los middlewares.
- `router/`: Directorio que contiene el manejador principal de las rutas.
- `models/`: Directorio que contiene las definiciones de los modelos de datos. Serializadores y logica de operación
- `controllers.go`: Archivo que controla las operaciones. Contiene las funciones y metodos.
- `middlewares/`: Directorio que contiene los middleware de todos los servicios.
- `helpers/`: Directorio que contiene fuciones de ayuda en los controladores
- Otros directorios y archivos según las necesidades de tu proyecto.

## Contribución

1. Haz un fork del repositorio.
2. Crea una nueva rama: `git checkout -b feature/nueva-caracteristica`.
3. Realiza tus cambios y realiza commits: `git commit -am 'Agrega nueva característica'`.
4. Envía tus cambios al repositorio remoto: `git push origin feature/nueva-caracteristica`.
5. Crea una solicitud de extracción en GitHub.

## Autor

- MBI-YAF
- Contacto: ingmbi8807@gmail.com, yasmaaf95@gmail.com


