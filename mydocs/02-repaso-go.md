# Fundamentos del Lenguaje Go - REPASO

1. [Introducción](#introducción)
2. [¿Qué es Go y por qué usarlo para Backend?](#qué-es-go-y-por-qué-usarlo-para-backend)
3. [Crear un proyecto en Go](#crear-un-proyecto-en-go)
4. [Variables, constantes y tipos de datos](#variables-constantes-y-tipos-de-datos)
5. [Control de flujo](#control-de-flujo)
6. [Funciones](#funciones)
7. [Structs y métodos](#structs-y-métodos)
8. [Punteros](#punteros)
9. [Interfaces](#interfaces)
10. [Manejo de errores](#manejo-de-errores)
11. [Paquetes y Módulos](#paquetes-y-módulos)
12. [JSON y Encoding](#json-y-encoding)
13. [Concurrencia](#concurrencia)
14. [Context (context.Context)](#context-contextcontext)
15. [Testing en Go](#testing-en-go)
16. [HTTP y net/http](#http-y-nethttp)
17. [Mini API RESTful](#mini-api-restful)
18. [Conclusión](#conclusión)


---
## Introducción
Go, también conocido como Golang, es un lenguaje de programación de código abierto desarrollado por Google. Fue diseñado para ser simple, eficiente y fácil de usar, con un enfoque en la concurrencia y el rendimiento. Go es especialmente popular para el desarrollo de aplicaciones backend debido a su capacidad para manejar múltiples tareas simultáneamente y su rendimiento comparable al de lenguajes compilados como C o C++.


---
## ¿Qué es Go y por qué usarlo para Backend?
Go es un lenguaje de programación desarrollado por Google que se destaca por su simplicidad, eficiencia y capacidad para manejar la concurrencia de manera efectiva. Fue diseñado para abordar las limitaciones de otros lenguajes en términos de rendimiento y escalabilidad, especialmente en el desarrollo de aplicaciones backend.

**Características clave de Go:**
- **Simplicidad:** Go tiene una sintaxis clara y concisa, lo que facilita su aprendizaje y uso. La simplicidad del lenguaje permite a los desarrolladores escribir código limpio y mantenible.
- **Rendimiento:** Go es un lenguaje compilado, lo que significa que el código se traduce directamente a código máquina. Esto resulta en un rendimiento rápido y eficiente, comparable al de lenguajes como C o C++.
- **Concurrencia:** Go tiene un modelo de concurrencia incorporado basado en goroutines y canales, lo que facilita la creación de aplicaciones que pueden manejar múltiples tareas simultáneamente sin complicaciones.
- **Gestión de dependencias:** Go utiliza un sistema de módulos que facilita la gestión de dependencias y la distribución de paquetes, lo que es esencial para proyectos grandes y colaborativos.
- **Ecosistema robusto:** Go cuenta con una amplia gama de bibliotecas y frameworks, como Gin, que facilitan el desarrollo de aplicaciones web y APIs RESTful.

**¿Por qué usar Go para Backend?**
- **Escalabilidad:** La capacidad de Go para manejar la concurrencia lo hace ideal para aplicaciones backend que requieren escalabilidad y rendimiento bajo carga. 
- **Desarrollo rápido:** La simplicidad del lenguaje y la disponibilidad de herramientas y bibliotecas permiten a los desarrolladores construir aplicaciones backend de manera rápida y eficiente.
- **Mantenimiento:** El código escrito en Go tiende a ser más fácil de mantener debido a su claridad y estructura, lo que reduce la deuda técnica a largo plazo.

---
## Crear un proyecto en Go

Para crear un proyecto básico en Go, sigue estos pasos:

```bash 
mkdir mi-proyecto-go
cd mi-proyecto-go
touch main.go
```
Dentro del archivo `main.go`, puedes escribir un programa simple como este:

```go
package main

import "fmt"

func main() {
	fmt.Println("¡Hola, Mundo!")
}
```

Para ejecutar el programa, usa el siguiente comando en la terminal:

```bash
go run main.go
```

- `main.go`: Es el archivo que contiene el código fuente de tu programa.
- `package main`: Define el paquete principal del programa.
- `import "fmt"`: Importa el paquete `fmt`, que contiene funciones para formatear y imprimir texto.
- `func main() { ... }`: Define la función principal donde comienza la ejecución del programa
- `fmt.Println`: Esta función imprime el texto en la consola.
- `fmt`: Es un paquete estándar de Go que proporciona funciones para imprimir texto, leer entrada, y formatear cadenas.
- `go run`: Este comando compila y ejecuta el archivo Go especificado.

```bash
go build main.go
./main
```

- `go build`: Compila el archivo Go y genera un ejecutable.
- `./main`: Ejecuta el archivo compilado en sistemas Unix/Linux. En Windows, usarías `main.exe`.

```bash
go build -o myapp main.go
./myapp
```

- `-o myapp`: Especifica el nombre del archivo ejecutable generado.
- `./myapp`: Ejecuta el archivo compilado con el nombre personalizado.
- `go build`: Se utiliza para compilar el código fuente de Go en un archivo ejecutable. Es útil para distribuir aplicaciones Go sin necesidad de recompilar el código fuente cada vez que se quiera ejecutar.

```bash
go mod init myproject
```
- `go mod init myproject`: Inicializa un nuevo módulo de Go llamado `myproject`, creando un archivo `go.mod` que gestiona las dependencias del proyecto.    
- `go.mod`: Es un archivo que contiene información sobre el módulo, incluyendo su nombre y las dependencias necesarias para el proyecto.
- `go.sum`: Es un archivo que contiene sumas de verificación criptográficas para las dependencias del módulo, asegurando la integridad y seguridad de las mismas.

**Instalación de Paquetes Externos**
Para instalar paquetes externos en tu proyecto Go, puedes usar el comando `go get`. Por ejemplo, para instalar el paquete `alexroel/gosaludos`, que es una paquete simple para saludar, con saudos aleatorios, puedes ejecutar:

```bash
go get -u github.com/alexroel/gosaludos@latest
```

- `go get -u`: Descarga e instala la última versión del paquete especificado, actualizando las dependencias si es necesario.
- `github.com/alexroel/gosaludos@latest`: Especifica la URL del paquete y la versión que deseas instalar (en este caso, la última versión disponible).
- `import "github.com/alexroel/gosaludos"`: Importa el paquete instalado para que puedas usar sus funciones en tu código.

```go
package main

import (
    "fmt"
    "github.com/alexroel/gosaludos"
)

func main() {
    mensaje := gosaludos.SaludaA("Alex")
	fmt.Println(mensaje)
}
```

- `mensaje := gosaludos.SaludaA("Alex")`: Llama a la función `SaludaA` del paquete `gosaludos`, pasando "Alex" como argumento, y almacena el resultado en la variable `mensaje`.
- `fmt.Println(mensaje)`: Imprime el mensaje de saludo generado por la función `


---
## Variables, constantes y tipos de datos

En Go, las variables y constantes son fundamentales para almacenar y manipular datos. A continuación, se describen cómo declararlas y los tipos de datos básicos disponibles en Go.

**Variables:**
```go
var nombre string
nombre = "Juan"
var edad int = 30
altura := 1.70 // Declaración corta y asignación
```

- `var nombre string`: Declara una variable llamada `nombre` de tipo `string`.
- `nombre = "Juan"`: Asigna el valor "Juan" a la variable `nombre`.
- `var edad int = 30`: Declara una variable `edad` de tipo `int` y la inicializa con el valor 30.
- `altura := 1.70`: Utiliza la declaración corta para crear e inicializar la variable `altura` con el valor 1.70 (tipo `float64`).
- `:=`: Es el operador de declaración corta que permite declarar e inicializar una variable en una sola línea sin especificar el tipo explícitamente.

**Constantes:**
```go
const pi float64 = 3.14
const saludo = "Hola, Mundo!"
```

- `const pi float64 = 3.14`: Declara una constante llamada `pi` de tipo `float64` con el valor 3.14.
- `const saludo = "Hola, Mundo!"`: Declara una constante `saludo` de tipo `string` con el valor "Hola, Mundo!".
- `const`: Es la palabra clave utilizada para declarar constantes en Go. Las constantes son valores que no pueden cambiar durante la ejecución del programa.

**Tipos de Datos Básicos:**
- `int`: Enteros (números sin decimales).
- `float64`: Números de punto flotante (números con decimales).
- `string`: Cadenas de texto.
- `bool`: Valores booleanos (`true` o `false`).  

**Más Tipados de Datos:**
- `byte`: Representa un solo byte (alias de `uint8`).
- `rune`: Representa un punto de código Unicode (alias de `int32`).
- `int8`, `int16`, `int32`, `int64`: Enteros con diferentes tamaños.
- `uint8`, `uint16`, `uint32`, `uint64`: Enteros sin signo con diferentes tamaños.
- `float32`: Números de punto flotante de menor precisión.
- `complex64`, `complex128`: Números complejos. 

**Conversión de Tipos:**
```go
var edadInt int = 25 
var edadFloat float64 = float64(edadInt)

// Conversión de string a int
import "strconv"
...

edadStr := "30"
edadConvertida, err := strconv.Atoi(edadStr)

// Conversión de int a string
edadStr2 := strconv.Itoa(edadInt)
```

- `float64(edadInt)`: Convierte la variable `edadInt` de tipo `int` a `float64`.
- `strconv.Atoi(edadStr)`: Convierte una cadena de texto `edadStr` a un entero. Devuelve el entero convertido y un error si la conversión falla.
- `strconv.Itoa(edadInt)`: Convierte un entero `edadInt` a una cadena de texto.

**Declaración Múltiple:**
```go
var a, b, c int = 1, 2, 3
d, e, f := "Hola", 3.14, true

var (
    x int = 10
    y string = "GoLang"
    z bool = false
)
```

- `var a, b, c int = 1, 2, 3`: Declara múltiples variables `a`, `b` y `c` de tipo `int` e inicializa con valores.
- `d, e, f := "Hola", 3.14, true`: Declara e inicializa múltiples variables con diferentes tipos utilizando la declaración corta.
- `var ( ... )`: Permite declarar múltiples variables en un bloque para mejorar la legibilidad del código.

**Values Cero:**
En Go, las variables no inicializadas tienen un "valor cero" predeterminado según su tipo:
- `int`: 0
- `float64`: 0.0
- `string`: ""
- `bool`: false

**Tipos Compuestos:**
- `arrays`: Colecciones de elementos del mismo tipo con tamaño fijo.
```go 
var edades [3] int
edades[0] = 25

fmt.Println(edades)

var nombres = [3]string{"Ana", "Luis", "Marta"}
fmt.Println(nombres)
```

- `slices`: Colecciones dinámicas que pueden crecer y reducirse en tamaño.
```go
var frutas []string
frutas = append(frutas, "Manzana")
frutas = append(frutas, "Banana", "Cereza")
fmt.Println(frutas)
```

- `maps`: Colecciones de pares clave-valor.
```go
var persona map[string]int
persona = make(map[string]int)
persona["edad"] = 30
persona["altura"] = 175
fmt.Println(persona)
```

- Aparte de `arrays`, `slices` y `maps`, Go también soporta otros tipos compuestos como `structs` (estructuras) y `interfaces`, que permiten crear tipos de datos personalizados y definir comportamientos comunes entre diferentes tipos.

---
## Control de flujo
Controlar el flujo de un programa es esencial para tomar decisiones y repetir acciones. Go proporciona varias estructuras de control de flujo, incluyendo condicionales y bucles.

**Condicionales:**
```go
if edad >= 18 {
    fmt.Println("Eres mayor de edad.")
} else {
    fmt.Println("Eres menor de edad.")
}
```

- `if edad >= 18 { ... }`: Evalúa si la variable `edad` es mayor o igual a 18. Si es verdadero, ejecuta el bloque de código dentro del `if`.
- `else { ... }`: Si la condición del `if` es falsa, ejecuta el bloque de código dentro del `else`.

**Sentencia switch:**
```go
switch dia {
case "Lunes":
    fmt.Println("Inicio de la semana.")
case "Viernes":
    fmt.Println("Fin de la semana.")
default:
    fmt.Println("Día normal.")
}
```

- `switch dia { ... }`: Evalúa la variable `dia` y ejecuta el bloque de código correspondiente al caso que coincida.
- `case "Lunes": ...`: Define un caso específico para el valor "Lunes".
- `default: ...`: Define un bloque de código que se ejecuta si ningún caso coincide.

**Bucles:**
En Go, el único tipo de bucle es el `for`, que puede usarse de varias maneras.
**Bucle for clásico:**
```go
for i := 0; i < 5; i++ {
    fmt.Println(i)
}
```
- `for i := 0; i < 5; i++ { ... }`: Inicia un bucle que se ejecuta mientras `i` sea menor que 5, incrementando `i` en 1 en cada iteración.

**Bucle for estilo while:**
```go
i := 0
for i < 5 {
    fmt.Println(i)
    i++
}
```
- `for i < 5 { ... }`: Inicia un bucle que se ejecuta mientras la condición `i < 5` sea verdadera.

**Bucle for infinito:**
```go
for {
    fmt.Println("Este bucle es infinito.")
    break // Rompe el bucle para evitar que sea realmente infinito
}
```
- `for { ... }`: Inicia un bucle infinito que se ejecuta indefinidamente hasta que se encuentre una instrucción `break`.

**Intrucción break y continue:**
```go
for i := 0; i < 10; i++ {
    if i%2 == 0 {
        continue // Salta a la siguiente iteración si i es par
    }
    if i == 7 {
        break // Rompe el bucle si i es 7
    }
    fmt.Println(i)
}
```
- `continue`: Salta a la siguiente iteración del bucle, omitiendo el código restante en la iteración actual.
- `break`: Rompe el bucle y sale de él inmediatamente.

**Uso de range:**
```go
frutas := []string{"🍎 Manzana", "🍌 Banana", "🍒 Cereza", "🍐 Pera"}
for index, fruta := range frutas {
    fmt.Printf("Fruta %d: %s\n", index, fruta)
}
```
- `for index, fruta := range frutas { ... }`: Itera sobre cada elemento del slice `frutas`, proporcionando el índice y el valor de cada elemento en cada iteración.

---
## Funciones
Las funciones en Go son bloques de código reutilizables que realizan tareas específicas. Se definen utilizando la palabra clave `func`, seguida del nombre de la función, los parámetros entre paréntesis y el tipo de retorno (si lo hay).

**Definición de una función simple:**
```go
func saludar(nombre string) {
    fmt.Printf("Hola, %s!\n", nombre)
}
```
- `func saludar(nombre string) { ... }`: Define una función llamada `saludar` que toma un parámetro `nombre` de tipo `string` y no devuelve ningún valor.

**Llamada a una función:**
```go
saludar("Juan")
```
- `saludar("Juan")`: Llama a la función `saludar`, pasando "

**Multiples valores de retorno:**
Ejemplo de una función que devuelve múltiples valores, una función que realiza una división y devuelve el cociente y el residuo:
```go
func dividir(a int, b int) (int, int) {
    cociente := a / b
    residuo := a % b
    return cociente, residuo
}
```
- `func dividir(a int, b int) (int, int) { ... }`: Define una función llamada `dividir` que toma dos parámetros `a` y `b` de tipo `int` y devuelve dos valores de tipo `int`.

Ahora, para llamar a esta función y capturar sus valores de retorno:
```go
cociente, residuo := dividir(10, 3)
fmt.Printf("Cociente: %d, Residuo: %d\n", cociente, residuo)
```
- `cociente, residuo := dividir(10, 3)`: Llama a la función `dividir` con los argumentos 10 y 3, y asigna los valores de retorno a las variables `cociente` y `residuo`. 

**Funciones variádicas:**
Una función variádica puede aceptar un número variable de argumentos del mismo tipo. Aquí tienes un ejemplo de una función que suma una cantidad indefinida de enteros:

```go
func sumar(numeros ...int) int {
    total := 0
    for _, numero := range numeros {
        total += numero
    }
    return total
}
```
- `func sumar(numeros ...int) int { ... }`: Define una función llamada `sumar` que acepta un número variable de argumentos de tipo `int` y devuelve un valor de tipo `int`.

**Llamada a una función variádica:**
```go
resultado := sumar(1, 2, 3, 4, 5)
fmt.Printf("La suma es: %d\n", resultado)
```
- `resultado := sumar(1, 2, 3, 4, 5)`: Llama a la función `sumar` con cinco argumentos y asigna el resultado a la variable `resultado`.

**Funciones anónimas**
Go permite definir funciones anónimas (sin nombre) y asignarlas a variables. Aquí tienes un ejemplo:
```go
mensaje := func(nombre string) string {
    return fmt.Sprintf("Hola, %s!", nombre)
}
fmt.Println(mensaje("Ana"))
```
- `mensaje := func(nombre string) string { ... }`: Define una función anónima que toma un parámetro `nombre` y devuelve un saludo formateado. Esta función se asigna a la variable `mensaje`.
- `fmt.Println(mensaje("Ana"))`: Llama a la función almacenada en `mensaje`, pasando "Ana" como argumento, e imprime el resultado.

**Funciones como argumentos**
En Go, las funciones pueden ser pasadas como argumentos a otras funciones. Aquí tienes un ejemplo:
```go
func ejecutarOperacion(a int, b int, operacion func(int, int) int) int {
    return operacion(a, b)
}
```
- `func ejecutarOperacion(a int, b int, operacion func(int, int) int) int { ... }`: Define una función llamada `ejecutarOperacion` que toma dos enteros y una función como parámetros, y devuelve un entero.
Ahora, puedes definir una función de suma y pasarla como argumento:
```go
suma := func(x int, y int) int {
    return x + y
}
resultado := ejecutarOperacion(5, 3, suma)
fmt.Printf("El resultado de la suma es: %d\n", resultado)
```
- `suma := func(x int, y int) int { ... }`: Define una función anónima para sumar dos enteros y la asigna a la variable `suma`. 

---
## Structs y métodos
Los structs en Go son tipos de datos compuestos que permiten agrupar múltiples campos bajo un mismo nombre. Son similares a las clases en otros lenguajes de programación, pero no tienen métodos asociados directamente. Sin embargo, puedes definir métodos para los structs. Aquí tienes un ejemplo de cómo definir un struct y agregarle métodos:

```go
type Persona struct {
    Nombre string
    Edad   int
}
```
- `type Persona struct { ... }`: Define un nuevo tipo de dato llamado `Persona` que tiene dos campos: `Nombre` de tipo `string` y `Edad` de tipo `int`.
- Para usar este struct, puedes crear una instancia de `Persona` y acceder a sus campos:
```go
persona := Persona{Nombre: "Juan", Edad: 30}
fmt.Printf("Nombre: %s, Edad: %d\n", persona.Nombre, persona.Edad)
```
- `persona := Persona{Nombre: "Juan", Edad: 30}`: Crea una instancia del struct `Persona` con el nombre "Juan" y la edad 30.
- `fmt.Printf("Nombre: %s, Edad: %d\n", persona.Nombre, persona.Edad)`: Imprime los valores de los campos `Nombre` y `Edad` de la instancia `persona`.

- `func (p Persona) Saludar() string { ... }`: Define un método llamado `Saludar` para el struct `Persona`. El receptor `p` es una instancia de `Persona`. Este método devuelve un saludo formateado como una cadena de texto.
```go
func (p Persona) Saludar() string {
    return fmt.Sprintf("Hola, mi nombre es %s y tengo %d años.", p.Nombre, p.Edad)
}
```

- `func (p *Persona) CumplirAnios() { ... }`: Define un método llamado `CumplirAnios` que incrementa la edad de la persona en 1. El receptor `p` es un puntero a una instancia de `Persona`, lo que permite modificar el valor original.
```go
func (p *Persona) CumplirAnios() {
    p.Edad++
}
```

- `persona := Persona{Nombre: "Juan", Edad: 30}`: Crea una instancia del struct `Persona` con el nombre "Juan" y la edad 30.

```go
func main() {
    persona := Persona{Nombre: "Juan", Edad: 30}
    fmt.Println(persona.Saludar()) // Llama al método Saludar

    persona.CumplirAnios() // Llama al método CumplirAnios
    fmt.Println(persona.Saludar()) // Verifica la edad actualizada
}
```


---
## Punteros
Los punteros en Go son temas importantes para entender cómo manejar la memoria y las referencias a variables. Un puntero es una variable que almacena la dirección de memoria de otra variable. Aquí tienes una explicación básica sobre punteros en Go:

```go
package main
import "fmt"

func main() {
    var x int = 42          // Declaración de una variable entera
    var p *int = &x        // Declaración de un puntero que apunta a la dirección de x

    fmt.Println("Valor de x:", x)          // Imprime el valor de x
    fmt.Println("Dirección de x:", &x)     // Imprime la dirección de memoria de x
    fmt.Println("Valor del puntero p:", p) // Imprime la dirección almacenada en p
    fmt.Println("Valor apuntado por p:", *p) // Desreferencia el puntero para obtener el valor de x

    *p = 100               // Modifica el valor de x a través del puntero
    fmt.Println("Nuevo valor de x:", x)    // Imprime el nuevo valor de x
}
```
- `var p *int = &x`: Declara un puntero `p` que apunta a la dirección de memoria de la variable `x`. El operador `&` se utiliza para obtener la dirección de una variable.
- `*p`: El operador de desreferenciación `*` se utiliza para acceder al valor almacenado en la dirección a la que apunta el puntero `p`.
- `*p = 100`: Modifica el valor de la variable `x a través del puntero `p`.

**Punteros en funciones:**
Los punteros son especialmente útiles cuando se pasan variables a funciones, ya que permiten modificar el valor original sin necesidad de devolverlo.
```go
func incrementar(valor *int) {
    *valor++ // Incrementa el valor al que apunta el puntero
}
```
- `func incrementar(valor *int) { ... }`: Define una función llamada `incrementar` que toma un puntero a un entero como parámetro.
```go
func main() {
    numero := 10
    fmt.Println("Antes de incrementar:", numero)
    incrementar(&numero) // Pasa la dirección de numero
    fmt.Println("Después de incrementar:", numero)
}
```
- `incrementar(&numero)`: Llama a la función `incrementar`, pasando la dirección de la variable `numero` utilizando el operador `&`.

**Punteros en structs:**
Los punteros también se pueden utilizar con structs para modificar sus campos directamente.
```go
type Punto struct {
    X int
    Y int
}

func mover(p *Punto, dx int, dy int) {
    p.X += dx
    p.Y += dy
}
```
- `func mover(p *Punto, dx int, dy int) { ... }`: Define una función llamada `mover` que toma un puntero a un struct `Punto` y dos enteros para desplazar las coordenadas.

```go
func main() {
    punto := Punto{X: 0, Y: 0}
    fmt.Printf("Antes de mover: %+v\n", punto)
    mover(&punto, 5, 10) // Pasa la dirección del struct punto
    fmt.Printf("Después de mover: %+v\n", punto)
}
```
- `mover(&punto, 5, 10)`: Llama a la función `mover`, pasando la dirección del struct `punto` para modificar sus campos directamente.

---
## Interfaces
Interfaces en Go son un tipo de dato que define un conjunto de métodos que un tipo debe implementar para satisfacer esa interfaz. Las interfaces permiten la abstracción y el polimorfismo, facilitando la escritura de código flexible y reutilizable. Aquí tienes una explicación básica sobre interfaces en Go:

```go
package main

import "fmt"

type Animal interface {
    HacerSonido() string
}

type Perro struct {
    Nombre string
}

func (p Perro) HacerSonido() string {
    return "Guau"
}

type Gato struct {
    Nombre string
}

func (g Gato) HacerSonido() string {
    return "Miau"
}

func main() {
    var animal Animal

    animal = Perro{Nombre: "Firulais"}
    fmt.Printf("%s dice: %s\n", animal.(Perro).Nombre, animal.HacerSonido())

    animal = Gato{Nombre: "Misu"}
    fmt.Printf("%s dice: %s\n", animal.(Gato).Nombre, animal.HacerSonido())
}
```
- `type Animal interface { ... }`: Define una interfaz llamada `Animal` que requiere que cualquier tipo que la implemente tenga un método `HacerSonido` que devuelva una cadena de texto.
- `type Perro struct { ... }`: Define un struct llamado `Perro` con un campo `Nombre`.
- `func (p Perro) HacerSonido() string { ... }`: Implementa el método `HacerSonido` para el struct `Perro`, devolviendo el sonido "Guau".
- `type Gato struct { ... }`: Define un struct llamado `Gato` con un campo `Nombre`.
- `func (g Gato) HacerSonido() string { ... }`: Implementa el método `HacerSonido` para el struct `Gato`, devolviendo el sonido "Miau".
- `var animal Animal`: Declara una variable `animal` de tipo `Animal`, que puede contener cualquier tipo que implemente la interfaz.
- `animal = Perro{Nombre: "Firulais"}`: Asigna una instancia de `Perro` a la variable `animal`.
- `animal.(Perro).Nombre`: Realiza una aserción de tipo para acceder al campo `Nombre` del struct `Perro`.
- `animal.HacerSonido()`: Llama al método `HacerSonido` de la interfaz `Animal`, que ejecuta la implementación correspondiente según el tipo concreto almacenado en `animal`.

---
## Manejo de errores
En Go, el manejo de errores es una parte fundamental del lenguaje y se realiza principalmente mediante el uso de valores de retorno. A diferencia de otros lenguajes que utilizan excepciones, Go prefiere un enfoque explícito para manejar errores. Aquí tienes una explicación básica sobre cómo manejar errores en Go:

**Filosofía del manejo de errores en Go:**
Go sigue la filosofía de "errores como valores", lo que significa que las funciones que pueden fallar devuelven un valor de error junto con el resultado esperado. Esto obliga a los desarrolladores a manejar los errores de manera explícita.

```go
package main
import (
    "errors"
    "fmt"
)

func dividir(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("no se puede dividir por cero")
    }
    return a / b, nil
}
func main() {
    resultado, err := dividir(10, 0)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Resultado:", resultado)
}
```
- `func dividir(a, b float64) (float64, error) { ... }`: Define una función llamada `dividir` que toma dos parámetros `a` y `b` de tipo `float64` y devuelve un `float64` y un `error`.
- `if b == 0 { ... }`: Verifica si el divisor `b` es cero. Si es así, devuelve un error utilizando `errors.New`.
- `return a / b, nil`: Si la división es válida, devuelve el resultado y `nil` para indicar que no hubo error.
- `resultado, err := dividir(10, 0)`: Llama a la función `dividir` y captura el resultado y el error.
- `if err != nil { ... }`: Verifica si hubo un error. Si es así, imprime el error y termina la ejecución.
- `fmt.Println("Resultado:", resultado)`: Si no hubo error, imprime el resultado de la división.    
---
## Paquetes y Módulos
Go utiliza paquetes y módulos para organizar y gestionar el código. Un paquete es una colección de archivos Go que se agrupan juntos, mientras que un módulo es una colección de paquetes versionados. Aquí tienes una explicación básica sobre paquetes y módulos en Go:

**Paquetes:**
Un paquete en Go es una forma de organizar el código en unidades reutilizables. Cada archivo Go comienza con una declaración de paquete que indica a qué paquete pertenece el archivo.
```go
package main

import "fmt"

func main() {
    fmt.Println("¡Hola, Mundo!")
}
```
- `package main`: Declara que este archivo pertenece al paquete `main`, que es el paquete especial que define un programa ejecutable en Go.
- `import "fmt"`: Importa el paquete `fmt`, que proporciona funciones para formatear y imprimir texto.

**Módulos:**
Un módulo en Go es una colección de paquetes que se versionan juntos. Los módulos se gestionan mediante el sistema de módulos de Go, que utiliza archivos `go.mod` para definir las dependencias del módulo.
Para crear un nuevo módulo, utiliza el comando `go mod init`:
```bash
go mod init mi-modulo
```
- `go mod init mi-modulo`: Inicializa un nuevo módulo llamado `mi-modulo, creando un archivo `go.mod` en el directorio actual.
El archivo `go.mod` contiene información sobre el módulo, incluyendo su nombre y las dependencias necesarias para el proyecto.
```go
module mi-modulo
go 1.16
```
- `module mi-modulo`: Define el nombre del módulo.
- `go 1.16`: Especifica la versión mínima de Go requerida para este módulo.

---
## JSON y Encoding
Go proporciona soporte integrado para trabajar con JSON a través del paquete `encoding/json`. Este paquete permite codificar (marshal) y decodificar (unmarshal) datos JSON de manera sencilla. Aquí tienes una explicación básica sobre cómo trabajar con JSON en Go:

```go
package main
import (
    "encoding/json"
    "fmt"
)

type Persona struct {
    Nombre string `json:"nombre"`
    Edad   int    `json:"edad"`
}

func main() {
    // Crear una instancia de Persona
    persona := Persona{Nombre: "Juan", Edad: 30}

    // Codificar (marshal) la estructura Persona a JSON
    jsonData, err := json.Marshal(persona)
    if err != nil {
        fmt.Println("Error al codificar a JSON:", err)
        return
    }
    fmt.Println("JSON codificado:", string(jsonData))

    // Decodificar (unmarshal) JSON a una estructura Persona
    var persona2 Persona
    err = json.Unmarshal(jsonData, &persona2)
    if err != nil {
        fmt.Println("Error al decodificar JSON:", err)
        return
    }
    fmt.Printf("Persona decodificada: %+v\n", persona2)
}
```
- `import "encoding/json"`: Importa el paquete `encoding/json`, que proporciona funciones para trabajar con JSON.
- `type Persona struct { ... }`: Define una estructura `Persona` con campos `Nombre` y `Edad`. Las etiquetas JSON (`json:"nombre"`) especifican cómo se deben nombrar los campos en el JSON.
- `json.Marshal(persona)`: Codifica la instancia `persona` a formato JSON. Devuelve los datos JSON y un error si ocurre.
- `json.Unmarshal(jsonData, &persona2)`: Decodifica los datos JSON en la estructura `persona2`. El segundo argumento es un puntero a la variable donde se almacenarán los datos decodificados.

**Trabajando con archivos JSON:**
Además de trabajar con JSON en memoria, también puedes leer y escribir archivos JSON utilizando el paquete `os` junto con `encoding/json`.
```go
import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
)  

func guardarJSON(nombreArchivo string, data interface{}) error {
    archivo, err := os.Create(nombreArchivo)
    if err != nil {
        return err
    }
    defer archivo.Close()

    encoder := json.NewEncoder(archivo)
    return encoder.Encode(data)
}

func cargarJSON(nombreArchivo string, data interface{}) error {
    archivo, err := os.Open(nombreArchivo)
    if err != nil {
        return err
    }
    defer archivo.Close()

    decoder := json.NewDecoder(archivo)
    return decoder.Decode(data)
}

type Persona struct {
    Nombre string `json:"nombre"`
    Edad   int    `json:"edad"`
}

func main() {
    persona := Persona{Nombre: "Roel", Edad: 31}
    err := guardarJSON("persona.json", persona)
    if err != nil {
        fmt.Println("Error al guardar JSON:", err)
        return
    }

    var personaCargada Persona
    err = cargarJSON("persona.json", &personaCargada)
    if err != nil {
        fmt.Println("Error al cargar JSON:", err)
        return
    }
    fmt.Printf("Persona cargada desde archivo: %+v\n", personaCargada)
}
```
- `os.Create(nombreArchivo)`: Crea un nuevo archivo con el nombre especificado para escribir datos JSON.
- `json.NewEncoder(archivo)`: Crea un codificador JSON que escribe en el archivo.
- `os.Open(nombreArchivo)`: Abre un archivo existente para leer datos JSON.
- `json.NewDecoder(archivo)`: Crea un decodificador JSON que lee desde el archivo.

---
## Concurrencia
Go tiene un modelo de concurrencia incorporado que facilita la creación de programas concurrentes utilizando goroutines y canales. Las goroutines son funciones que se ejecutan de manera concurrente, mientras que los canales permiten la comunicación segura entre goroutines. Aquí tienes una explicación básica sobre cómo usar la concurrencia en Go:

**Sin concurrencia:**
```go
package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	start := time.Now()

	apis := []string{
		"https://management.azure.com",
		"https://dev.azure.com",
		"https://api.github.com",
		"https://outlook.office.com/",
		"https://api.somewhereintheinternet.com/",
		"https://graph.microsoft.com",
	}

	// Recorreer los apis
	for _, api := range apis {
		checkAPI(api)
	}

	elapsed := time.Since(start)
	fmt.Printf("¡Listo! ¡Tomó %v segundos!\n", elapsed.Seconds())
}

// Función que verifica los APIS
func checkAPI(api string) {
	_, err := http.Get(api)
	if err != nil {
		fmt.Printf("ERROR: ¡%s está caído!\n", api)
		return
	}

	fmt.Printf("SUCCESS: ¡%s está en funcionamiento!\n", api)
}
```
- `for _, api := range apis { ... }`: Recorre cada API en la lista y llama a la función `checkAPI` de manera secuencial.

**Agregando concurrencia**
Para crear una goroutine, es necesario usar la palabra clave go antes de llamar a una función.
```go
	// Recorreer los apis
	for _, api := range apis {
		go checkAPI(api)
	}
```

Vuelva a ejecutar el programa y observe lo que sucede. Parece que el programa ya no comprueba las API, ¿verdad? Es posible que vea algo parecido a la salida siguiente:

```bash
¡Listo! ¡Tomó 2.7371e-05 segundos!
```
¡Muy rápido! ¿Qué ha ocurrido? Verá el mensaje final que indica que el programa ha finalizado porque Go ha creado una goroutine para cada sitio dentro del bucle e inmediatamente a pasado a la siguiente línea.

Aunque no parece que la función checkAPI se esté ejecutando, realmente sí lo está haciendo. Simplemente no tuvo tiempo de finalizarse. Observe lo que ocurre si incluye un temporizador de suspensión justo después del bucle:

```go
	// Recorreer los apis
	for _, api := range apis {
		go checkAPI(api)
	}

	time.Sleep(5 * time.Second)

```
Ahora, cuando vuelva a ejecutar el programa, podría ver una salida similar a la siguiente:
```bash
ERROR: ¡https://api.somewhereintheinternet.com/ está caído!
SUCCESS: ¡https://api.github.com está en funcionamiento!
SUCCESS: ¡https://dev.azure.com está en funcionamiento!
SUCCESS: ¡https://management.azure.com está en funcionamiento!
SUCCESS: ¡https://outlook.office.com/ está en funcionamiento!
SUCCESS: ¡https://graph.microsoft.com está en funcionamiento!
¡Listo! ¡Tomó 5.002491318 segundos!
```
Parece que funciona, ¿verdad? En realidad, no exactamente. ¿Qué ocurre si desea agregar un nuevo sitio a la lista? Quizás tres segundos no son suficientes. ¿Cómo podría saberlo? No puede. Debe haber una manera mejor, y eso es lo que analizaremos en la sección siguiente cuando hablemos de los canales.

**Uso de canales:**
En Go, los canales son una característica fundamental para la comunicación y sincronización entre goroutines (subprocesos ligeros) dentro de un programa concurrente. Un canal es una estructura que permite enviar y recibir valores entre goroutines, actuando como un conducto a través del cual fluye la información.

Para crear un canal en Go, se utiliza la función make() con la siguiente sintaxis:
```go
canal := make(chan tipoDato)
```

Donde tipoDato especifica el tipo de datos que se enviarán a través del canal. Puede ser cualquier tipo de datos válido en Go, como int, string, struct, etc.

Una vez creado el canal, se pueden enviar y recibir datos utilizando la notación de flecha <-. Por ejemplo:

```go
// Crear un canal de tipo entero
canal := make(chan int)

// Enviar un valor a través del canal
canal <- 10

// Recibir un valor del canal
valor := <-canal
```
La operación <- se utiliza para enviar un valor al canal (colocándolo a la izquierda de la flecha) o recibir un valor del canal (colocándolo a la derecha de la flecha).

**Canales y concurrencia:**
En el programa use canales para quitar la funcionalidad de suspensión. En primer lugar, vamos a crear un canal de cadena en la función main, como se indica a continuación:
```go
// Crear un canal de tipo string
ch := make(chan string)
```

Y quitaremos la línea de suspensión `time.Sleep(5 * time.Second)`.

Ahora, podemos usar canales para comunicarse entre goroutines. En lugar de imprimir el resultado en la función checkAPI, se refactorizará el código y ese mensaje se enviará por el canal. Para usar el canal desde esa función, debe agregar el canal como parámetro. La función checkAPI debe tener el siguiente aspecto:

```go
// Función que verifica los APIS
func checkAPI(api string, ch chan string) {
	_, err := http.Get(api)
	if err != nil {
		ch <- fmt.Sprintf("ERROR: ¡%s está caído!\n", api)
		return
	}

	ch <- fmt.Sprintf("SUCCESS: ¡%s está en funcionamiento!\n", api)
}
```
Tenga en cuenta que es necesario usar la función fmt.Sprintf porque no quiere imprimir ningún texto, simplemente enviar texto con formato por el canal. Además, observe que usamos el operador <- después de la variable de canal para enviar datos.

Ahora debe cambiar la función main para enviar la variable de canal y recibir los datos para imprimirla, como se muestra a continuación:

```go
	// Recorreer los apis
	for _, api := range apis {
		go checkAPI(api, ch)
	}

	// Leer datos de canal
	fmt.Println(<-ch)

```
Observe cómo usamos el operador <- antes de que el canal indique que queremos leer datos del canal. Cuando vuelva a ejecutar el programa, verá una salida similar a la siguiente:
```bash
ERROR: ¡https://api.somewhereintheinternet.com/ está caído!

¡Listo! ¡Tomó 0.009662042 segundos!ç
```

Al menos funciona sin una llamada a una función de suspensión, ¿no? Pero todavía no hace lo que queremos. Vemos la salida solo de una de las goroutines, pero creamos cinco. En la siguiente clase descubriremos por qué este programa funciona de esta manera.

**Canales no almacenados en búfer:**
Cuando se crea un canal mediante la función make(), se crea un canal no almacenado en búfer, que es el comportamiento predeterminado. Los canales no almacenados en búfer bloquean la operación de envío hasta que algún componente esté listo para recibir los datos. Como se ha afirmado antes, el envío y la recepción son operaciones de bloqueo. Esta operación de bloqueo también es la razón por la que el programa de la sección anterior se ha detenido en cuanto ha recibido el primer mensaje.

Podemos empezar diciendo que fmt.Print(<-ch) bloquea el programa porque está leyendo de un canal y espera a que lleguen algunos datos. En cuanto hay algunos, continúa con la línea siguiente y el programa finaliza.

¿Qué ha ocurrido con el resto de las goroutines? Todavía se están ejecutando, pero ya no hay ninguna escuchando. Y dado que el programa terminó pronto, algunas goroutines no pudieron enviar datos. Para demostrar esto, vamos a agregar otra línea fmt.Print(<-ch), como se indica a continuación:
```go
	ch := make(chan string)

	// Recorreer los apis
	for _, api := range apis {
		go checkAPI(api, ch)
	}

	fmt.Print(<-ch)
	fmt.Print(<-ch)
```

Cuando vuelva a ejecutar el programa, verá una salida similar a la siguiente:
```bash
ERROR: ¡https://api.somewhereintheinternet.com/ está caído!
SUCCESS: ¡https://api.github.com está en funcionamiento!
¡Listo! ¡Tomó 0.48367305 segundos!
```

Observe que ahora verá la salida de dos API. Si continúa agregando más líneas fmt.Print(<-ch), acabará leyendo todos los datos que se envían al canal. Pero ¿qué ocurre si intenta leer más datos y ya no hay ninguna goroutine que envíe datos? Por ejemplo:
```go
ch := make(chan string)

for _, api := range apis {
    go checkAPI(api, ch)
}

fmt.Print(<-ch)
fmt.Print(<-ch)
fmt.Print(<-ch)
fmt.Print(<-ch)
fmt.Print(<-ch)
fmt.Print(<-ch)

fmt.Print(<-ch)
```
Cuando vuelva a ejecutar el programa, verá una salida similar a la siguiente:
```bash
ERROR: ¡https://api.somewhereintheinternet.com/ está caído!
SUCCESS: ¡https://api.github.com está en funcionamiento!
SUCCESS: ¡https://management.azure.com está en funcionamiento!
SUCCESS: ¡https://dev.azure.com está en funcionamiento!
SUCCESS: ¡https://graph.microsoft.com está en funcionamiento!
SUCCESS: ¡https://outlook.office.com/ está en funcionamiento!
```
Funciona, pero el programa no finaliza. La última línea de impresión lo está bloqueando porque está esperando recibir datos. Tendrá que cerrar el programa con un comando como Ctrl+C.

El ejemplo anterior simplemente demuestra que la lectura y recepción de datos son operaciones de bloqueo. Para corregir este problema, podría cambiar el código a un bucle for y recibir solo los datos que sabe con certeza que va a enviar, como en este ejemplo:

```go
for i := 0; i < len(apis); i++ {
    fmt.Print(<-ch)
}
```
El programa está haciendo lo que se supone que debe hacer. Ya no usa una función de suspensión; usa canales. Observe también que ahora se tardan aproximadamente 1.357984 segundos en finalizar en lugar de los casi 5 segundos cuando no se usaba la simultaneidad.

---
## Context (context.Context)
El paquete `context` en Go proporciona una forma de manejar la cancelación, los plazos y los valores asociados con las solicitudes y operaciones concurrentes. Es especialmente útil en aplicaciones web y servicios donde las operaciones pueden necesitar ser canceladas o tener un tiempo límite. Aquí tienes una explicación básica sobre cómo usar `context.Context` en Go:

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    // Crear un contexto con un plazo de 2 segundos
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel() // Asegura que se liberen los recursos del contexto

    // Simular una operación que toma tiempo
    resultadoCh := make(chan string)
    go func() {
        time.Sleep(3 * time.Second) // Simula una operación larga
        resultadoCh <- "Operación completada"
    }()

    select {
    case resultado := <-resultadoCh:
        fmt.Println(resultado)
    case <-ctx.Done():
        fmt.Println("Operación cancelada o tiempo agotado:", ctx.Err())
    }
}
```

- `context.WithTimeout(context.Background(), 2*time.Second)`: Crea un nuevo contexto que se cancelará automáticamente después de 2 segundos. `context.Background()` es el contexto raíz.
- `defer cancel()`: Asegura que se liberen los recursos asociados con el contexto cuando la función `main` termine.
- `go func() { ... }()`: Inicia una goroutine que simula una operación larga.
- `select { ... }`: Permite esperar en múltiples canales. En este caso, espera a que la operación se complete o a que el contexto se cancele.
- `case <-ctx.Done()`: Se activa cuando el contexto se cancela, ya sea por el plazo agotado o por una cancelación explícita.

---
## Testing en Go
El paquete `testing` en Go proporciona un marco integrado para escribir y ejecutar pruebas unitarias. Las pruebas en Go se escriben en archivos separados con el sufijo `_test.go` y utilizan funciones que comienzan con `Test`. Aquí tienes una explicación básica sobre cómo escribir y ejecutar pruebas en Go:
```go
package main
import "testing"
func Sumar(a, b int) int {
    return a + b
}
func TestSumar(t *testing.T) {
    resultado := Sumar(2, 3)
    esperado := 5
    if resultado != esperado {
        t.Errorf("Sumar(2, 3) = %d; se esperaba %d", resultado, esperado)
    }
}
```
- `package main`: Declara el paquete principal.
- `import "testing"`: Importa el paquete `testing`, que proporciona las herramientas necesarias para escribir pruebas.
- `func Sumar(a, b int) int { ... }`: Define una función simple que suma dos enteros.
- `func TestSumar(t *testing.T) { ... }`: Define una función de prueba llamada `TestSumar`. La función debe comenzar con `Test` y aceptar un parámetro de tipo `*testing.T`.
- `t.Errorf(...)`: Registra un error en la prueba si el resultado no coincide con el valor esperado.   

---
## HTTP y net/http
**¿Qué es HTTP?**
HTTP (Hypertext Transfer Protocol) es el protocolo de comunicación utilizado en la World Wide Web para la transferencia de datos entre clientes (como navegadores web) y servidores web. Es un protocolo basado en texto que define cómo se formatean y transmiten los mensajes, así como las acciones que deben tomarse en respuesta a diversas solicitudes.

**Métodos HTTP comunes:**
- `GET`: Solicita la representación de un recurso específico. Las solicitudes GET solo deben recuperar datos y no deben tener efectos secundarios.
- `POST`: Envía datos al servidor para crear o actualizar un recurso. Las solicitudes POST pueden tener efectos secundarios en el servidor.
- `PUT`: Reemplaza todas las representaciones actuales del recurso de destino con los datos de la solicitud.
- `DELETE`: Elimina el recurso especificado.
- `PATCH`: Aplica modificaciones parciales a un recurso.

**Códigos de estado HTTP comunes:**
- `200 OK`: La solicitud se ha procesado correctamente.
- `201 Created`: La solicitud se ha completado y se ha creado un nuevo recurso.
- `400 Bad Request`: La solicitud no se pudo entender o fue malformada.
- `401 Unauthorized`: La solicitud requiere autenticación del usuario.
- `403 Forbidden`: El servidor entendió la solicitud, pero se niega a autorizarla.
- `404 Not Found`: El recurso solicitado no se encontró en el servidor.
- `500 Internal Server Error`: El servidor encontró una condición inesperada que le impidió cumplir con la solicitud.

**Tipo de contenido común:**
- `application/json`: Indica que el cuerpo del mensaje contiene datos en formato JSON.
- `text/html`: Indica que el cuerpo del mensaje contiene datos en formato HTML.

**Paquete net/http en Go:**
El paquete `net/http` en Go proporciona funcionalidades para construir clientes y servidores HTTP. Aquí tienes una explicación básica sobre cómo usar `net/http` para crear un servidor web simple:

```go
package main
import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "¡Hola, Mundo!")
}

func main() {
    http.HandleFunc("/", handler) // Asocia la ruta raíz con el manejador
    fmt.Println("Servidor escuchando en http://localhost:8080")
    http.ListenAndServe(":8080", nil) // Inicia el servidor en el puerto 8080
}
```

- `import "net/http"`: Importa el paquete `net/http`, que proporciona funcionalidades para trabajar con HTTP.
- `func handler(w http.ResponseWriter, r *http.Request) { ... }`: Define una función manejadora que responde a las solicitudes HTTP. El parámetro `w` se utiliza para escribir la respuesta, y `r` contiene la solicitud entrante.
- `http.HandleFunc("/", handler)`: Asocia la ruta raíz ("/") con la función manejadora `handler`.

---
## Mini API RESTful
**¿Qué es una API RESTful?**
Una API RESTful (Representational State Transfer) es un conjunto de convenciones y principios para diseñar servicios web que permiten la comunicación entre sistemas a través de HTTP. Las APIs RESTful utilizan los métodos HTTP estándar (GET, POST, PUT, DELETE) para realizar operaciones sobre recursos representados en formato JSON, XML u otros.

**Características clave de una API RESTful:**
- **Recursos**: Los recursos son las entidades que la API expone, como usuarios, productos o pedidos. Cada recurso se identifica mediante una URL única.
- **Métodos HTTP**: Los métodos HTTP se utilizan para realizar operaciones sobre los recursos. Por ejemplo, GET para recuperar datos, POST para crear nuevos recursos, PUT para actualizar recursos existentes y DELETE para eliminar recursos.
- **Stateless**: Las APIs RESTful son sin estado, lo que significa que cada solicitud del cliente al servidor debe contener toda la información necesaria para procesar la solicitud. El servidor no mantiene el estado entre solicitudes.
- **Representaciones**: Los recursos pueden tener múltiples representaciones, como JSON o XML. El cliente puede especificar el formato deseado mediante encabezados HTTP. 

**Ejemplo de una API RESTful simple en Go:**
Para esto utilizaremos el paquete `net/http` para crear una API RESTful básica que maneje operaciones CRUD (Crear, Leer, Actualizar, Eliminar) para un recurso llamado "Artículo" en memoria.

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "sync"
)

type Articulo struct {
    ID    int    `json:"id"`
    Titulo string `json:"titulo"`
    Contenido string `json:"contenido"`
}

var (
    articulos = make(map[int]Articulo)
    idCounter = 1
)

// Handler para crear un nuevo artículo
func crearArticulo(w http.ResponseWriter, r *http.Request) {
    var articulo Articulo
    err := json.NewDecoder(r.Body).Decode(&articulo)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    articulo.ID = idCounter
    articulos[idCounter] = articulo
    idCounter++


    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(articulo)
}

// Handler para obtener todos los artículos
func obtenerArticulos(w http.ResponseWriter, r *http.Request) {
    var lista []Articulo
    for _, articulo := range articulos {
        lista = append(lista, articulo)
    }
    json.NewEncoder(w).Encode(lista)
}
    
// Handler para obtener un artículo por ID
func obtenerArticuloPorID(w http.ResponseWriter, r *http.Request) {
    // articulo/{id}
    idStr := r.Params["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    articulo, exists := articulos[id]
    if !exists {
        http.Error(w, "Artículo no encontrado", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(articulo)
}

// Handler para actualizar un artículo por ID
func actualizarArticulo(w http.ResponseWriter, r *http.Request) {
    idStr := r.Params["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }
    var articuloActualizado Articulo
    err = json.NewDecoder(r.Body).Decode(&articuloActualizado)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    articuloActualizado.ID = id
    articulos[id] = articuloActualizado
    json.NewEncoder(w).Encode(articuloActualizado)
}

// Handler para eliminar un artículo por ID
func eliminarArticulo(w http.ResponseWriter, r *http.Request) {
    idStr := r.Params["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil { 
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }
    delete(articulos, id)
    w.WriteHeader(http.StatusNoContent)
}

// Configuración de rutas y servidor
func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("POST /articulos", crearArticulo)
    mux.HandleFunc("GET /articulos", obtenerArticulos)
    mux.HandleFunc("GET /articulos/{id}", obtenerArticuloPorID)
    mux.HandleFunc("PUT /articulos/{id}", actualizarArticulo)
    mux.HandleFunc("DELETE /articulos/{id}", eliminarArticulo)

    fmt.Println("Servidor escuchando en http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
```
- `type Articulo struct { ... }`: Define una estructura `Articulo` con campos `ID`, `Titulo` y `Contenido`.
- `var articulos = make(map[int]Articulo)`: Crea un mapa para almacenar los artículos en memoria.
- `func crearArticulo(w http.ResponseWriter, r *http.Request) { ... }`: Define un manejador para crear un nuevo artículo.
- `func obtenerArticulos(w http.ResponseWriter, r *http.Request) { ... }`: Define un manejador para obtener todos los artículos.
- `func obtenerArticuloPorID(w http.ResponseWriter, r *http.Request) { ... }`: Define un manejador para obtener un artículo por su ID.
- `func actualizarArticulo(w http.ResponseWriter, r *http.Request) { ... }`: Define un manejador para actualizar un artículo por su ID.
- `func eliminarArticulo(w http.ResponseWriter, r *http.Request) { ... }`: Define un manejador para eliminar un artículo por su ID.
- `http.NewServeMux()`: Crea un nuevo multiplexor de solicitudes HTTP para manejar las rutas.
- `mux.HandleFunc(...)`: Asocia las rutas con sus respectivos manejadores.
- `http.ListenAndServe(":8080", mux)`: Inicia el servidor HTTP en el puerto 8080.

---
## Conclusión
En esta guía, hemos explorado varios conceptos fundamentales de Go, incluyendo punteros, interfaces, manejo de errores, paquetes y módulos, JSON y encoding, concurrencia, context.Context, testing en Go, HTTP y net/http, y la creación de una mini API RESTful. Estos conceptos son esenciales para desarrollar aplicaciones robustas y eficientes en Go. A medida que continúes aprendiendo y practicando, te familiarizarás más con las características avanzadas del lenguaje y podrás construir aplicaciones más complejas y escalables. ¡Feliz codificación!


