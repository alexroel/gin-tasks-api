# 🚀 Gin Tasks API

API REST para gestión de tareas construida con Go y el framework Gin. Implementa autenticación JWT, arquitectura limpia y buenas prácticas de desarrollo.

## 📋 Tabla de Contenidos

- [Características](#-características)
- [Tecnologías](#-tecnologías)
- [Arquitectura](#-arquitectura)
- [Instalación](#-instalación)
- [Configuración](#-configuración)
- [Ejecución](#-ejecución)
- [API Endpoints](#-api-endpoints)
- [Testing](#-testing)
- [Documentación](#-documentación)

## ✨ Características

- ✅ CRUD completo de tareas
- ✅ Autenticación y autorización con JWT
- ✅ Registro e inicio de sesión de usuarios
- ✅ Hash seguro de contraseñas con bcrypt
- ✅ Validación de datos de entrada
- ✅ Documentación automática con Swagger
- ✅ Arquitectura limpia (Clean Architecture)
- ✅ Testing unitario, de integración y E2E
- ✅ Middleware de autenticación

## 🛠 Tecnologías

| Tecnología | Versión | Descripción |
|------------|---------|-------------|
| Go | 1.25.5 | Lenguaje de programación |
| Gin | 1.11.0 | Framework web HTTP |
| GORM | 1.31.1 | ORM para Go |
| MySQL | - | Base de datos relacional |
| JWT | 5.3.0 | Autenticación con tokens |
| Swagger | 1.16.6 | Documentación de API |
| Testify | 1.11.1 | Framework de testing |

## 🏗 Arquitectura

El proyecto sigue los principios de **Clean Architecture**:

```
gin-tasks-api/
├── cmd/                    # Punto de entrada de la aplicación
│   └── main.go
├── internal/               # Código privado de la aplicación
│   ├── config/            # Configuración y conexión a BD
│   ├── domain/            # Entidades del dominio
│   ├── handler/           # Controladores HTTP
│   ├── middleware/        # Middlewares (auth, etc.)
│   ├── repository/        # Acceso a datos
│   │   └── mocks/         # Mocks para testing
│   └── service/           # Lógica de negocio
│       └── mocks/         # Mocks para testing
├── pkg/                   # Código reutilizable público
│   ├── jwt/              # Utilidades JWT
│   └── utils/            # Utilidades generales
├── tests/                 # Tests E2E e integración
│   └── e2e/
├── docs/                  # Documentación Swagger generada
└── mydocs/               # Documentación del curso
```

### Flujo de Datos

```
Request → Handler → Service → Repository → Database
                ↓
            Response
```

## 📦 Instalación

### Prerrequisitos

- Go 1.25.5 o superior
- MySQL 8.0 o superior
- Git

### Clonar el repositorio

```bash
git clone https://github.com/alexroel/gin-tasks-api.git
cd gin-tasks-api
```

### Instalar dependencias

```bash
go mod download
```

## ⚙️ Configuración

### 1. Crear archivo de entorno

Copia el archivo de ejemplo y configura las variables:

```bash
cp .env.example .env
```

### 2. Editar variables de entorno

```env
# Servidor
PORT=8080
GIN_MODE=debug

# Base de datos
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=tu_password
DB_NAME=gin_tasks_db

# JWT
JWT_SECRET=tu_clave_secreta_muy_segura
JWT_EXPIRE_IN=24
```

### 3. Crear base de datos

```sql
CREATE DATABASE gin_tasks_db;
```

## 🚀 Ejecución

### Modo desarrollo

```bash
go run cmd/main.go
```

### Compilar y ejecutar

```bash
go build -o bin/api cmd/main.go
./bin/api
```

### Con hot-reload (usando air)

```bash
# Instalar air
go install github.com/air-verse/air@latest

# Ejecutar
air
```

## 📡 API Endpoints

### Autenticación

| Método | Endpoint | Descripción | Auth |
|--------|----------|-------------|------|
| POST | `/api/auth/register` | Registrar usuario | ❌ |
| POST | `/api/auth/login` | Iniciar sesión | ❌ |

### Tareas

| Método | Endpoint | Descripción | Auth |
|--------|----------|-------------|------|
| GET | `/api/tasks` | Listar tareas del usuario | ✅ |
| GET | `/api/tasks/:id` | Obtener tarea por ID | ✅ |
| POST | `/api/tasks` | Crear nueva tarea | ✅ |
| PUT | `/api/tasks/:id` | Actualizar tarea | ✅ |
| DELETE | `/api/tasks/:id` | Eliminar tarea | ✅ |

### Ejemplos de uso

#### Registro de usuario

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "usuario",
    "email": "usuario@email.com",
    "password": "password123"
  }'
```

#### Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@email.com",
    "password": "password123"
  }'
```

#### Crear tarea (con token)

```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <tu_token>" \
  -d '{
    "title": "Mi nueva tarea",
    "description": "Descripción de la tarea",
    "completed": false
  }'
```

## 🧪 Testing

### Ejecutar todos los tests

```bash
go test ./... -v
```

### Tests con cobertura

```bash
go test ./... -cover
```

### Reporte de cobertura detallado

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Ejecutar tests específicos

```bash
# Tests de un paquete
go test ./internal/service/... -v

# Tests E2E
go test ./tests/e2e/... -v

# Tests por nombre
go test ./... -run TestLogin -v
```

### Benchmarks

```bash
go test ./tests/e2e/... -bench=. -benchmem
```

### Cobertura actual

| Paquete | Cobertura |
|---------|-----------|
| middleware | 100% |
| service | 86.7% |
| handler | 75.3% |
| config | 60%+ |
| pkg/utils | 70%+ |
| pkg/jwt | 80%+ |

## 📚 Documentación

### Swagger UI

Una vez la aplicación esté corriendo, accede a la documentación interactiva:

```
http://localhost:8080/swagger/index.html
```

### Regenerar documentación Swagger

```bash
# Instalar swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generar docs
swag init -g cmd/main.go -o docs
```

## 🔒 Seguridad

- Las contraseñas se hashean con bcrypt (cost factor 14)
- Los tokens JWT expiran según configuración
- Validación de entrada en todos los endpoints
- Middleware de autenticación para rutas protegidas

## 📝 Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para más detalles.

## 👤 Autor

**Alex Roel**

- GitHub: [@alexroel](https://github.com/alexroel)

---

⭐ Si este proyecto te fue útil, ¡dale una estrella!
