package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Uso: go run main.go <NATS_SERVER> <CHAT_CHANNEL> <USER_NAME>")
		os.Exit(1)
	}

	natsServer := os.Args[1]
	chatChannel := os.Args[2]
	userName := os.Args[3]

	// Conexión al servidor NATS
	nc, err := nats.Connect(natsServer)
	if err != nil {
		fmt.Println("Error conectando al servidor NATS:", err)
		os.Exit(1)
	}
	defer nc.Close()

	// Conexión a JetStream
	js, err := nc.JetStream()
	if err != nil {
		fmt.Println("Error conectando a JetStream:", err)
		os.Exit(1)
	}

	// Verificar o crear el Stream
	streamName := "CHAT"
	streamInfo, err := js.StreamInfo(streamName)
	if err != nil && err != nats.ErrStreamNotFound {
		fmt.Println("Error verificando el Stream:", err)
		return
	}

	// Si el Stream no existe, crearlo
	if streamInfo == nil {
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{chatChannel},
			Storage:  nats.FileStorage,
			MaxAge:   time.Hour, // Retener mensajes solo de la última hora
		})
		if err != nil {
			fmt.Println("Error creando el Stream:", err)
			return
		}
		fmt.Println("Stream creado exitosamente.")
	} else {
		fmt.Println("El Stream ya existe. Usándolo.")
	}

	fmt.Printf("Conectado al servidor NATS en %s\n", natsServer)
	fmt.Printf("Uniéndose al canal '%s' como '%s'\n", chatChannel, userName)

	// Recuperar mensajes históricos
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	consumerName := userName + "_consumer"
	_, err = js.AddConsumer(streamName, &nats.ConsumerConfig{
		Durable:       consumerName,
		DeliverPolicy: nats.DeliverAllPolicy,
		AckPolicy:     nats.AckExplicitPolicy,
	})
	if err != nil {
		fmt.Println("Error creando el consumidor:", err)
		return
	}

	sub, err := js.PullSubscribe(chatChannel, consumerName)
	if err != nil {
		fmt.Println("Error creando el suscriptor Pull:", err)
		return
	}

	// Intentar recuperar mensajes históricos
	msgs, err := sub.Fetch(10, nats.Context(ctx))
	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Println("No se encontraron mensajes históricos. Continuando...")
		} else {
			fmt.Println("Error obteniendo mensajes históricos:", err)
		}
	} else {
		for _, m := range msgs {
			// Filtrar mensajes con base en el timestamp
			timestamp, _ := time.Parse(time.RFC3339, m.Header.Get("timestamp"))
			if time.Since(timestamp) <= time.Hour {
				fmt.Printf("Histórico %s: %s\n", m.Header.Get("user"), string(m.Data))
				m.Ack()
			}
		}
	}

	// Suscribirse al canal solo para mensajes nuevos
	_, err = js.Subscribe(chatChannel, func(msg *nats.Msg) {
		fmt.Printf("%s: %s\n", msg.Header.Get("user"), string(msg.Data))
	}, nats.DeliverNew()) // Importante: solo mensajes nuevos
	if err != nil {
		fmt.Println("Error suscribiéndose al canal:", err)
		os.Exit(1)
	}

	// Publicar mensajes en el canal
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			message := scanner.Text()
			msg := &nats.Msg{
				Subject: chatChannel,
				Data:    []byte(message),
				Header: nats.Header{
					"user":      []string{userName},
					"timestamp": []string{time.Now().Format(time.RFC3339)},
				},
			}
			err := nc.PublishMsg(msg)
			if err != nil {
				fmt.Println("Error enviando mensaje:", err)
			}
		}
	}()

	// Mantener la conexión activa
	select {}
}





