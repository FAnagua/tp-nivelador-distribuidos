# TP NIVELADOR

Nombre y Apellido: Facundo Anagua Rocabado

Padrón: 109641

## Protocolo

### Cliente

* Comando para enviar el numero de agencia del cliente:
    ```
    [comando:1byte -> 0x01]
    [agency_id:1bytes]
    ```

* Comando para el envio de cada batch:
    ```
    [comando:1byte -> 0x02]
    [cantidad_de_apuestas:2bytes]
    [largo_nombre:2bytes]
    [nombre]
    [largo_apellido:2bytes]
    [apellido]
    [documento:4bytes]
    [largo_dia_cumpleaños:2bytes]
    [dia_cumpleaños]
    [numero:4bytes]
    ```
    En el batch se repite desde  `largo_nombre` hasta  `numero` dependiendo la cantidad de apuestas dentro de la misma.

* Comando para pedir los resultados del sorteo:
    ```
    [comando:1byte -> 0x05]
    ```

* Comando de finalización de comunicación con el servidor:
    ```
    [comando:1byte -> 0x04]
    ```

### Servidor

* Comando para confirmar que recibios el batch con las apuestas:
    ```
    [comando:1byte -> 0x03]
    ```
* Comando para enviar los resultados del sorteo:
    ```
    [comando:1byte -> 0x06]
    [cantidad_de_ganadores:2bytes]
    [largo_nombre:2bytes]
    [nombre]
    [largo_apellido:2bytes]
    [apellido]
    [documento:4bytes]
    [largo_dia_cumpleaños:2bytes]
    [dia_cumpleaños]
    [numero:4bytes]
    ```
    Se repite desde  `largo_nombre` hasta  `numero` dependiendo la cantidad de ganadores por agencia.

## Mecanismo de Sincronización

Para sincronizar el acceso a escritura y lectura del archivo que almacena las apuestas usamos un monitor para protoger el recurso.

Como mecanismo de concurrencia usamos Threads, no nos afecta en gran medidad los limtes de Python por parte del GIL porque se realizan muchas operaciones O/I dando la aportunidad de poder liberarlo.

Luego para la lectura de los ganadores se uso una barrera permite realizar la lectura respecto a `AGENCY_QUORUM_MIN`.
