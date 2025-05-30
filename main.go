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
	fmt.Println(config.StringConexaoBanco)

	fmt.Println("Rodando API")
	r := router.Gerar()

	// Subindo o server em uma porta
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Porta), r))
}
