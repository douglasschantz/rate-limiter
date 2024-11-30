package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	redis "github.com/redis/go-redis/v9"
)

var limite int = 5

func main() {

	// Iniciar o servidor HTTP.
	http.HandleFunc("/", Handler)
	fmt.Println("Servidor escutando em http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}

func Handler(w http.ResponseWriter, r *http.Request) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer client.Close()

	identifier := r.RemoteAddr
	ip := fmt.Sprintf("ratelimit:%s", identifier)

	ctx := context.Background()

	// Incrementar a contagem de requisições por IP
	// Usaremos um hash do Redis para manter a contagem
	_, err := client.HIncrBy(ctx, "requests", ip, 1).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Obter a contagem atual de requisições para o IP
	count, err := client.HGet(ctx, "requests", ip).Int()

	if err == redis.Nil {
		client.Set(ctx, ip, 1, 1*time.Second)
	} else if err != nil {
		http.Error(w, "Redis Error", http.StatusInternalServerError)
	} else if count >= limite {
		// Limite Excedido
		http.Error(w, "Too many requests - Limite Excedido", http.StatusTooManyRequests)
	} else {
		// Incrementando ip
		client.Incr(ctx, ip)
	}

	// Processar a requisição normalmente
	fmt.Fprintf(w, "Requisição processada com sucesso para %s", ip)

}
