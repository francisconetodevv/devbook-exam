// Inicializando a API
package main

import (
	"api/src/config"
	"api/src/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Carregando as variáveis de ambiente
	config.Carregar()
	r := router.Gerar()

	fmt.Printf("Escutando na porta %d", config.Porta)

	// Subindo o server em uma porta
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Porta), r))
}
