# Documentación de APIs con Swagger

---
## Introducción

Cuando creamos una API, es importante que otros desarrolladores (o nosotros mismos en el futuro) sepan cómo usarla. Swagger nos ayuda a crear documentación automática y una interfaz visual para probar nuestra API.

---
## ¿Qué es Swagger / OpenAPI?

**Swagger** es una herramienta para documentar APIs REST. Ahora se llama **OpenAPI**.

**¿Qué hace?**
- Genera documentación automática de tu API
- Crea una interfaz web interactiva para probar endpoints
- Permite exportar la documentación en JSON o YAML

**Beneficios:**
- Los desarrolladores entienden rápidamente cómo usar tu API
- Puedes probar endpoints sin usar Postman
- La documentación siempre está actualizada con el código

---
## Instalación de Swagger en Go

### 1. Instalar la herramienta CLI

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Esta herramienta lee los comentarios de tu código y genera la documentación.

### 2. Instalar las librerías para Gin

```bash
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

---
## Configuración con Gin

### 1. Agregar anotaciones principales en main.go

Antes de `package main`, agrega la información general de tu API:

```go
// @title Tasks API
// @version 1.0
// @description API para la gestión de tareas con autenticación JWT.

// @contact.name Tu Nombre
// @contact.email tu@email.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Tipo de token JWT con el prefijo 'Bearer '

package main
```

### 2. Importar las librerías

```go
import (
    _ "tu-proyecto/docs"  // Documentación generada
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)
```

### 3. Agregar la ruta de Swagger

```go
router := gin.Default()

// Ruta de documentación
router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
```

---
## Documentación de Endpoints

Cada handler necesita comentarios especiales antes de la función.

### Estructura básica

```go
// NombreFuncion godoc
// @Summary      Título corto
// @Description  Descripción detallada
// @Tags         NombreGrupo
// @Accept       json
// @Produce      json
// @Router       /ruta [metodo]
func (h *Handler) NombreFuncion(c *gin.Context) {
```

### Ejemplo real: Login

```go
// Login godoc
// @Summary      Iniciar sesión
// @Description  Autentica un usuario y retorna un token JWT
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body domain.UserLogin true "Credenciales"
// @Success      200 {object} utils.Response "Login exitoso"
// @Failure      401 {object} utils.Response "Credenciales incorrectas"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
```

---
## Documentación de parámetros y Body

### Parámetros en la ruta (path)

```go
// @Param id path int true "ID de la tarea"
// @Router /tasks/{id} [get]
```

### Parámetros en query string

```go
// @Param page query int false "Número de página"
// @Param limit query int false "Límite por página"
```

### Body JSON

```go
// @Param request body domain.CreateTask true "Datos de la tarea"
```

**Formato:** `@Param nombre ubicación tipo requerido "descripción"`

| Ubicación | Uso |
|-----------|-----|
| `path` | En la URL: `/tasks/{id}` |
| `query` | Query string: `?page=1` |
| `body` | Cuerpo JSON |
| `header` | En headers |

---
## Documentación de respuestas

### Respuesta exitosa simple

```go
// @Success 200 {object} utils.Response "Mensaje de éxito"
```

### Respuesta con datos específicos

```go
// @Success 200 {object} utils.Response{data=domain.UserResponse} "Usuario obtenido"
```

### Respuesta con array

```go
// @Success 200 {object} utils.Response{data=[]domain.TaskResponse} "Lista de tareas"
```

### Múltiples respuestas de error

```go
// @Failure 400 {object} utils.Response "Datos inválidos"
// @Failure 401 {object} utils.Response "No autenticado"
// @Failure 404 {object} utils.Response "No encontrado"
// @Failure 500 {object} utils.Response "Error interno"
```

---
## Seguridad JWT en Swagger

### 1. Definir el esquema de seguridad (en main.go)

```go
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Tipo de token JWT con el prefijo 'Bearer '
```

### 2. Aplicar seguridad a endpoints protegidos

```go
// @Security BearerAuth
// @Router /tasks [get]
```

### 3. Usar en Swagger UI

1. Haz login y copia el token
2. Click en el botón **"Authorize"** 🔒
3. Escribe: `Bearer tu_token_aqui`
4. Click en **"Authorize"**
5. Ahora puedes probar endpoints protegidos

---
## Publicación de documentación

### Generar la documentación

```bash
swag init -g cmd/main.go -o docs
```

Esto crea la carpeta `docs/` con:
- `docs.go` - Código Go
- `swagger.json` - Formato JSON
- `swagger.yaml` - Formato YAML

### Ejecutar y acceder

```bash
go run cmd/main.go
```

Abre en tu navegador:
```
http://localhost:8080/swagger/index.html
```

### Regenerar después de cambios

Cada vez que modifiques las anotaciones, ejecuta:

```bash
swag init -g cmd/main.go -o docs
```

---
## ¿Cómo usar IA para generar documentación Swagger?

Puedes pedirle a la IA que genere las anotaciones. Ejemplo de prompt:

```
Genera la documentación Swagger para este handler de Go:

func (h *TaskHandler) Create(c *gin.Context) {
    // ... código
}

El endpoint:
- Crea una nueva tarea
- Requiere autenticación JWT
- Recibe domain.CreateTask en el body
- Retorna domain.TaskResponse
```

**Tips:**
- Describe qué hace el endpoint
- Menciona si necesita autenticación
- Indica los tipos de entrada y salida
- Especifica los posibles errores

---
## Conclusión

Swagger es una herramienta esencial para documentar APIs profesionales:

✅ **Aprendiste:**
- Qué es Swagger y para qué sirve
- Cómo instalarlo en un proyecto Go + Gin
- Cómo documentar endpoints, parámetros y respuestas
- Cómo agregar seguridad JWT
- Cómo generar y acceder a la documentación

📌 **Recuerda:**
- Ejecuta `swag init` después de cada cambio
- Las rutas en `@Router` no incluyen el `@BasePath`
- Usa `@Security BearerAuth` en endpoints protegidos
- Los tipos deben existir en tu código (no usar funciones)

🚀 **Siguiente paso:** Probar tu API desde Swagger UI en `http://localhost:8080/swagger/index.html`