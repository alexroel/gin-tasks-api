# Introducción a Gin Framework

1. [Introducción](#introducción)
2. [¿Qué es Gin y cómo funciona?](#qué-es-gin-y-cómo-funciona)
3. [Gin vs net/http](#gin-vs-nethttp)
4. [Primer servidor HTTP con Gin](#primer-servidor-http-con-gin)
5. [Rutas básicas](#rutas-básicas)
6. [Uso de parámetros de ruta](#uso-de-parámetros-de-ruta)
7. [Query Strings](#query-strings)
8. [Petición y Respuestas JSON](#petición-y-respuestas-json)
9. [Manejo básico de errores HTTP](#manejo-básico-de-errores-http)
10. [Mini API REST](#mini-api-rest)
11. [Terminar CRUD con IA](#terminar-crud-con-ia)
12. [Conclusión](#conclusión)


---
## Introducción
Bienvenidos a la sección de introducción al framework Gin. En esta sección, aprenderemos los conceptos básicos de Gin, un framework web ligero y rápido para Go. Cubriremos desde la creación de un servidor HTTP básico hasta la construcción de una mini API RESTful.

---
## ¿Qué es Gin y cómo funciona?
Gin es un framework web escrito en Go que se destaca por su rendimiento y facilidad de uso. Utiliza un enrutador basado en árboles para manejar las rutas de manera eficiente y proporciona una serie de características útiles como middleware, validación de datos y manejo de errores.

**Características principales de Gin:**
- Alto rendimiento
- Enrutamiento rápido basado en árboles
- Middleware integrado
- Soporte para JSON y otros formatos de respuesta
- Manejo de errores simplificado 

**Cómo funciona Gin:**
Gin utiliza un enrutador que mapea las rutas HTTP a funciones manejadoras. Cuando una solicitud llega al servidor, Gin busca la ruta correspondiente y ejecuta la función asociada, pasando un contexto que contiene información sobre la solicitud y la respuesta.

**Qué necesitas para empezar:**
- Conocimientos básicos de Go
- Tener Go instalado en tu máquina

**Qué pudes hacer con Gin:**
- Crear servidores web rápidos y eficientes
- Construir APIs RESTful
- Manejar solicitudes y respuestas HTTP de manera sencilla
- Implementar middleware para funcionalidades adicionales
- Gestionar errores de forma efectiva
- Microservicios escalables 
- Integrar con bases de datos y otros servicios



---
## Gin vs net/http
Gin es un framework web construido sobre el paquete net/http de Go, pero ofrece varias ventajas que facilitan el desarrollo de aplicaciones web. A continuación, se presentan algunas diferencias clave entre Gin y net/http:

| Característica          | net/http                          | Gin                                 |
|------------------------|----------------------------------|------------------------------------|
| Facilidad de uso       | Requiere más configuración       | Proporciona una API sencilla y amigable |
| Enrutamiento           | Básico, basado en patrones       | Enrutamiento basado en árboles, más eficiente |
| Middleware             | No integrado, requiere implementación manual | | Middleware integrado y fácil de usar |
| Manejo de errores      | Manual                           | Manejo de errores simplificado |
| Rendimiento            | Bueno                            | Muy alto, optimizado para velocidad |
| Soporte para JSON      | Manual                           | Soporte nativo para JSON y otros formatos |

**Cuándo usar Gin sobre net/http:**
- Cuando se necesita un desarrollo rápido y eficiente de aplicaciones web.
- Cuando se requiere un enrutamiento avanzado y manejo de middleware.
- Cuando se busca un framework con características integradas para facilitar el desarrollo de APIs RESTful.

*Página oficial de Gin:* [https://gin-gonic.com/](https://gin-gonic.com/)
*Documentación de net/http:* [https://pkg.go.dev/net/http](https://pkg.go.dev/net/http)

Veras que Gin simplifica muchas tareas comunes en el desarrollo web, permitiéndote concentrarte en la lógica de tu aplicación en lugar de en los detalles de bajo nivel del manejo de solicitudes HTTP.
---
## Primer servidor HTTP con Gin
Ahora, vamos a crear nuestro primer servidor HTTP utilizando el framework Gin. Sigue los pasos a continuación para configurar y ejecutar un servidor básico que responda a las solicitudes GET.

1. **Crear un nuevo proyecto:**
   Abre tu terminal y crea un nuevo directorio para tu proyecto. Luego, navega a ese directorio.

   ```bash
   mkdir gin-users-api
   cd gin-users-api
   ```

2. **Inicializar un módulo Go:**
   Ejecuta el siguiente comando para inicializar un nuevo módulo Go:
    ```bash
    go mod init gin-users-api
    ```
3. **Instalar Gin:**
    Usa el siguiente comando para instalar el paquete Gin:
    ```bash
    go get -u github.com/gin-gonic/gin
    ```

4. **Crear el archivo principal:**
   Crea un archivo llamado `main.go` en el directorio del proyecto y abrelo en tu editor de código favorito.

5. **Escribir el código del servidor:**
    Copia y pega el siguiente código en `main.go`:
    ```go
    package main

    import (
        "github.com/gin-gonic/gin"
    )

    func main() {
        // Crear una instancia del router Gin
        r := gin.Default()

        // Definir una ruta GET para la raíz
        r.GET("/", func(c *gin.Context) {
            c.String(200, "¡Hola, Mundo desde Gin!")
        })

        // Ejecutar el servidor en el puerto 8080 por defecto
        r.Run()
    }
    ```

6. **Ejecutar el servidor:**
   Guarda el archivo y vuelve a la terminal. Ejecuta el siguiente comando para iniciar el servidor:
   ```bash
   go run main.go
   ```
7. **Probar el servidor:**
   Abre tu navegador web o una herramienta como Postman y navega a `http://localhost:8080`. Deberías ver el mensaje "¡Hola, Mundo desde Gin!".

**explicación del código:**
- `gin.Default()`: Crea una instancia del router Gin con middleware predeterminado (logger y recuperación).
- `r.GET("/", ...)`: Define una ruta GET para la raíz (`/`) que responde con un mensaje de texto.
- `c.String(200, ...)`: Envía una respuesta de texto con el código de estado HTTP 200.
- `r.Run()`: Inicia el servidor en el puerto 8080.
- `logger y recuperación`: El middleware predeterminado incluye un logger que registra las solicitudes entrantes y un middleware de recuperación que maneja los pánicos y evita que el servidor se caiga.

¡Felicidades! Has creado y ejecutado tu primer servidor HTTP con Gin. En las siguientes secciones, exploraremos más características de Gin para construir aplicaciones web más complejas.

---
## Rutas básicas

En Gin, las rutas son la forma en que defines cómo tu aplicación responde a las solicitudes HTTP. Puedes definir rutas para diferentes métodos HTTP como GET, POST, PUT, DELETE, entre otros. A continuación, te mostramos cómo crear rutas básicas en Gin.

```go
package main
import (
    "github.com/gin-gonic/gin"
)

func CrateUser(c *gin.Context) {
    c.String(200, "Crear un nuevo usuario")
}

func GetUser(c *gin.Context) {
    c.String(200, "Obtener información del usuario")
}


func UpdateUser(c *gin.Context) {
    c.String(200, "Actualizar usuario con ID: ")
}

func DeleteUser(c *gin.Context) {
    c.String(200, "Eliminar usuario con ID: ")
}

func main() {
    r := gin.Default()

    // Definir rutas básicas
    r.GET("/users", GetUser)         // Obtener un usuario por ID
    r.POST("/users", CrateUser)         // Crear un nuevo usuario
    r.PUT("/users", UpdateUser)     // Actualizar un usuario por ID
    r.DELETE("/users", DeleteUser)  // Eliminar un usuario por ID

    r.Run() // Ejecutar el servidor en el puerto 8080 por defecto
}
```
Para probar estas rutas, puedes usar herramientas como Postman o cURL para enviar solicitudes HTTP a tu servidor Gin. Cada ruta está asociada con una función manejadora que define la lógica para esa ruta específica.


---
## Uso de parámetros de ruta
En Gin, puedes definir parámetros de ruta para capturar valores dinámicos en las rutas de tu aplicación. Los parámetros de ruta se definen utilizando dos puntos (`:`) seguidos del nombre del parámetro en la ruta. A continuación, te mostramos cómo usar parámetros de ruta en Gin.

```go
func GetUserByID(c *gin.Context) {
    // Obtener el valor del parámetro de ruta "id"
    id := c.Param("id")
    c.String(200, "Obtener usuario con ID: %s", id)
}

func UpdateUser(c *gin.Context) {
    // Obtener el valor del parámetro de ruta "id"
    id := c.Param("id")
    c.String(200, "Actualizar usuario con ID: %s", id)
}

func DeleteUser(c *gin.Context) {
    // Obtener el valor del parámetro de ruta "id"
    id := c.Param("id")
    c.String(200, "Eliminar usuario con ID: %s", id)
}

func main() {
    r := gin.Default()

    // Definir una ruta con un parámetro de ruta
    r.GET("/users/:id", GetUserByID) // Obtener un usuario por ID
    r.PUT("/users/:id", UpdateUser)     // Actualizar un usuario por ID
    r.DELETE("/users/:id", DeleteUser)  // Eliminar un usuario por ID

    r.Run() // Ejecutar el servidor en el puerto 8080 por defecto
}
```

Lo mismo puedes hacer con otros métodos HTTP como PUT o DELETE, utilizando parámetros de ruta para identificar recursos específicos.


---
## Query Strings
En Gin, puedes manejar query strings (cadenas de consulta) para capturar parámetros adicionales enviados en la URL de una solicitud HTTP. Los query strings se encuentran después del signo de interrogación (`?`) en la URL y consisten en pares clave-valor separados por `&`. A continuación, te mostramos cómo manejar query strings en Gin.

```go
// Obtener usuarios por email
func GetUserByEmail(c *gin.Context) {
    // Obtener el valor del query string "email"
    email := c.Query("email")
    if email == "" {
        c.String(400, "El parámetro 'email' es obligatorio")
        return
    }
    c.String(200, "Obtener usuario con email: %s", email)
}
func main() {
    r := gin.Default()

    // Definir una ruta que maneje query strings
    r.GET("/users", GetUserByEmail) // Obtener usuario por email

    r.Run() // Ejecutar el servidor en el puerto 8080 por defecto
}
```

Para probar esta ruta, puedes enviar una solicitud GET a `http://localhost:808/users?email=roelcode@ejemplo.com`. La función `GetUserByEmail` capturará el valor del query string `email` y responderá en consecuencia.

---
## Petición y Respuestas JSON
En Gin, puedes manejar fácilmente peticiones y respuestas en formato JSON. Gin proporciona métodos integrados para serializar y deserializar datos JSON, lo que facilita el trabajo con APIs RESTful. A continuación, te mostramos cómo manejar peticiones y respuestas JSON en Gin.

```go
import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type User struct {
    ID     string  `json:"id"`
    Name  string  `json:"name"`
    Email string  `json:"email"`
    Password string  `json:"password"`
}

// Manejar una petición POST para crear un nuevo usuario
func CreateUser(c *gin.Context) {
    var newUser User
    // Deserializar el cuerpo de la solicitud JSON en la estructura User
    if err := c.BindJSON(&newUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    // Simular la creación del usuario (aquí podrías guardar en una base de datos)
    c.JSON(http.StatusCreated, gin.H{"message": "Usuario creado", "user": newUser})
}

func main() {
    r := gin.Default()

    // Definir una ruta POST para crear un nuevo usuario
    r.POST("/users", CreateUser) // Crear un nuevo usuario

    r.Run() // Ejecutar el servidor en el puerto 8080 por defecto
}
```

- `c.BindJSON(&newUser)`: Deserializa el cuerpo JSON de la solicitud en la estructura `User`.
- `gin.H{"error": err.Error()}`: Crea un mapa para enviar mensajes de error en formato JSON.
- `c.JSON(http.StatusCreated, ...)`: Envía una respuesta JSON con el código de estado HTTP 201 (Creado).

---
## Manejo básico de errores HTTP
En Gin, puedes manejar errores HTTP de manera sencilla utilizando el contexto (`gin.Context`). Puedes enviar respuestas de error con códigos de estado HTTP apropiados y mensajes descriptivos. A continuación, te mostramos cómo manejar errores HTTP básicos en Gin.

```go
import (
    "net/http"
    "github.com/gin-gonic/gin"
)

// Manejar una petición GET para obtener un usuario por ID
func GetAlbumByID(c *gin.Context) {
    id := c.Param("id")
    // Simular la búsqueda de un usuario (aquí podrías consultar una base de datos)
    if id != "1" {
        // Si el usuario no se encuentra, enviar un error 404
        c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
        return
    }
    // Si el usuario se encuentra, enviar la información del usuario
    user := gin.H{
        "id":    "1",
        "name":  "Roelcode",
        "email": "roelcode@ejemplo.com",
    }
    c.JSON(http.StatusOK, user)
}

func main() {
    r := gin.Default()

    // Definir una ruta GET para obtener un usuario por ID
    r.GET("/users/:id", GetAlbumByID) // Obtener un usuario por ID

    r.Run() // Ejecutar el servidor en el puerto 8080 por defecto
}
```
- Si el álbum no se encuentra, se envía una respuesta JSON con el código de estado HTTP 404 (No Encontrado) y un mensaje de error.

---
## Mini API REST
En esta sección, construiremos una mini API RESTful utilizando Gin. Nuestra API gestionará una colección de álbumes de música, permitiendo operaciones CRUD (Crear, Leer, Actualizar, Eliminar). A continuación, te mostramos cómo implementar esta API.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Estructura para representar un usuario
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Colección de Usuarios en memoria
var users []User
var nextID = 1

// Crear un nuevo usuario
func CreateUser(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUser.ID = fmt.Sprintf("%d", nextID)
	nextID++
	users = append(users, newUser)
	c.JSON(http.StatusCreated, newUser)
}

// Obtener todos los usuarios
func GetUsers(c *gin.Context) {
    if len(users) == 0 {
        c.JSON(http.StatusOK, gin.H{"message": "No hay usuarios disponibles"})
        return
    }
	c.JSON(http.StatusOK, users)
}

// Obtener un usuario por ID
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	for _, user := range users {
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
	c.String(200, "Actualizar usuario con ID: "+ id)
}

// Eliminar un usuario
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	c.String(200, "Eliminar usuario con ID: "+ id)
}

func main() {
	r := gin.Default()

	// Definir rutas para la API RESTful
	r.POST("/users", CreateUser)
	r.GET("/users", GetUsers)
	r.GET("/users/:id", GetUserByID)
	r.PUT("/users/:id", UpdateUser)
	r.DELETE("/users/:id", DeleteUser)

	r.Run() // Ejecutar el servidor en el puerto 8080 por defecto
}
```
---
## Terminar CRUD con IA
Ahoara que hemos implementado las operaciones de creación y lectura en nuestra mini API RESTful, vamos a completar las operaciones de actualización y eliminación (CRUD) utilizando Gin. A continuación, te mostramos cómo implementar estas funcionalidades.

Para terminar el CRUD, añadiremos dos nuevas funciones: `UpdateUser` y `DeleteUser` le pediremos a la IA que nos ayude con el siguiente pront:
```
Explora el archivo main.go de la mini API RESTful que hemos creado con Gin. y completa las funciones UpdateUser y DeleteUser para permitir la actualización y eliminación de usuarios en la colección. Asegúrate de manejar los casos en los que el usuario no se encuentre y de enviar respuestas JSON adecuadas.
```

```go
// Actualizar un usuario existente
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var updatedUser User
	if err := c.BindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, user := range users {
		if user.ID == id {
			updatedUser.ID = id
			users[i] = updatedUser
			c.JSON(http.StatusOK, updatedUser)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
}

// Eliminar un usuario
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
}
```

¡Felicidades! Ahora has completado las operaciones CRUD en tu mini API RESTful utilizando Gin. Puedes probar las rutas de actualización y eliminación utilizando herramientas como Postman o cURL para enviar solicitudes HTTP a tu servidor Gin.

---
## Conclusión
En esta sección, hemos aprendido los conceptos básicos del framework Gin y cómo construir una mini API RESTful. Hemos cubierto desde la creación de un servidor HTTP básico hasta la implementación de operaciones CRUD completas. Gin facilita el desarrollo de aplicaciones web en Go gracias a su rendimiento, facilidad de uso y características integradas.A medida que avances en tu aprendizaje de Gin, podrás explorar características más avanzadas como middleware, validación de datos y manejo global de errores. ¡Sigue practicando y construyendo aplicaciones web con Gin! 