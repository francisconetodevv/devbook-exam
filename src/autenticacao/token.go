package autenticacao

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// CriarToken retorna um token assinado com as permissões do usuario
func CriarToken(usuarioID uint64) (string, error) {
	permisoes := jwt.MapClaims{}
	permisoes["authorized"] = true
	permisoes["exp"] = time.Now().Add(time.Hour * 6).Unix()
	permisoes["usuarioId"] = usuarioID

	// Secret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permisoes)
	return token.SignedString([]byte("Secret")) // Secret
}
