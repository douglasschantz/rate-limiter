package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	redis "github.com/redis/go-redis/v9"
)

var limite int = 2

func main() {

	// Iniciar o servidor HTTP.
	http.HandleFunc("/", Handler)
	fmt.Println("Servidor escutando em http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}

func Handler(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer client.Close()
	ctx := context.Background()

	// Incrementar a contagem de requisições por IP
	// Usaremos um hash do Redis para manter a contagem
	_, err := client.HIncrBy(ctx, "requests", ip, 1).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Definir um tempo de expiração para o hash (1 segundo)
	_, err = client.Expire(ctx, "requests", time.Duration(limite)*time.Second).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Obter a contagem atual de requisições para o IP
	count, err := client.HGet(ctx, "requests", ip).Int()

	if err != nil && err != redis.Nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verificar se a contagem excedeu o limite de 10 requisições por segundo
	if count > 10 {
		http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
		return
	}

	// Processar a requisição normalmente
	fmt.Fprintf(w, "Requisição processada com sucesso para %s", ip)

}
