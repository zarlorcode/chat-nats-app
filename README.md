# Chat CLI con NATS

## Descripción
Esta es una aplicación de chat colaborativo en línea basada en Go y NATS. Utiliza el modelo PUB/SUB para mensajes en tiempo real y JetStream para persistencia.

## Requisitos
- Go 1.20+
- Docker

## Configuración
1. Clona el repositorio:
   ```bash
   git clone https://github.com/zarlorcode/chat-nats-app.git
   cd chat-nats-go
   
2. Instalar dependencias de Go
    ```bash
    go get github.com/nats-io/nats.go

3. Levantar el servidor NATS
    ```bash
    docker-compose up -d
    
4. Ejecutar la aplicación
    ```bash
    go run main.go <NATS_SERVER> <CHAT_CHANNEL> <USER_NAME>
    
Parámetros:
    <NATS_SERVER>: La dirección del servidor NATS (por ejemplo, nats://localhost:4222).
    <CHAT_CHANNEL>: El nombre del canal de chat (por ejemplo, chatroom).
    <USER_NAME>: El nombre o identificador que se usará en el chat.
    
5. Ejemplo de uso
    Vamos a lanzar 3 usuarios desde 3 terminales diferentes.
    
    En la primera terminal:
    ```bash
    go run main.go nats://localhost:4222 chatroom usuario1
    
    En la segunda terminal:
    ```bash
    go run main.go nats://localhost:4222 chatroom usuario2
    
    Ahora escribimos un par de mensajes para comprobar que el chat funciona.
    A continuación añadiremos otro usuario para comprobar la persistencia del historico del chat.
    
    En la tercera terminal:
    ```bash
    go run main.go nats://localhost:4222 chatroom usuario3
    
    El usuario 3 deberá haber recibido los mensajes anteriores de los otros 2 usuarios.

## Componentes de la aplicación

1. main.go

    Descripción: Este es el archivo principal de la aplicación, escrito en Go, que implementa las funcionalidades principales del chat.
    
    Funciones clave:
        - Se conecta al servidor NATS proporcionado por el usuario.
        - Publica mensajes en el canal especificado, con encabezados para el usuario y la marca de tiempo.
        - Recupera mensajes históricos utilizando JetStream, filtrando aquellos enviados en la última hora.
        - Se suscribe al canal para recibir mensajes en tiempo real.
        - Maneja tanto la publicación como la recepción de mensajes a través de la línea de comandos.
        
2. docker-compose.yml

    Descripción: Archivo de configuración para Docker Compose que levanta un servidor NATS con soporte para JetStream.
    
    Funciones clave:
        - Configura el servidor NATS para ejecutarse localmente en el puerto 4222.
        - Habilita JetStream para el almacenamiento persistente de mensajes.

3. go.mod:

    Archivo de configuración de módulos de Go, generado automáticamente al iniciar el proyecto con go mod init. Contiene las dependencias necesarias, como github.com/nats-io/nats.go.

    


