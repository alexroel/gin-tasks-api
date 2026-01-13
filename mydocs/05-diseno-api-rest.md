# Diseño de APIs REST


---
## Introducción


---
## ¿Qué es una API REST?
Una API REST (Representational State Transfer) es un conjunto de reglas y convenciones para construir servicios web que permiten la comunicación entre sistemas a través de HTTP. Las APIs RESTful utilizan métodos HTTP estándar como GET, POST, PUT, DELETE, entre otros, para realizar operaciones sobre recursos representados en formato JSON o XML.

**Principios de REST:**
1. **Cliente-Servidor:** La arquitectura REST separa la interfaz de usuario del almacenamiento de datos, lo que permite una mayor flexibilidad y escalabilidad.
2. **Sin estado:** Cada solicitud del cliente al servidor debe contener toda la información necesaria para entender y procesar la solicitud. El servidor no debe almacenar ningún estado del cliente entre solicitudes.
3. **Caché:** Las respuestas deben ser etiquetadas como cachéables o no cachéables para mejorar el rendimiento.
4. **Interfaz uniforme:** REST define una interfaz uniforme entre los componentes, lo que simplifica la arquitectura y mejora la visibilidad de las interacciones.
5. **Sistema en capas:** La arquitectura REST puede estar compuesta por múltiples capas, lo que permite la escalabilidad y la separación de responsabilidades.
6. **Código bajo demanda (opcional):** Los servidores pueden enviar código ejecutable al cliente para extender su funcionalidad.

---
## Presentación del proyecto de API REST
En este proyecto, construiremos una API REST básica utilizando el framework Gin en Go. La API permitirá gestionar recursos de tareas (tasks) y usuarios (users) con operaciones CRUD (Crear, Leer, Actualizar, Eliminar).

**Objetivos del proyecto:**
- Implementar endpoints RESTful para gestionar tareas y usuarios.
- Utilizar métodos HTTP adecuados para cada operación.
- Manejar respuestas y errores de manera estándar.
- Aplicar buenas prácticas en el diseño de la API.

