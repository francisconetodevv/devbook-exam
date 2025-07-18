package modelos

import (
	"errors"
	"strings"
	"time"
)

// Publicação representa uma publicação feita por um usuário
type Publicacao struct {
	ID        uint64    `json:"id,omitempy"`
	Titulo    string    `json:"titulo,omitempy"`
	Conteudo  string    `json:"conteudo,omitempy"`
	AutorID   uint64    `json:"autorId,omitempy"`
	AutorNick string    `json:"autorNick,omitempy"`
	Curtidas  uint64    `json:"curtidas"`
	CriadaEm  time.Time `json:"criadaEm,omitempy"`
}

func (publicacao *Publicacao) Preparar() error {
	if erro := publicacao.validar(); erro != nil {
		return erro
	}

	publicacao.formatar()
	return nil
}

func (publicacao *Publicacao) validar() error {
	if publicacao.Titulo == "" {
		return errors.New("O título é obrigatório e não pode está em branco!")
	}

	if publicacao.Conteudo == "" {
		return errors.New("O conteudo é obrigatório e não pode estar em branco!")
	}

	return nil
}

func (publicacao *Publicacao) formatar() {
	publicacao.Titulo = strings.TrimSpace(publicacao.Titulo)
	publicacao.Conteudo = strings.TrimSpace(publicacao.Conteudo)
}
