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

### `internal/emailservice`

Capa de implementacion abstracta de diferentes proveedores de mailing services. (Solo implementado AWS SES)

### `internal/emailsender`

Package que resuleve el envio de los emails alamacenados en el storage configurado. implementado un servicio de almacenamiento en disco y otro en AWS SQS

### `internal/fileprocessor`

Este paquete maneja el procesamiento de archivos CSV desde un directorio montado, incluyendo la lectura de archivos y la gestión de archivos de bloqueo para evitar la doble lectura.

### `internal/sqsclient`

implementacion de un cliente de sqs para su uso en otros pkg


### API HTTP Accounts

- `POST /create-account`: Crea un registro en la base de datos 
  - Parámetros (en el cuerpo de la solicitud):
    - `name` (string): Nombre del poseedor de la cuenta
    - `email` (string): correo electronico del poseedo de la cuenta

- `GET /find/{account_id}`: Busca los datos de una cuenta. 

- `GET /list`: Busca todos los registros en la base de datos. 

### API HTTP Transacctions

- `POST /process`: Procesa los archivos en el directorio especificado y calcula el resumen de las transacciones.
  - Parámetros (en el cuerpo de la solicitud):
    - `directory` (string): El directorio donde se encuentran los archivos CSV.
    - `accountId` (string): El ID de la cuenta para procesar las transacciones.

## Ejecución del Código

### Requisitos
- Docker
- Docker Compose
- AWS SERVICES: 
    SES
    SQS
    DYNAMODB
    EFS
    
- AWS ENVS:
    AWS_ACCOUNT_KEY
    AWS_ACCOUNT_SECRET
    AWS_REGION
    SQS_REGION 
    SQS_URL
    AWS_SES_SENDER_EMAIL


- AWS DynamoDB: Configura una tabla de DynamoDB llamada *accountIndo*:
 ``` key: id ```
- AWS DynamoDB: Configura una tabla de DynamoDB llamada *transactions*: 
 ``` key: account ```

### Instrucciones

1. Clona este repositorio en tu máquina local.
2. Navega al directorio raíz del proyecto.
3. Asegúrate de que Docker esté instalado y funcionando en tu máquina.
4. Completa las variables de entorno en el archivo docker-compose. (example.docker-compose.yml)
5. Ejecuta el siguiente comando para construir y levantar los servicios:

```sh
docker-compose up --build