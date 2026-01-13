# Profundizando en Gin

1. [Introducción](#introducción)
2. [Mejorando la estructura del proyecto](#mejorando-la-estructura-del-proyecto)
3. [Manejo de Request Body](#manejo-de-request-body)
4. [Grupos de rutas](#grupos-de-rutas)
5. [Middlewares: concepto y uso](#middlewares-concepto-y-uso)
6. [Creación de middleware personalizado](#creación-de-middleware-personalizado)
7. [Manejo de CORS](#manejo-de-cors)
8. [Manejo de Headers](#manejo-de-headers)
9. [Manejo global de errores](#manejo-global-de-errores)
10. [Buenas prácticas con Gin](#buenas-prácticas-con-gin)
11. [Conclusión](#conclusión)
 

---
## Introducción
Bienvenido a la sección de Profundizando en Gin. En esta sección, exploraremos características avanzadas del framework Gin para construir aplicaciones web robustas y eficientes en Go. Aprenderemos sobre el manejo de rutas, middlewares, validación de datos, manejo de errores y buenas prácticas para trabajar con Gin.

---
## Mejorando la estructura del proyecto
Para empezar a profundizar en Gin, es esencial comprender el concepto de `gin.Context`. Primero, vamos a modularizar nuetro proyecto: 

Creares una carpeta llamada `storage` para manejar la persistencia de datos. Dentro de esta carpeta, crea un archivo `storage.go` donde implementaremos funciones para guardar y cargar datos desde un archivo JSON.

```go
// storage/storage.go
package storage

import (
	"encoding/json"
	"myapi/models"
	"os"
)

const dataFile = "data/users.json"

// Data estructura para almacenar todos los datos
type Data struct {
	Users  []models.User `json:"users"`
	NextID int           `json:"next_id"`
}

// LoadData carga los datos desde el archivo JSON
func LoadData() error {
	// Leer el archivo
	file, err := os.ReadFile(dataFile)
	if err != nil {
		return err
	}

	// Si el archivo está vacío, inicializar con datos vacíos
	if len(file) == 0 {
		models.Users = []models.User{}
		models.NextID = 1
		return nil
	}

	// Deserializar los datos
	var data Data
	if err := json.Unmarshal(file, &data); err != nil {
		return err
	}

	models.Users = data.Users
	models.NextID = data.NextID

	return nil
}

// SaveData guarda los datos en el archivo JSON
func SaveData() error {
	data := Data{
		Users:  models.Users,
		NextID: models.NextID,
	}

	// Serializar los datos con formato legible
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Escribir al archivo
	return os.WriteFile(dataFile, jsonData, 0644)
}
```
 Esto nos permitirá mantener nuestros datos persistentes entre reinicios de la aplicación. 

Luego creresps una carpeta llamada `models` para definir nuestras estructuras de datos. Dentro de esta carpeta, crea un archivo `user.go` para definir la estructura del usuario.

```go
// models/user.go
package models

// Estructura para representar un usuario
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Colección de Usuarios en memoria
var Users []User
var NextID = 1
```

Ahoara para manejar los handlers de usuario, crea una carpeta llamada `handlers` y dentro un archivo `user_handlers.go`.

```go
// handlers/user_handlers.go
package handlers

import (
	"fmt"
	"myapi/models"
	"myapi/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Crear un nuevo usuario
func CreateUser(c *gin.Context) {
	var newUser models.User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUser.ID = fmt.Sprintf("%d", models.NextID)
	models.NextID++
	models.Users = append(models.Users, newUser)

	// Guardar en archivo JSON
	if err := storage.SaveData(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar datos"})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

// Obtener todos los usuarios
func GetUsers(c *gin.Context) {
	if len(models.Users) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No hay usuarios disponibles"})
		return
	}
	c.JSON(http.StatusOK, models.Users)
}

// Obtener un usuario por ID
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	for _, user := range models.Users {
		if user.ID == id {
			c.JSON(http.StatusOK, user)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
}

// Actualizar un usuario existente
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var updatedUser models.User
	if err := c.BindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, user := range models.Users {
		if user.ID == id {
			updatedUser.ID = id
			models.Users[i] = updatedUser

			// Guardar en archivo JSON
			if err := storage.SaveData(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar datos"})
				return
			}

			c.JSON(http.StatusOK, updatedUser)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
}

// Eliminar un usuario
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	for i, user := range models.Users {
		if user.ID == id {
			models.Users = append(models.Users[:i], models.Users[i+1:]...)

			// Guardar en archivo JSON
			if err := storage.SaveData(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar datos"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
}
```

Y nuestro archivo principal `main.go` se verá así:

```go
// main.go
package main

import (
	"log"
	"myapi/handlers"
	"myapi/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	// Cargar datos desde el archivo JSON al iniciar
	if err := storage.LoadData(); err != nil {
		log.Printf("Advertencia: No se pudieron cargar los datos: %v", err)
	}

	r := gin.Default()

	// Definir rutas para la API RESTful
	r.POST("/users", handlers.CreateUser)
	r.GET("/users", handlers.GetUsers)
	r.GET("/users/:id", handlers.GetUserByID)
	r.PUT("/users/:id", handlers.UpdateUser)
	r.DELETE("/users/:id", handlers.DeleteUser)

	r.Run() // Ejecutar el servidor en el puerto 8080 por defecto
}
```

Con esta estructura modular, nuestro proyecto es más organizado y fácil de mantener. Ahora, profundicemos en el uso de `gin.Context`.

---
## Manejo de Request Body
En Gin, el `gin.Context` proporciona métodos convenientes para manejar el cuerpo de las solicitudes HTTP. Uno de los métodos más comunes es `BindJSON`, que se utiliza para enlazar el cuerpo de una solicitud JSON a una estructura Go. Pero también existen otros métodos como `ShouldBindJSON`, `BindXML`, `ShouldBindXML`, entre otros, que permiten manejar diferentes formatos de datos.

`ShouldBindJSON` es similar a `BindJSON`, pero en lugar de devolver un error directamente, devuelve un error que puede ser manejado por el desarrollador. Tambien puedes usar etiquetas de validación en las estructuras para validar los datos entrantes automáticamente.

```go
// models/user.go
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// handlers/user_handlers.go
func CreateUser(c *gin.Context) {
    var newUser models.User
    if err := c.ShouldBindJSON(&newUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    // Lógica para crear el usuario...
}

```
Con `binding:"required"` y otras etiquetas, Gin validará automáticamente los campos cuando se llame a `ShouldBindJSON`. Si algún campo no cumple con las reglas de validación, se devolverá un error que podemos manejar y responder adecuadamente.

---
## Grupos de rutas
Primero, vamos hacer un amodificación en nuestro proyecto, crearemos otro archivo para registrar usuarios y hacer login. Crea un archivo llamado `auth_handlers.go` dentro de la carpeta `handlers`.

```go
// handlers/auth_handlers.go
package handlers

// Crear un nuevo usuario
func Register(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUser.ID = fmt.Sprintf("%d", models.NextID)
	models.NextID++
	models.Users = append(models.Users, newUser)

	// Guardar en archivo JSON
	if err := storage.SaveData(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar datos"})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

// Login de usuario
func Login(c *gin.Context) {
    var loginData struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required"`
    }
    if err := c.ShouldBindJSON(&loginData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    for _, user := range models.Users {
        if user.Email == loginData.Email && user.Password == loginData.Password {
            c.JSON(http.StatusOK, gin.H{"message": "login successful"})
            return
        }
    }
    c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
}
´´´
Ahora en el archivo `main.go`, registraremos estas nuevas rutas utilizando grupos de rutas:

```go
// main.go

func main() {
    
	// Rutas de autenticación
	r.POST("/auth/register", handlers.Register)
	r.POST("/auth/login", handlers.Login)
}
```
En Gin, los grupos de rutas permiten organizar y agrupar rutas relacionadas bajo un prefijo común. Esto es útil para estructurar aplicaciones más grandes y aplicar middlewares específicos a un conjunto de rutas.
Por ejemplo, podemos crear un grupo de rutas para todas las operaciones relacionadas con usuarios:

```go
	// Definir rutas para la API RESTful
	userGroup := r.Group("/users")
	{
		userGroup.POST("", handlers.CreateUser)
		userGroup.GET("", handlers.GetUsers)
		userGroup.GET("/:id", handlers.GetUserByID)
		userGroup.PUT("/:id", handlers.UpdateUser)
		userGroup.DELETE("/:id", handlers.DeleteUser)
	}

	// Rutas de autenticación
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", handlers.Register)
		authGroup.POST("/login", handlers.Login)
	} 
```

De esta manera, todas las rutas relacionadas con usuarios están agrupadas bajo el prefijo `/users`, y las rutas de autenticación bajo `/auth`. Esto facilita la gestión y el mantenimiento del código.



---
## Middlewares: concepto y uso
Los middlewares en Gin son funciones que se ejecutan antes o después de que una solicitud HTTP sea procesada por un manejador de rutas. Permiten realizar tareas comunes como autenticación, registro, manejo de errores, entre otros, de manera centralizada.
Un middleware recibe un `gin.Context` como parámetro y puede modificar la solicitud o la respuesta, o incluso abortar la solicitud si es necesario.
Por ejemplo, un middleware simple que registra el tiempo que tarda en procesarse una solicitud podría verse así:

```go
// handlers/middleware.go
func LoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        startTime := time.Now()
        log.Println("LoggerMiddleware está procesando la solicitud...")
        
        // Procesar la solicitud
        c.Next()
        
        duration := time.Since(startTime)
        log.Printf("Request %s %s took %v", c.Request.Method, c.Request.URL.Path, duration)
    }
}
```
- `time.Now()`: Captura el tiempo actual al inicio de la solicitud.
- `c.Next()`: Llama al siguiente middleware o al manejador de la ruta.
- Después de que la solicitud ha sido procesada, calcula la duración y registra el tiempo que tomó.
- `time.Since(startTime)`: Calcula la duración desde el inicio hasta el final de la solicitud.

Para usar este middleware, simplemente lo agregamos al router o a un grupo de rutas:

```go
// main.go
func ping(c *gin.Context) {
    c.JSON(200, gin.H{"message": "pong"})
}

    // Aplicar middleware de ping
	r.Use(handlers.LoggerMiddleware())
	r.GET("/ping", ping)
```

Este middleware registrará el tiempo de cada solicitud que pase por el router. Para probarlo, puedes hacer una solicitud GET a `/ping` y ver el registro en la consola.

---
## Creación de middleware personalizado
Crear un middleware personalizado en Gin es sencillo. Un middleware es simplemente una función que devuelve un `gin.HandlerFunc`. Esta función puede realizar operaciones antes y después de que la solicitud sea manejada por el controlador principal.

Aquí tienes un ejemplo de cómo crear un middleware personalizado que verifica si una solicitud tiene un encabezado específico, por ejemplo, un token de autenticación:

```go
// handlers/middleware.go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token != "mysecrettoken" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            return
        }
        c.Next() // Continuar con el siguiente middleware o controlador
    }
}
```
- `c.GetHeader("Authorization")`: Obtiene el valor del encabezado `Authorization` de la solicitud.
- Si el token no coincide con el valor esperado (`mysecrettoken`), se aborta la solicitud y se envía una respuesta 401 Unauthorized.
- `c.AbortWithStatusJSON(...)`: Aborta la solicitud y envía una respuesta JSON con el código de estado especificado.
- `c.Next()`: Si el token es válido, se llama a `c.Next()` para continuar con el siguiente middleware o el controlador de la ruta.

Para usar este middleware, simplemente lo agregamos al router o a un grupo de rutas, en este caso a todas las rutas de usuarios:

```go
	// Definir rutas para la API RESTful
	userGroup := r.Group("/users")
	userGroup.Use(handlers.AuthMiddleware())
	{
		userGroup.POST("", handlers.CreateUser)
		userGroup.GET("", handlers.GetUsers)
		userGroup.GET("/:id", handlers.GetUserByID)
		userGroup.PUT("/:id", handlers.UpdateUser)
		userGroup.DELETE("/:id", handlers.DeleteUser)
	}
```

Si es a una sola ruta, lo podemos hacer de pin:

```go
r.GET("/ping", handlers.AuthMiddleware(), ping)
```

Para probar este middleware, puedes hacer una solicitud GET a `/users` sin el encabezado `Authorization` o con un valor incorrecto, y deberías recibir una respuesta 401 Unauthorized. Si proporcionas el encabezado correcto, la solicitud debería pasar y obtener la lista de usuarios.


---
## Manejo de CORS
Los CORS (Cross-Origin Resource Sharing) son un mecanismo que permite controlar cómo los recursos en un servidor web pueden ser solicitados desde otro dominio distinto al del servidor. En aplicaciones web modernas, es común que el frontend y el backend estén alojados en dominios diferentes, por lo que es esencial configurar CORS adecuadamente para permitir estas solicitudes cruzadas.
En Gin, podemos manejar CORS utilizando un middleware. A continuación, te muestro cómo crear un middleware para manejar CORS:

```go
// handlers/middleware.go
import "github.com/gin-contrib/cors"

func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://example.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	})
}
```
- `AllowOrigins`: Especifica los orígenes permitidos para realizar solicitudes. Puedes agregar los dominios de tu frontend aquí.
- `AllowMethods`: Define los métodos HTTP permitidos para las solicitudes.
- `AllowHeaders`: Especifica los encabezados permitidos en las solicitudes.
- `ExposeHeaders`: Define los encabezados que pueden ser expuestos al navegador.
- `AllowCredentials`: Indica si se permiten las credenciales (cookies, encabezados de autorización) en las solicitudes.
- `MaxAge`: Define el tiempo que el navegador puede almacenar en caché la respuesta de preflight.

Luego, para usar este middleware, simplemente lo agregamos al router principal en `main.go`:

```go
// main.go
func main() {
	r := gin.Default()
	// Aplicar middleware de CORS
	r.Use(handlers.CORSMiddleware())
	// Definir rutas para la API RESTful
	// ...
	r.Run() // Ejecutar el servidor en el puerto 8080 por defecto
}
```
Con esta configuración, tu servidor Gin permitirá solicitudes CORS desde los orígenes especificados, facilitando la comunicación entre tu frontend y backend.

---
## Manejo de Headers
Los encabezados HTTP (Headers) son componentes clave en las solicitudes y respuestas HTTP. Permiten enviar información adicional entre el cliente y el servidor, como detalles sobre el tipo de contenido, autenticación, control de caché, entre otros. En Gin, podemos manejar los encabezados de manera sencilla utilizando el `gin.Context`. A continuación, te muestro cómo trabajar con los encabezados en Gin:

```go
// handlers/header_handlers.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Obtener un encabezado específico
func GetCustomHeader(c *gin.Context) {
	headerValue := c.GetHeader("X-Custom-Header")
	if headerValue == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Custom-Header not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"X-Custom-Header": headerValue})
}

// Establecer un encabezado en la respuesta
func SetCustomHeader(c *gin.Context) {
	c.Header("X-Custom-Response-Header", "CustomValue")
	c.JSON(http.StatusOK, gin.H{"message": "Custom header set"})
}
```
- `c.GetHeader("X-Custom-Header")`: Obtiene el valor del encabezado `X-Custom-Header` de la solicitud. Si no se encuentra, devuelve una respuesta de error.
- `c.Header("X-Custom-Response-Header", "CustomValue")`: Establece un encabezado `X-Custom-Response-Header` en la respuesta con el valor `CustomValue`.

Luego, en el archivo `main.go`, registramos estas nuevas rutas para manejar los encabezados:

```go
// main.go
func main() {
	r := gin.Default()
	// Rutas para manejar encabezados
	r.GET("/get-header", handlers.GetCustomHeader)
	r.GET("/set-header", handlers.SetCustomHeader)
	// Definir otras rutas para la API RESTful
	r.Run() // Ejecutar el servidor en el puerto 8080 por defecto
}
```

Con estas rutas, puedes probar lo siguiente:
- Hacer una solicitud GET a `/get-header` con un encabezado `X-Custom-Header` para obtener su valor.
- Hacer una solicitud GET a `/set-header` para establecer un encabezado personalizado en la respuesta.

El manejo adecuado de los encabezados HTTP es esencial para la comunicación efectiva entre el cliente y el servidor, y Gin facilita este proceso con su API intuitiva.

---
## Manejo global de errores

En Gin, podemos manejar errores de manera global utilizando middlewares. Esto nos permite capturar cualquier error que ocurra durante el procesamiento de una solicitud y responder de manera consistente. A continuación, te muestro cómo implementar un middleware para el manejo global de errores:

```go
// handlers/middleware.go
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Procesar la solicitud

		// Verificar si hubo errores durante el procesamiento
		if len(c.Errors) > 0 {
			// Aquí puedes personalizar la respuesta de error
			c.JSON(-1, gin.H{"errors": c.Errors.JSON()})
		}
	}
}
```
- `c.Next()`: Procesa la solicitud y llama al siguiente middleware o controlador.
- `c.Errors`: Contiene una lista de errores que ocurrieron durante el procesamiento de la solicitud.
- Si hay errores, se envía una respuesta JSON con los detalles de los errores.
Luego, para usar este middleware, simplemente lo agregamos al router principal en `main.go`:

```go
// main.go
func main() {
	r := gin.Default()
	// Aplicar middleware de manejo global de errores
	r.Use(handlers.ErrorHandlingMiddleware())
	// Definir rutas para la API RESTful
	// ...
	r.Run() // Ejecutar el servidor en el puerto 8080 por defecto
}
```

Con esta configuración, cualquier error que ocurra durante el procesamiento de una solicitud será capturado por el middleware y se enviará una respuesta JSON con los detalles del error. Esto facilita el manejo de errores y proporciona una experiencia consistente para los clientes de la API.

---
## Buenas prácticas con Gin
Al trabajar con el framework Gin en Go, es importante seguir ciertas buenas prácticas para garantizar que nuestras aplicaciones sean eficientes, mantenibles y seguras. A continuación, se presentan algunas recomendaciones clave:

1. **Estructura del Proyecto**: Organiza tu proyecto en carpetas claras como `handlers`, `models`, `storage`, y `middlewares`. Esto facilita la navegación y el mantenimiento del código.
2. **Uso de Middlewares**: Aprovecha los middlewares para manejar tareas comunes como autenticación, registro, manejo de errores y CORS. Esto ayuda a mantener el código limpio y modular.
3. **Validación de Datos**: Siempre valida los datos entrantes utilizando binding y etiquetas de validación. Considera el uso de validaciones personalizadas para requisitos específicos.
4. **Manejo de Errores**: Implementa un manejo global de errores para capturar y responder de manera consistente a los errores que ocurran durante el procesamiento de solicitudes.
5. **Documentación de la API**: Utiliza herramientas como Swagger para documentar tu API. Esto facilita la comprensión y el uso de la API por parte de otros desarrolladores.
6. **Pruebas Unitarias**: Escribe pruebas unitarias para tus handlers y middlewares. Esto asegura que tu código funcione correctamente y facilita la detección de errores.
7. **Seguridad**: Asegúrate de implementar medidas de seguridad adecuadas, como la validación de entradas, el uso de HTTPS, y la protección contra ataques comunes como inyección SQL y XSS.
8. **Manejo de Configuraciones**: Utiliza archivos de configuración o variables de entorno para gestionar configuraciones sensibles como credenciales de base de datos y claves API.
9. **Optimización del Rendimiento**: Monitorea y optimiza el rendimiento de tu aplicación utilizando herramientas de profiling y logging para identificar cuellos de botella.
10. **Actualización Regular**: Mantén tus dependencias actualizadas para beneficiarte de las últimas mejoras y correcciones de seguridad.
Siguiendo estas buenas prácticas, podrás desarrollar aplicaciones web robustas y eficientes utilizando el framework Gin en Go.

---
## Conclusión
En esta sección de Profundizando en Gin, hemos explorado diversas características avanzadas del framework Gin para construir aplicaciones web robustas y eficientes en Go. Hemos aprendido sobre el manejo de rutas, la creación y uso de middlewares, la validación de datos, el manejo global de errores y las buenas prácticas para trabajar con Gin. Al aplicar estos conceptos y técnicas, estarás mejor equipado para desarrollar aplicaciones web escalables y mantenibles utilizando Gin. ¡Sigue practicando y explorando más funcionalidades para convertirte en un experto en Gin!