package modelos

import "time"

// Usuário representa um usuário utilizando a rede social
type Usuario struct {
	ID        uint64    `json:"id,omitempty"`
	Nome      string    `json:"nome,omitempty"`
	Nick      string    `json:"nick,omitempty"`
	Email     string    `json:"email,omitempty"`
	Senha     string    `json:"senha,omitempty"`
	CreatedOn time.Time `json:"createdOn,omitempty"`
}
