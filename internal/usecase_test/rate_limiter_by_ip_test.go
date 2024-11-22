// rate_limiter_by_ip_test.go

package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/douglasschantz/rate-limiter/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestCheckRateLimitByIP(t *testing.T) {
	// Mock do cliente Redis para os testes
	mockRedisClient := &mockRedisClient{}

	// Caso de uso para teste
	limiter := usecase.CreateLimiterIPUseCase(mockRedisClient)

	t.Run("Limit not reached", func(t *testing.T) {
		// Configuração do caso de teste
		ip := "192.168.1.1"
		mockRedisClient.mockGet = "5" // Definindo o contador de requisições para 5

		// Execução do caso de teste
		result, err := limiter.CheckRateLimitByIP(context.Background(), ip)

		// Verificação do resultado
		assert.NoError(t, err)                       // Verifica se não houve erro
		assert.True(t, result)                       // Verifica se o resultado é true (limite não atingido)
		assert.Equal(t, ip, mockRedisClient.lastKey) // Verifica se a chave correta foi usada para recuperar o contador de requisições
	})

	t.Run("Limit reached", func(t *testing.T) {
		// Configuração do caso de teste
		ip := "192.168.1.2"
		mockRedisClient.mockGet = "100" // Definindo o contador de requisições para 100 (acima do limite)

		// Execução do caso de teste
		result, err := limiter.CheckRateLimitByIP(context.Background(), ip)

		// Verificação do resultado
		assert.NoError(t, err)                       // Verifica se não houve erro
		assert.False(t, result)                      // Verifica se o resultado é false (limite atingido)
		assert.Equal(t, ip, mockRedisClient.lastKey) // Verifica se a chave correta foi usada para recuperar o contador de requisições
	})

	t.Run("Error from Redis", func(t *testing.T) {
		// Configuração do caso de teste
		ip := "192.168.1.3"
		mockRedisClient.mockErr = errors.New("redis error") // Simulando um erro ao acessar o Redis

		// Execução do caso de teste
		result, err := limiter.CheckRateLimitByIP(context.Background(), ip)

		// Verificação do resultado
		assert.Error(t, err)                         // Verifica se ocorreu um erro
		assert.False(t, result)                      // Verifica se o resultado é false (limite atingido devido ao erro)
		assert.Equal(t, ip, mockRedisClient.lastKey) // Verifica se a chave correta foi usada para recuperar o contador de requisições
	})
}

// Mock do cliente Redis para os testes
type mockRedisClient struct {
	mockGet string // Mock para o valor retornado pelo método Get
	mockErr error  // Mock para o erro retornado pelo método Get
	lastKey string // Última chave usada no método Get
}

// Método mock para o método Get do cliente Redis
func (m *mockRedisClient) Get(ctx context.Context, key string) (string, error) {
	m.lastKey = key
	return m.mockGet, m.mockErr
}
