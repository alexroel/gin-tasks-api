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






