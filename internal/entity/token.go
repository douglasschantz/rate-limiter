package entity

type Token struct {
	value  string
	userID string
}

func (t *Token) CheckRateLimitToken(count, limit int) bool {
	return count >= limit
}
