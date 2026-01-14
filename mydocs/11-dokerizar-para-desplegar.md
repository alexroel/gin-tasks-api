# Dockerizar para desplegar

---
## Introducción


---
## Instalación de Docker


---
## Creación de un contenedor de PostgreSQL con Docker Compose
Crearemos un archivo llamado `docker-compose.yml` en la raíz del proyecto con el siguiente contenido:

```yml
services:
  postgres:
    build: .
    container_name: postgres-tasks-db
    ports:
      - "5432:5432"
    environment:
      POSTGRES_DB: tasks_db
      POSTGRES_USER: roelcode
      POSTGRES_PASSWORD: 123456
```

Ahora para crear una imagen de Docker y ejecutar el contenedor, utilizamos el siguiente comando en la terminal:

```bash
docker-compose up -d --build
```
Esto construirá la imagen de Docker y levantará el contenedor de PostgreSQL con la configuración especificada.

Podemos verificar que el contenedor esté corriendo con el comando:

```bash
docker ps
```

Para detener y eliminar el contenedor, utilizamos:

```bash
docker-compose down
```
Para listar e eliminar imágenes de Docker que ya no necesitamos, podemos usar los siguientes comandos:

```bash
docker images
docker rmi <image_id>
```
Reemplazando `<image_id>` con el ID de la imagen que queremos eliminar.

Ahoara que tenemos nuestro contenedor de PostgreSQL corriendo, podemos proceder a configurar nuestra aplicación para que se conecte a esta base de datos en el contenedor.

Para esto vamos a `.env` y actualizamos la variable `URL_DATABASE` para que apunte al contenedor de Docker:

```env
URL_DATABASE=host=localhost user=roelcode password=123456 dbname=tasks_db port=5432 sslmode=disable TimeZone=America/Lima
```

---
## Cargar variables de entorno desde un archivo `.env` en Docker Compose
Para cargar las variables de entorno desde un archivo `.env` en Docker Compose, primero creamos un archivo llamado `.env` en la raíz del proyecto con el siguiente contenido:

```env
PQSL_DB=tasks_db
PQSL_USER=roelcode
PQSL_PASSWORD=123456
```

Luego, en el archivo `docker-compose.yml`, actualizamos la sección de `environment` del servicio de PostgreSQL para que utilice estas variables de entorno:

```yml
services:
  postgres:
    build: .
    container_name: postgres-tasks-db
    ports:
      - "5432:5432"
    environment:
      PQSL_DB: ${PQSL_DB}
      PQSL_USER: ${PQSL_USER}
      PQSL_PASSWORD: ${PQSL_PASSWORD}
```

Además, actualizamos el comando de `healthcheck` para que también utilice estas variables de entorno:

```yml
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${PQSL_USER} -d ${PQSL_DB}"]
      interval: 5s
      timeout: 5s
      retries: 10
```

- `healthcheck`: Define una verificación de salud para el contenedor.
  - `test`: Comando para verificar si PostgreSQL está listo, utilizando las variables de entorno.
  - `interval`: Intervalo entre cada verificación.
  - `timeout`: Tiempo máximo para esperar una respuesta.
  - `retries`: Número de intentos antes de considerar que el contenedor no está saludable.

Con estos cambios, Docker Compose cargará las variables de entorno desde el archivo `.env` y las utilizará para configurar el contenedor de PostgreSQL.

---
## Dockerizando una API en Go
Para dockerizar una API en Go, primero necesitamos crear un archivo llamado `Dockerfile` en la raíz del proyecto con el siguiente contenido:

```Dockerfile
# ============================================
# STAGE 1: Build
# ============================================
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copiar dependencias primero (mejor cache)
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar binario estático
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/main.go

# ============================================
# STAGE 2: Runtime (imagen mínima)
# ============================================
FROM alpine:3.19

# Agregar certificados SSL y timezone
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copiar solo el binario compilado
COPY --from=builder /app/main .

# Usuario no-root para seguridad
RUN adduser -D -g '' appuser
USER appuser

EXPOSE 8080

CMD ["./main"]
```

- `FROM golang:1.25-alpine AS builder`: Utiliza una imagen base de Go para compilar la aplicación.
- `WORKDIR /app`: Establece el directorio de trabajo dentro del contenedor.
- `COPY go.mod go.sum ./` y `RUN go mod download`: Copia los archivos de dependencias y descarga las dependencias.
- `COPY . .`: Copia todo el código fuente al contenedor.
- `RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/main.go`: Compila la aplicación en un binario estático.
- `FROM alpine:3.19`: Utiliza una imagen mínima de Alpine Linux para el entorno de ejecución.
- `RUN apk --no-cache add ca-certificates tzdata`: Instala certificados SSL y datos de zona horaria.
- `COPY --from=builder /app/main .`: Copia el binario compilado desde la etapa de construcción. 
- `RUN adduser -D -g '' appuser` y `USER appuser`: Crea y utiliza un usuario no-root para mayor seguridad.
- `EXPOSE 8080`: Expone el puerto 8080 para la aplicación.
- `CMD ["./main"]`: Comando para ejecutar la aplicación cuando se inicie el contenedor.

Ahora, para construir la imagen de Docker y ejecutar el contenedor de la API en Go, utilizamos el siguiente comando en la terminal:

```bash
docker build -t gin-tasks-api .
docker run -d -p 8080:8080 --name gin-tasks-api-container gin-tasks-api
```
Esto construirá la imagen de Docker y levantará el contenedor de la API en Go.

---
## Agregar la API de Go al archivo Docker Compose
Para agregar la API de Go al archivo `docker-compose.yml`, actualizamos el archivo para incluir un nuevo servicio para la API:

Esto es bueno para enlazar ambos contenedores y facilitar su gestión conjunta.

```yml
services:
  postgres:
    image: postgres:17-alpine
    container_name: postgres-tasks-db
    restart: unless-stopped
    environment:
      POSTGRES_DB: ${PQSL_DB}
      POSTGRES_USER: ${PQSL_USER}
      POSTGRES_PASSWORD: ${PQSL_PASSWORD}
    ports:
      - "5433:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${PQSL_USER} -d ${PQSL_DB}"]
      interval: 5s
      timeout: 5s
      retries: 10

  api:
    build: .
    container_name: gin-tasks-api
    restart: unless-stopped
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      URL_DATABASE: "postgresql://${PQSL_USER}:${PQSL_PASSWORD}@postgres:5432/${PQSL_DB}?sslmode=disable"
      JWT_SECRET: ${JWT_SECRET}
      GIN_MODE: ${GIN_MODE}
      PORT: ":8080"

volumes:
  postgres_data:
```

- `api`: Define un nuevo servicio para la API en Go.
  - `build: .`: Construye la imagen de Docker utilizando el `Dockerfile` en la raíz del proyecto.
  - `container_name: gin-tasks-api`: Asigna un nombre al contenedor de la API.
  - `restart: unless-stopped`: Configura el contenedor para que se reinicie automáticamente a menos que se detenga manualmente.
  - `ports`: Mapea el puerto 8080 del contenedor al puerto 8080 del host.
  - `depends_on`: Asegura que el servicio de PostgreSQL esté saludable antes de iniciar la API.
  - `environment`: Define las variables de entorno necesarias para la API, incluyendo la cadena de conexión a la base de datos utilizando las variables definidas en el archivo `.env`.

Con estos cambios, podemos levantar ambos contenedores (PostgreSQL y la API en Go) utilizando el comando:

```bash
docker-compose up -d --build
```
Esto construirá las imágenes de Docker y levantará ambos contenedores con la configuración especificada.

---
## Crear cuenta gratuita en Render.com
Render.com es una plataforma de alojamiento en la nube que permite desplegar aplicaciones web, bases de datos y otros servicios de manera sencilla. Para crear una cuenta gratuita en Render.com, sigue estos pasos:

1. Visita el sitio web de Render.com en [https://render.com/](https://render.com/).
2. Haz clic en el botón "Sign Up" o "Get Started for Free".
3. Completa el formulario de registro con tu dirección de correo electrónico, nombre de usuario y contraseña, o utiliza una cuenta de GitHub o Google para registrarte.
4. Verifica tu dirección de correo electrónico si es necesario.
5. Una vez registrado, inicia sesión en tu cuenta de Render.com.
Ahora tienes una cuenta gratuita en Render.com y puedes comenzar a desplegar tus aplicaciones y servicios en la plataforma.

---
## Crear una base de datos PostgreSQL en Render.com
Para crear una base de datos PostgreSQL en Render.com, sigue estos pasos:

1. Inicia sesión en tu cuenta de Render.com.
2. En el panel de control, haz clic en el botón "New" y selecciona "Database".
3. Elige "PostgreSQL" como el tipo de base de datos.
4. Configura los detalles de la base de datos:
   - Nombre de la base de datos: Ingresa un nombre para tu base de datos.
   - Plan: Selecciona el plan gratuito (Free).
   - Región: Elige la región más cercana a tus usuarios para mejorar el rendimiento.
5. Haz clic en "Create Database" para crear la base de datos.
6. Una vez creada la base de datos, Render.com te proporcionará la cadena de conexión que necesitarás para conectar tu aplicación a la base de datos. Asegúrate de copiar esta información y guardarla en un lugar seguro. 

---
## Desplegar la API en Go en Render.com
Para desplegar la API en Go en Render.com, primero prepararemos nuestro proyecto para el despliegue. Asegúrate de que tu proyecto tenga un archivo `Dockerfile` y un archivo `docker-compose.yml` correctamente configurados.

Luego, sigue estos pasos para desplegar la API en Render.com:

1. Crear un archivo `.render.yaml` en la raíz del proyecto con el siguiente contenido:

```yaml
# render.yaml - Configuración de despliegue en Render.com
# Documentación: https://render.com/docs/blueprint-spec

services:
  # Web Service para la API
  - type: web
    name: gin-tasks-api
    runtime: docker
    region: oregon
    plan: free
    branch: main
    healthCheckPath: /swagger/index.html
    envVars:
      - key: URL_DATABASE
        sync: false  # Se configura manualmente en el dashboard
      - key: JWT_SECRET
        generateValue: true  # Render genera un valor seguro automáticamente
      - key: JWT_EXPIRE_IN
        value: 48h
      - key: GIN_MODE
        value: release
      - key: PORT
        value: 8080
```
Si no creas este archivo, Render.com intentará desplegar tu aplicación utilizando configuraciones predeterminadas, lo que podría no ser adecuado para tu proyecto específico.

2. Inicializa un repositorio Git en tu proyecto si aún no lo has hecho:

```bash
git init
git add .
git commit -m "Initial commit"
```

3. Sube tu código a un repositorio en GitHub, GitLab o Bitbucket.
4. En Render.com, haz clic en "New" y selecciona "Web Service".
5. Conecta tu cuenta de GitHub, GitLab o Bitbucket y selecciona el repositorio donde está tu proyecto.
6. Render.com detectará automáticamente el archivo `render.yaml` y utilizará la configuración especificada para desplegar tu API en Go.
7. Configura las variables de entorno necesarias en el dashboard de Render.com, especialmente la variable `URL_DATABASE` con la cadena de conexión a la base de datos PostgreSQL que creaste anteriormente.
8. Haz clic en "Create Web Service" para iniciar el proceso de despliegue.
9. Render.com construirá la imagen de Docker y desplegará tu API en Go. Puedes monitorear el progreso del despliegue en el dashboard de Render.com.

---
## Verificar el despliegue
Una vez que el despliegue se haya completado, puedes verificar que tu API en Go esté funcionando correctamente accediendo a la URL proporcionada por Render.com. Deberías poder ver la documentación Swagger de tu API en la ruta `/swagger/index.html`.

---
## Conclusión
En esta guía, hemos aprendido cómo dockerizar una API en Go y una base de datos PostgreSQL, y cómo desplegarlos en Render.com utilizando Docker Compose y un archivo de configuración `render.yaml`. Con estos conocimientos, ahora puedes desplegar tus propias aplicaciones web de manera eficiente y escalable en la nube.






