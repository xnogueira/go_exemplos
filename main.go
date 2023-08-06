package main

// ! Nogueira com CRUD GO  -  Junho 2023

//todo: IMPLEMENTAR: ====================================
//todo: 1-No novo cadastro, verificar se nome e email já existe
//todo: 2-Poder mudar de mysql local para servidor
//todo: ====================================================

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

// ? cria acesso ao database:
func conexaoBD() (conexao *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Senha := "anju3BOM"
	Database := "mc2_kodular"

	conexao, err := sql.Open(Driver, Usuario+":"+Senha+"@tcp(127.0.0.1)/"+Database)

	if err != nil {
		panic(err.Error())
	}

	return conexao
}

// ? cria acesso a todos os templates:
var templates = template.Must(template.ParseGlob("templates/*"))

// inicio:
func main() {
	http.HandleFunc("/", Inicio)
	http.HandleFunc("/criar", Criar)
	http.HandleFunc("/inserir", Inserir)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/apagar", Apagar)
	http.HandleFunc("/atualizar", Atualizar)

	//abre servidor local:
	fmt.Println("Servidor em localhost:3000")
	http.ListenAndServe("localhost:3000", nil)
}

// ? cria estrutura da tabela:
type Teste struct {
	Id    int
	Nome  string
	Email string
}

// ? função inicial:
func Inicio(w http.ResponseWriter, r *http.Request) {
	conexaoOK := conexaoBD()
	registros, err := conexaoOK.Query("SELECT * FROM teste")
	if err != nil {
		panic(err.Error())
	}
	teste := Teste{}
	allTeste := []Teste{}

	for registros.Next() {
		var id int
		var nome, email string
		err = registros.Scan(&id, &nome, &email)
		if err != nil {
			panic(err.Error())
		}
		teste.Id = id
		teste.Nome = nome
		teste.Email = email

		allTeste = append(allTeste, teste)
	}
	templates.ExecuteTemplate(w, "inicio", allTeste)
}

// ? função de criação:
func Criar(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "criar", nil)
}

// ? função de inserção:
func Inserir(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nome := r.FormValue("nome")
		email := r.FormValue("email")

		conexaoOK := conexaoBD()
		inserirReg, err := conexaoOK.Prepare("INSERT INTO teste (nome,email) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		inserirReg.Exec(nome, email)

		http.Redirect(w, r, "/", 301)
	}
}

// ? função de apagar:
func Apagar(w http.ResponseWriter, r *http.Request) {
	idTeste := r.URL.Query().Get("id") //pega parametro ID

	conexaoOK := conexaoBD()
	apagarReg, err := conexaoOK.Prepare("DELETE FROM teste WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	apagarReg.Exec(idTeste)

	http.Redirect(w, r, "/", 301)
}

// ? função de edição:
func Editar(w http.ResponseWriter, r *http.Request) {
	idTeste := r.URL.Query().Get("id") //pega parametro ID

	conexaoOK := conexaoBD()
	editarReg, err := conexaoOK.Query("SELECT * FROM teste WHERE id=?", idTeste)
	if err != nil {
		panic(err.Error())
	}
	teste := Teste{}

	for editarReg.Next() {
		var id int
		var nome, email string
		err = editarReg.Scan(&id, &nome, &email)
		if err != nil {
			panic(err.Error())
		}
		teste.Id = id
		teste.Nome = nome
		teste.Email = email
	}
	templates.ExecuteTemplate(w, "editar", teste)
}

// ? função de atualização:
func Atualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		nome := r.FormValue("nome")
		email := r.FormValue("email")

		conexaoOK := conexaoBD()
		alterarReg, err := conexaoOK.Prepare("UPDATE teste SET nome=?, email=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		alterarReg.Exec(nome, email, id) //na mesma ordem da query

		http.Redirect(w, r, "/", 301)
	}
}

//! fim do programa
