package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckRateLimitByToken_LimitNotReached(t *testing.T) {
	// Mock do repositório
	mockRepo := &mockRequestRepository{
		mockGetToken: true, // Simula que o limite não foi atingido
	}

	// Criar instância do use case com o mock repository
	usecase := NewRateLimiterByToken(mockRepo)

	// Chamar o método de verificação do limite
	reached, err := usecase.CheckRateLimitByToken(context.Background(), "sample_token")

	// Verificar se não ocorreu nenhum erro
	assert.NoError(t, err)

	// Verificar se o limite não foi atingido
	assert.False(t, reached)
}

func TestCheckRateLimitByToken_LimitReached(t *testing.T) {
	// Mock do repositório
	mockRepo := &mockRequestRepository{
		mockGetToken: false, // Simula que o limite foi atingido
	}

	// Criar instância do use case com o mock repository
	usecase := NewRateLimiterByToken(mockRepo)

	// Chamar o método de verificação do limite
	reached, err := usecase.CheckRateLimitByToken(context.Background(), "sample_token")

	// Verificar se não ocorreu nenhum erro
	assert.NoError(t, err)

	// Verificar se o limite foi atingido
	assert.True(t, reached)
}

func TestCheckRateLimitByToken_ErrorFromRepository(t *testing.T) {
	// Mock do repositório
	mockRepo := &mockRequestRepository{
		mockError: errors.New("error fetching token"), // Simula um erro ao buscar o token
	}

	// Criar instância do use case com o mock repository
	usecase := NewRateLimiterByToken(mockRepo)

	// Chamar o método de verificação do limite
	reached, err := usecase.CheckRateLimitByToken(context.Background(), "sample_token")

	// Verificar se ocorreu um erro e se foi o esperado
	assert.Error(t, err)
	assert.False(t, reached) // Se ocorrer um erro, o limite não foi atingido
}
