package shared

import (
	"sync"
)

type TokenManager struct {
	tokenMap sync.Map
}

var tokenManager *TokenManager

func InitializeTokenManager() {
	tokenManager = &TokenManager{}
}

func (tm *TokenManager) AddToken(token string) {
	tm.tokenMap.Store(token, true)
}

func (tm *TokenManager) TokenExists(token string) bool {
	_, exist := tm.tokenMap.Load(token)
	return exist
}

func (tm *TokenManager) RemoveToken(token string) {
	tm.tokenMap.Delete(token)
}

func GetTokenManager() *TokenManager {
	return tokenManager
}
