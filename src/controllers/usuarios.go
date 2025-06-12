package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"api/src/seguranca"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// Criar usuario - Inserir usuario no banco de dados
func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	corpoRequest, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario modelos.Usuario
	if erro = json.Unmarshal(corpoRequest, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// Teste para validar se o JSON enviado pelo usuário está passando as infos no padrão, esse teste é realizado quando chama a func Preparar()
	if erro = usuario.Preparar("cadastro"); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuario(db)
	usuarioID, erro := repositorio.Criar(usuario)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	usuario.ID = usuarioID
	respostas.JSON(w, http.StatusCreated, usuario)
}

// Busca todos os usuários no banco
func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	nomeOuNick := strings.ToLower(r.URL.Query().Get("usuario"))

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuario(db)
	usuarios, erro := repositorio.Buscar(nomeOuNick)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, usuarios)
}

// Busca um único usuário no banco
func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuario(db)
	usuario, erro := repositorio.BuscarPorID(usuarioID)

	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, usuario)
}

// Altera as info de um usuário no banco
func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)

	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	usuarioIDNoToken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	if usuarioID != usuarioIDNoToken {
		respostas.Erro(w, http.StatusForbidden, errors.New("Não é possível atualizar um usuário que não seja o seu!"))
		return
	}

	corpoRequest, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario modelos.Usuario
	if erro = json.Unmarshal(corpoRequest, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = usuario.Preparar("edicao"); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuario(db)
	if erro = repositorio.Atualizar(usuarioID, usuario); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)

}

// Exclui as info de um usuário no banco
func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	usuarioIDNoToken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	if usuarioID != usuarioIDNoToken {
		respostas.Erro(w, http.StatusForbidden, errors.New("Não é possível deletar um usuário que não seja seu!"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuario(db)
	if erro = repositorio.Deletar(usuarioID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// SeguirUsuario permite que um usuário siga outro
func SeguirUsuario(w http.ResponseWriter, r *http.Request) {
	seguidorID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if seguidorID == usuarioID {
		respostas.Erro(w, http.StatusForbidden, errors.New("Não é possível seguir você mesmo"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuario(db)
	if erro = repositorio.Seguir(usuarioID, seguidorID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// ParaDeSeguirUsuario faz com que para-se de seguir outro usuario
func PararDeSeguirUsuario(w http.ResponseWriter, r *http.Request) {
	seguidorID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)

		return
	}

	if seguidorID == usuarioID {
		respostas.Erro(w, http.StatusForbidden, errors.New("Não é possível para de seguir você mesmo"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuario(db)
	if erro = repositorio.PararDeSeguirUsuario(usuarioID, seguidorID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// BuscarSeguidores busca todos os seguidors de um usuário
func BuscarSeguidores(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuario(db)
	seguidores, erro := repositorio.BuscarSeguidores(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, seguidores)

}

// BuscarSeguindo busca todos os usuários que um determinado usuário está seguindo
func BuscarSeguindo(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuario(db)
	usuarios, erro := repositorio.BuscarSeguindo(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, usuarios)
}

// AtualizarSenha atualiza a senha do usuário
func AtualizarSenha(w http.ResponseWriter, r *http.Request) {
	usuarioIDNoToken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if usuarioIDNoToken != usuarioID {
		respostas.Erro(w, http.StatusForbidden, errors.New("Não é possível atualizar uma senha de um usuário que não é seu!"))
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(r.Body)

	var senha modelos.Senha
	if erro = json.Unmarshal(corpoRequisicao, &senha); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuario(db)
	senhaSalvaNoBanco, erro := repositorio.BuscarSenha(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if erro = seguranca.VerificarSenha(senhaSalvaNoBanco, senha.Atual); erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, errors.New("A senha atual não é igual a salva no banco!"))
		return
	}

	senhaComHash, erro := seguranca.Hash(senha.Nova)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = repositorio.AtualizarSenha(usuarioID, string(senhaComHash)); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}
