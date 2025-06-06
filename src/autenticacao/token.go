package autenticacao

import (
	"api/src/config"
	"errors"
	"fmt"
	"net/http"
	"strings"
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
	return token.SignedString([]byte(config.SecretKey)) // Secret
}

// Validar token verificar se o token é valido
func ValidarToken(r *http.Request) error {
	tokenString := extrairToken(r)
	token, erro := jwt.Parse(tokenString, retornarChaveVerificacao)

	if erro != nil {
		return erro
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New("Token inválido")
}

func extrairToken(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if len(strings.Split(token, "")) == 2 {
		return strings.Split(token, "")[1]
	}

	return ""
}

func retornarChaveVerificacao(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Método de assinatura inesperado! %v", token.Header["alg"])
	}

	return config.SecretKey, nil
}
