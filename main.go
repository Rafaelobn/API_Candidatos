/* Essa API que cadastra, exibe e deleta dados de candidatos
no banco de dados com interface no browser em HTML */

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql" // Importando o driver MySQL
)

// Definindo o tipo "Candidatos"
type Candidatos struct {
	Id            int     `json:"id"`
	NomeCompleto  string  `json:"nomecompleto"`
	Email         string  `json:"email"`
	Idade         int     `json:"idade"`
	Formação      string  `json:"formação"`
	UltimoSalario float64 `json:"ultimosalario"`
	Telefone      string  `json:"telefone"`
}

var nextID int = 1
var db *sql.DB

// Inicia e testa a conexão com o DB MySQL
func initdb() (*sql.DB, error) {
	db, err := sql.Open("mysql", "user:password@tcp(host:port)/dbname") // ("Driver para o MySQL", "String de conexão")
	/* Alteração acima:
	   - Aqui, substitua "user" pelo seu nome de usuário do MySQL, "password" pela senha,
	     e "dbname" pelo nome do banco de dados que você quer usar. */
	if err != nil {
		log.Fatalf("Não foi possível conectar ao banco de dados MySQL: %s", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Erro ao verificar conexão com o banco: %s", err)
	}
	fmt.Println("Conectado ao banco de dados MySQL")
	return db, nil
	//"retorno"

}

func main() {

	// Cria a conexão com o banco de dados
	var err error
	db, err = initdb()
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}

	if err != nil {
		log.Fatalf("Erro ao criar tabela: %v", err)
	}

	// Criando rotas na porta local 8080
	http.HandleFunc("/candidatos", pagInicial)
	http.HandleFunc("/opcaoEscolhida", opcaoEscolhida)
	http.HandleFunc("/cadastrarCandidato", cadastrarCandidato)
	http.HandleFunc("/salvarCandidato", salvarCandidato)
	http.HandleFunc("/exibirCandidatos", exibirCandidatos)
	http.HandleFunc("/excluirCandidato", excluirCandidato)
	fmt.Println("Servidor iniciado na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil)) // Inicia na localhost de porta 8080
}

func pagInicial(w http.ResponseWriter, r *http.Request) {
	// Printando em HTML as opções
	fmt.Fprintf(w, `
        <html>
        <body>
            <h1>Bem-vindo ao banco de dados dos candidatos!</h1>
            <form action="/opcaoEscolhida" method="GET">
                <label for="opcao">Escolha uma das opções abaixo:</label><br>
                <input type="radio" id="opcao1" name="opcao" value="1">
                <label for="opcao1">1 - Cadastrar novo candidato</label><br>
                <input type="radio" id="opcao2" name="opcao" value="2">
                <label for="opcao2">2 - Listar candidatos já cadastrados</label><br>
                <input type="radio" id="opcao3" name="opcao" value="3">
                <label for="opcao3">3 - Excluir candidatos cadastrados</label><br>
                <input type="submit" value="Enviar">
            </form>
        </body>
        </html>`)
}

// Tratamento da função escolhida
func opcaoEscolhida(w http.ResponseWriter, r *http.Request) {
	opcao := r.URL.Query().Get("opcao") // Captura a opção escolhida

	switch opcao {
	case "1":
		http.Redirect(w, r, "/cadastrarCandidato", http.StatusSeeOther)
	case "2":
		http.Redirect(w, r, "/exibirCandidatos", http.StatusSeeOther)
	case "3":
		http.Redirect(w, r, "/excluirCandidato", http.StatusSeeOther)
	default:
		http.Error(w, "Opção inválida", http.StatusBadRequest)
	}
}

// Função para cadastrar novos candidatos no DB
func cadastrarCandidato(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Exibe, em HTML, os campos a serem preenchidos
		fmt.Fprintf(w, `
            <html>
            <body>
                <h1>Cadastrar Novo Candidato</h1>
                <form action="/salvarCandidato" method="POST">
                    <label for="nomeCompleto">Nome Completo:</label><br>
                    <input type="text" id="nomeCompleto" name="nomeCompleto" required><br>
                    <label for="email">Email:</label><br>
                    <input type="email" id="email" name="email" required><br>
                    <label for="idade">Idade:</label><br>
                    <input type="number" id="idade" name="idade" required><br>
                    <label for="formacao">Formação:</label><br>
                    <input type="text" id="formacao" name="formacao" required><br>
                    <label for="ultimoSalario">Último Salário:</label><br>
                    <input type="number" step="0.01" id="ultimoSalario" name="ultimoSalario" required><br>
                    <label for="telefone">Telefone:</label><br>
                    <input type="text" id="telefone" name="telefone" required><br><br>
                    <input type="submit" value="Salvar">
                </form>
            </body>
            </html>`)
	}
}

func salvarCandidato(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		// Processar os dados do candidato
		idade, err := strconv.Atoi(r.FormValue("idade")) // Convertendo idade para int
		if err != nil {
			http.Error(w, "Idade inválida", http.StatusBadRequest)
			return
		}

		ultimoSalario, err := strconv.ParseFloat(r.FormValue("ultimoSalario"), 64) // Convertendo último salário para float64
		if err != nil {
			http.Error(w, "Último Salário inválido", http.StatusBadRequest)
			return
		}

		// Verifica se db está conectado
		err = db.Ping()
		if err != nil {
			http.Error(w, "Erro de conexão com o banco de dados", http.StatusInternalServerError)
			return
		}

		// Insere os dados no DB
		_, err = db.Exec("INSERT INTO candidatos (nome, email, idade, formacao, salario, telefone) VALUES (?, ?, ?, ?, ?, ?)",
			r.FormValue("nomeCompleto"), r.FormValue("email"), idade, r.FormValue("formacao"), ultimoSalario, r.FormValue("telefone"))
		if err != nil {
			http.Error(w, "Erro ao salvar candidato no banco de dados", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Candidato %s cadastrado com sucesso!", r.FormValue("nomeCompleto"))
	}
}

func exibirCandidatos(w http.ResponseWriter, r *http.Request) {
	// Executa a consulta SQL para buscar todos os candidatos no DB
	rows, err := db.Query("SELECT id, nome, email, idade, formacao, salario, telefone FROM candidatos")
	if err != nil {
		http.Error(w, "Erro ao buscar candidatos no banco de dados", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Cria a tabela HTML
	fmt.Fprintf(w, `
		<html>
		<body>
			<h1>Lista de Candidatos</h1>
			<table border="1">
				<tr>
					<th>ID</th>
					<th>Nome Completo</th>
					<th>Email</th>
					<th>Idade</th>
					<th>Formação</th>
					<th>Último Salário</th>
					<th>Telefone</th>
				</tr>`)

	// Itera sobre os resultados da consulta e exibe cada candidato na tabela HTML
	for rows.Next() {
		var candidato Candidatos
		err := rows.Scan(&candidato.Id, &candidato.NomeCompleto, &candidato.Email, &candidato.Idade, &candidato.Formação, &candidato.UltimoSalario, &candidato.Telefone)
		if err != nil {
			http.Error(w, "Erro ao processar dados dos candidatos", http.StatusInternalServerError)
			return
		}

		// Adiciona cada candidato na tabela
		fmt.Fprintf(w, `
			<tr>
				<td>%d</td>
				<td>%s</td>
				<td>%s</td>
				<td>%d</td>
				<td>%s</td>
				<td>%.2f</td>
				<td>%s</td>
			</tr>`, candidato.Id, candidato.NomeCompleto, candidato.Email, candidato.Idade, candidato.Formação, candidato.UltimoSalario, candidato.Telefone)
	}

	// Finaliza a tabela e a página HTML
	fmt.Fprintf(w, `
			</table>
		</body>
		</html>`)
}

func excluirCandidato(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Exibe o HTML para o usuário inserir o ID do candidato a ser excluído
		fmt.Fprintf(w, `
            <html>
            <body>
                <h1>Excluir Candidato</h1>
                <form action="/excluirCandidato" method="POST">
                    <label for="id">ID do Candidato:</label><br>
                    <input type="number" id="id" name="id" required><br><br>
                    <input type="submit" value="Excluir">
                </form>
            </body>
            </html>`)
	} else if r.Method == http.MethodPost {
		// Processa a exclusão após o envio do ID
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		// Exclui o usuário do Banco de dados
		result, err := db.Exec("DELETE FROM candidatos WHERE id = ?", id)
		if err != nil {
			http.Error(w, "Erro ao excluir candidato do banco de dados", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, "Erro ao verificar exclusão", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			fmt.Fprintf(w, "Nenhum candidato encontrado com o ID %d.", id)
		} else {
			fmt.Fprintf(w, "Candidato do ID %d foi excluído com sucesso!", id)
		}
	}
}
