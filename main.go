// Inicializando a API
package main

import (
	"api/src/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Rodando API")
	r := router.Gerar()

	// Subindo o server em uma porta
	log.Fatal(http.ListenAndServe(":5000", r))
}
