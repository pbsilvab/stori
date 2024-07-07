# Prueba Tecnica Para Stori

Este proyecto es una solución para un desafío de fintech que procesa archivos CSV desde un directorio montado, calcula un resumen de las transacciones y envía la información resumida por correo electrónico.

## Estructura del Proyecto

root/
│
├── cmd/
│ └── main.go
│
├── internal/
│ ├── account/
│ ├── emailtemplate/
│ └── fileprocessor/
|
|__tmp/
|   |__transactions/
|   |__emails
│
├── Dockerfile
└── docker-compose.yml

Nota: El Directorio tmp/transactions simula un FS montado que se le pasa como volumen al container de docker.


## Descripción de los Paquetes

### `internal/account`

Este paquete maneja las operaciones relacionadas con las cuentas, incluyendo la adición de transacciones y el cálculo del balance y los resúmenes de las transacciones.

### `internal/emailtemplate`

Este paquete maneja la generación de plantillas de correo electrónico y la creación de archivos de correo electrónico con información resumida.

### `internal/fileprocessor`

Este paquete maneja el procesamiento de archivos CSV desde un directorio montado, incluyendo la lectura de archivos y la gestión de archivos de bloqueo para evitar la doble lectura.

### API HTTP

- `POST /process`: Procesa los archivos en el directorio especificado y calcula el resumen de las transacciones.
  - Parámetros (en el cuerpo de la solicitud):
    - `directory` (string): El directorio donde se encuentran los archivos CSV.
    - `accountId` (string): El ID de la cuenta para procesar las transacciones.
    - `name` (string): Nombre del cliente para propositos de personalizar el email.

## Ejecución del Código

### Requisitos
- Docker
- Docker Compose

- AWS ENVS:
    AWS_ACCOUNT_KEY
    AWS_ACCOUNT_SECRET
    AWS_REGION

- AWS DynamoDB: Configura una tabla de DynamoDB llamada: transactions

### Instrucciones

1. Clona este repositorio en tu máquina local.
2. Navega al directorio raíz del proyecto.
3. Asegúrate de que Docker esté instalado y funcionando en tu máquina.
4. Completa las variables de entorno en el archivo docker-compose. (example.docker-compose.yml)
5. Ejecuta el siguiente comando para construir y levantar los servicios:

```sh
docker-compose up --build